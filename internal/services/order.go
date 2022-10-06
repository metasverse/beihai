package services

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/google/uuid"

	"lihood/conf"
	"lihood/g"
	"lihood/internal/enum"
	"lihood/internal/models"
	"lihood/internal/repository"
	"lihood/pkg/chain"
	"lihood/pkg/pay"
)

type OrderService interface {
	NewOrder(uid int64, pid int64, payType enum.PayType) ([]byte, error)
	OrderCallback(orderID string) error
	//QueryOrder
}

func NewOrderService(session g.Session) OrderService {
	return &orderService{session: session}
}

type orderService struct {
	session g.Session
}

func (o orderService) OrderCallback(orderID string) error {
	orderRepo := repository.NewProductOrderRepository(o.session)
	order, err := orderRepo.QueryByOID(orderID)
	if err != nil {
		return err
	}
	if order.Status == enum.ProductPaid {
		return g.Error("订单已支付")
	}
	// 修改订单的状态
	if err := orderRepo.UpdateStatusByID(order.ID, enum.ProductPaid); err != nil {
		return err
	}
	// 查询作品信息
	userRepo := repository.NewUserProductRepository(o.session)
	userPro, err := userRepo.GetByID(order.PID)
	if err != nil {
		return err
	}
	product, err := repository.NewProductRepository(o.session).QueryByID(userPro.PID)
	if err != nil {
		return err
	}
	// 增加拥有者的的收入记录
	income := models.AccountIncome{
		UID:        userPro.UID,
		Type:       enum.Income,
		Amount:     product.Price,
		Remark:     fmt.Sprintf("%s 作品卖出收入", product.Name),
		CreateTime: time.Now().Unix(),
	}
	if err := repository.NewAccountIncomeRepository(o.session).Create(&income); err != nil {
		return err
	}
	// 增加拥有者的账户余额
	if err = repository.NewAccountRepository(o.session).UpdateAmountById(userPro.UID, product.Price); err != nil {
		return err
	}
	// 给当前购买这增加一笔支出记录
	expense := models.AccountIncome{
		UID:        order.UID,
		Type:       enum.Expense,
		Amount:     product.Price,
		Remark:     fmt.Sprintf("%s 作品购入", product.Name),
		CreateTime: time.Now().Unix(),
	}
	if err := repository.NewAccountIncomeRepository(o.session).Create(&expense); err != nil {
		return err
	}
	count, err := userRepo.CountByPID(product.ID)
	// 给购买者增加一个作品
	item := models.UserProduct{
		PID:        order.PID,
		UID:        order.UID,
		Times:      count + 1,
		Status:     true,
		CreateTime: time.Now().Unix(),
		CID:        uuid.New().String(),
	}
	account, err := repository.NewAccountRepository(o.session).GetByID(order.UID)
	if err != nil {
		return err
	}
	// 上链
	client := chain.NewChainClient()
	cb := g.ChainCallback(item.CID)
	resp, err := client.NewProduct(item.CID, product.Image, account.BsnAddress, product.Description, cb)
	if err != nil {
		return err
	}
	if !resp.Success {
		log.Println("上链失败")
		return errors.New(resp.ErrMsg)
	}
	item.TxID = resp.Data.TxId
	item.TokenID = resp.Data.TokenId
	return userRepo.Create(&item)
}

func (o orderService) NewOrder(uid int64, pid int64, payType enum.PayType) ([]byte, error) {
	user, err := repository.NewAccountRepository(o.session).GetByID(uid)
	if err != nil {
		return nil, err
	}
	if user.IDCardNum == "" {
		return nil, g.Error("请先实名认证")
	}
	// 先判断当前的商品有没有存库
	fmt.Println(uid, pid, payType)
	userProRepo := repository.NewUserProductRepository(o.session)
	userPro, err := userProRepo.GetByID(pid)
	if err == sql.ErrNoRows {
		return nil, g.Error("商品不存在")
	}
	if err != nil {
		return nil, err
	}
	// 查询相关作品信息
	product, err := repository.NewProductRepository(o.session).QueryByID(userPro.PID)
	if err != nil {
		return nil, err
	}
	// 判断库存是否足够
	count, err := userProRepo.CountByPID(userPro.PID)
	if err != nil {
		return nil, err
	}
	if count-1 >= product.Stock {
		return nil, errors.New("该作品已经销罄")
	}
	// 创建订单
	order := models.ProductOrder{
		OID:        uuid.New().String()[:8],
		PayType:    payType,
		PID:        pid,
		UID:        uid,
		CreateTime: time.Now().Unix(),
	}
	if err := repository.NewProductOrderRepository(o.session).Create(&order); err != nil {
		return nil, err
	}

	cb := conf.Instance.Server.Domain + fmt.Sprintf("/api/v1/order/callback/%s", order.OID)

	// 100 => 1.00
	price := strconv.Itoa(int(product.Price * 100))
	if len(price) <= 2 {
		price = "0." + price
	} else {
		price = price[:len(price)-2] + "." + price[len(price)-2:]
	}
	// 取前4位

	desc := product.Name[:4]
	fmt.Println(desc, user.IDCardNum, user.Name, price, cb, "2")
	resp, err := pay.HfPay(desc, user.IDCardNum, user.Name, price, cb, "2")
	if err != nil {
		return nil, err
	}

	return resp, nil
}
