package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/eatmoreapple/sqlbuilder"
	"github.com/google/uuid"

	"lihood/g"
	"lihood/internal/entity"
	"lihood/internal/enum"
	"lihood/internal/models"
	"lihood/internal/repository"
	"lihood/pkg/chain"
	"lihood/utils"
)

type ProductService interface {
	CreateProduct(product *models.Product, payType enum.PayType) (g.ResponseWriter, error)
	ProductDetail(id int64, uid int64) (g.ResponseWriter, error)
	ProductLike(id int64, uid int64) (g.ResponseWriter, error)
	// QueryList 查询列表
	QueryList(uid int64, page, pageSize int, orderBy string) (g.ResponseWriter, error)
	// SalesRank 销售排行榜
	SalesRank(page, pageSize int) (g.ResponseWriter, error)
	// PublishList  发布列表
	PublishList(uid int64, page, pageSize int) (g.ResponseWriter, error)
	// BuyList 购买列表
	BuyList(uid int64, page, pageSize int) (g.ResponseWriter, error)
	// LikeList 点赞列表
	LikeList(uid int64, page, pageSize int) (g.ResponseWriter, error)
	// SaleList 出售列表
	SaleList(uid int64, page, pageSize int) (g.ResponseWriter, error)
	// PublicCreateProduct 公开开放铸件
	PublicCreateProduct(pid string, product *models.Product) (g.ResponseWriter, error)
}

func NewProductService(session g.Session) ProductService {
	return &productService{session: session}
}

type productService struct {
	session g.Session
}

func (p productService) SaleList(uid int64, page, pageSize int) (g.ResponseWriter, error) {
	// 先查询订单状态是已支付的并且是当前用户的
	limit, offset := pageSize, (page-1)*pageSize
	ids, err := repository.NewProductOrderRepository(p.session).QueryProductId(uid, limit, offset)
	if err != nil {
		return nil, err
	}
	if len(ids) == 0 {
		return g.EmptyListRespWriter, nil
	}
	count, err := repository.NewProductOrderRepository(p.session).CountByUID(uid)
	if err != nil {
		return nil, err
	}
	result, err := repository.NewUserProductRepository(p.session).QueryProductListByIDs(uid, ids)
	if err != nil {
		return nil, err
	}
	return g.NewListResponseWriter(result, count), nil
}

func (p productService) PublicCreateProduct(pid string, product *models.Product) (g.ResponseWriter, error) {
	// 获取当前的创作者
	user, err := repository.NewAccountRepository(p.session).GetByID(product.AuthorID)
	if err != nil {
		return nil, err
	}
	product.Status = enum.ProductPaid

	orderId := uuid.New().String()
	product.OrderNo = orderId
	// 上链
	cid := uuid.New().String()
	client := chain.NewChainClient()
	cb := g.ChainCallback(cid)
	resp, err := client.NewProduct(pid, product.Image, user.BsnAddress, product.Description, cb)
	if err != nil {
		return nil, err
	}
	product.TokenID = resp.Data.TokenId
	product.TxID = resp.Data.TxId
	// 保存到数据库
	if err = repository.NewProductRepository(p.session).Create(product); err != nil {
		return nil, err
	}
	// 创建一条拥有者记录
	history := &models.UserProduct{
		PID:        product.ID,
		UID:        user.ID,
		CreateTime: time.Now().Unix(),
		Times:      1,
		CID:        cid,
		TokenID:    resp.Data.TokenId,
		TxID:       resp.Data.TxId,
		Display:    true,
		CName:      fmt.Sprintf("%s #%d", product.Cname, 1),
	}
	// todo 这块改掉
	if strings.HasPrefix(product.Cname, "KT") {
		history.IsAirDrop = true
	}
	if err = repository.NewUserProductRepository(p.session).Create(history); err != nil {
		return nil, err
	}
	return g.NewRespWriter[any](nil), nil
}

func (p productService) LikeList(uid int64, page, pageSize int) (g.ResponseWriter, error) {
	limit, offset := pageSize, (page-1)*pageSize
	repo := repository.NewProductLikesRepository(p.session)
	likes, err := repo.QueryByUID(uid, limit, offset)
	if err == sql.ErrNoRows || len(likes) == 0 {
		return g.EmptyListRespWriter, nil
	}
	count, err := repo.CountByUID(uid)
	if err != nil {
		return nil, err
	}
	var ids []int64
	for _, like := range likes {
		ids = append(ids, like.PID)
	}
	fmt.Println(ids)
	result, err := repository.NewUserProductRepository(p.session).QueryProductListByIDs(uid, ids)
	if err != nil {
		return nil, err
	}
	return g.NewListResponseWriter(result, count), nil
}

func (p productService) BuyList(uid int64, page, pageSize int) (g.ResponseWriter, error) {
	limit, offset := pageSize, (page-1)*pageSize
	repo := repository.NewUserProductRepository(p.session)
	histories, err := repo.QueryNotFirstRecordByUID(uid, limit, offset)
	if err == sql.ErrNoRows || len(histories) == 0 {
		return g.EmptyListRespWriter, nil
	}
	count, err := repo.QueryNotFirstRecordCountByUID(uid)
	if err != nil {
		return nil, err
	}
	var ids []int64
	for _, history := range histories {
		ids = append(ids, history.ID)
	}
	result, err := repo.QueryProductListByIDs(uid, ids)
	if err != nil {
		return nil, err
	}
	return g.NewListResponseWriter(result, count), nil
}

func (p productService) PublishList(uid int64, page, pageSize int) (g.ResponseWriter, error) {
	limit, offset := pageSize, (page-1)*pageSize
	repo := repository.NewUserProductRepository(p.session)
	histories, err := repo.QueryFirstRecordByUID(uid, limit, offset)
	if err == sql.ErrNoRows || len(histories) == 0 {
		return g.EmptyListRespWriter, nil
	}
	count, err := repo.QueryFirstRecordCountByUID(uid)
	if err != nil {
		return nil, err
	}
	var ids []int64
	for _, history := range histories {
		ids = append(ids, history.ID)
	}
	result, err := repo.QueryProductListByIDs(uid, ids)
	if err != nil {
		return nil, err
	}
	return g.NewListResponseWriter(result, count), nil
}

func (p productService) SalesRank(page, pageSize int) (g.ResponseWriter, error) {
	limit, offset := pageSize, (page-1)*pageSize
	result, err := repository.NewProductRepository(p.session).SalesRank(limit, offset)
	if err == sql.ErrNoRows {
		return g.NewRespWriter[[]int](make([]int, 0)), nil
	}
	return g.NewRespWriter[[]*entity.SalesRank](result), nil
}

func (p productService) QueryList(uid int64, page, pageSize int, orderBy string) (g.ResponseWriter, error) {
	limit, offset := pageSize, pageSize*(page-1)
	var idList []int64
	var count int64
	if strings.Contains(orderBy, "index") {
		// 根据推荐列表查询
		repoRepo := repository.NewRecommendRepository(p.session)
		recommends, err := repoRepo.QueryList(limit, offset)
		if err == sql.ErrNoRows {
			return g.EmptyListRespWriter, nil
		}
		if err != nil {
			return nil, err
		}
		if len(recommends) == 0 {
			return g.EmptyListRespWriter, nil
		}
		count, err = repoRepo.Count()
		if err != nil {
			return nil, err
		}
		for _, item := range recommends {
			idList = append(idList, item.ProductID)
		}
	} else {
		var orderExpr string
		switch orderBy {
		case "price":
			orderExpr = "price"
		// 根据价格查询
		case "-price":
			orderExpr = "price DESC"
		// 根据价格倒叙
		case "likes":
			orderExpr = "likes"
		// 根据点赞查询
		case "-likes":
			orderExpr = "likes DESC"
		// 根据点赞倒序
		case "create_time": // 根据创建时间查询
			orderExpr = "create_time"
		case "-create_time": // 根据创建时间倒序
			orderExpr = "create_time DESC"
		default:
			orderExpr = "create_time"
		}
		// 直接这里写
		builder := sqlbuilder.NewSelect("?")
		builder.From(models.UserProduct{}.TableName(), "a")
		builder.LeftJoin(models.Product{}.TableName(), "a.pid = b.id", "b")
		builder.LeftJoin(models.ProductLikes{}.TableName(), "a.pid = c.pid", "c")
		builder.Fields("a.id id", "b.price price", sqlbuilder.Count("c.id", "likes"), "a.create_time create_time")
		builder.Where("a.status = 1").And("a.display = 1")
		builder.GroupBy("a.id")
		builder.Limit(limit).Offset(offset).OrderBy(orderExpr)
		type Result struct {
			Id         int64 `column:"id"`
			Price      int64 `column:"price"`
			Likes      int64 `column:"likes"`
			CreateTime int64 `column:"create_time"`
		}
		results, err := sqlbuilder.BuilderScanner[Result](p.session, builder).List(context.Background())
		if err == sql.ErrNoRows {
			return g.EmptyListRespWriter, nil
		}
		for _, result := range results {
			idList = append(idList, result.Id)
		}
		count, err = repository.NewUserProductRepository(p.session).Count()
		if err != nil {
			return nil, err
		}
	}
	// 根据historyIDList去查列表
	result, err := repository.NewUserProductRepository(p.session).QueryProductListByIDs(uid, idList)
	if err != nil {
		return nil, err
	}
	newResult := make([]*entity.ProductList, len(result))
	for index, id := range idList {
		for _, item := range result {
			if int64(item.ID) == id {
				newResult[index] = item
				break
			}
		}
	}
	// 返回result
	return g.NewListResponseWriter(newResult, count), nil
}

// ProductLike 用户点赞和取消点赞
func (p productService) ProductLike(id int64, uid int64) (g.ResponseWriter, error) {
	dao := repository.NewProductLikesRepository(p.session)
	item, err := dao.GetByUIDAndPID(uid, id)
	if err == sql.ErrNoRows {
		// 创建新的
		model := models.ProductLikes{UID: uid, PID: id, CreateTime: time.Now().Unix()}
		return g.NewRespWriter[any](nil), dao.Create(&model)
	}
	if err != nil {
		return nil, err
	}
	if item.UID != uid {
		return nil, g.Error("没有权限")
	}
	return g.NewRespWriter[any](nil), dao.DeleteByID(item.ID)
}

// ProductDetail 获取作品详情
func (p productService) ProductDetail(id int64, uid int64) (g.ResponseWriter, error) {
	// 根据传来的history来查询对应作品的详情
	history, err := repository.NewUserProductRepository(p.session).GetByID(id)
	if err != nil {
		return nil, err
	}
	// 根据history的pid来查询作品信息
	product, err := repository.NewProductRepository(p.session).QueryByID(history.PID)
	if err != nil {
		return nil, err
	}
	var result = entity.ProductDetail{
		ID:          id,
		Name:        product.Name,
		Price:       product.Price,
		Stock:       product.Stock,
		Image:       product.Image,
		AuthorID:    product.AuthorID,
		OwnerID:     history.UID,
		Description: product.Description,
		CreateTime:  product.CreateTime,
		Hash:        history.Hash,
	}

	// 查询商品点赞数据
	likeRepo := repository.NewProductLikesRepository(p.session)
	result.Likes, err = likeRepo.CountByPID(id)
	if err != nil {
		return nil, err
	}

	// 查询自己有没有点赞
	_, err = likeRepo.GetByUIDAndPID(uid, id)
	result.Liked = err == nil

	// 查询作者和拥有者相关信息
	accountRepo := repository.NewAccountRepository(p.session)
	// 查询作者
	author, err := accountRepo.GetByID(product.AuthorID)
	if err != nil {
		return nil, err
	}
	result.AuthorName = author.Nickname
	result.AuthorAvatar = author.Avatar
	// 查询拥有者
	owner, err := accountRepo.GetByID(history.UID)
	if err != nil {
		return nil, err
	}
	result.OwnerName = owner.Nickname
	result.AuthorDesc = owner.Description
	result.OwnerAvatar = owner.Avatar
	result.TokenID = "#" + utils.ZeroFill(history.Times, len(strconv.FormatInt(product.Stock, 10)))

	result.Sales, err = repository.NewUserProductRepository(p.session).CountByPID(id)
	if err != nil {
		return nil, err
	}
	result.Sales = result.Sales - 1

	// 查询这个作品是否可以购买
	now := time.Now().Unix()

	if now > product.SaleTime {
		result.CanBuy = true
	} else {
		// 判断是否到了提前购的时间
		if now > product.SaleTime-int64(product.AdvanceHour*3600) {

		}
		// 查询当前用户是否在白名单中
		ok, err := repository.NewWhiteListRepository(p.session).IsWhiteList(id, uid)
		if err != nil {
			log.Println(err)
		}
		if ok {
			result.CanBuy = now+int64(product.AdvanceHour*3600) > product.SaleTime
		}
	}

	// 查询自己有没有购买过
	ok, err := repository.NewUserProductRepository(p.session).ExistsByUIDAndPid(uid, history.PID)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
	}
	if ok {
		result.CanBuy = false
	}

	return g.NewRespWriter(result), nil
}

// CreateProduct 创建作品
func (p productService) CreateProduct(product *models.Product, payType enum.PayType) (g.ResponseWriter, error) {
	// 保存到数据库
	if err := repository.NewProductRepository(p.session).Create(product); err != nil {
		return nil, err
	}

	return nil, nil
}

//type ProductService interface {
//	QueryList(uid int64, page, pageSize int, orderBy string) (g.ResponseWriter, error)
//	Count() (int64, error)
//	CreateProduct(product *models.Product, payType enum.PayType) (interface{}, error)
//	QueryByID(pk int64) (product *models.Product, err error)
//	Detail(uid int64, pk int64) (*entity.ProductDetail, error)
//	PublishList(uid int64, page, pageSize int) ([]*entity.ProductList, error)
//	PublishCount(uid int64) (int64, error)
//	BuyList(uid int64, page, pageSize int) ([]*entity.ProductList, error)
//	BuyCount(uid int64) (int64, error)
//	LikeList(uid int64, page, pageSize int) ([]*entity.ProductList, error)
//	OwnerList(uid int64, page, pageSize int) ([]*entity.ProductList, error)
//	SalesRank(page, pageSize int) ([]*entity.SalesRank, error)
//	ActiveProduct(pid int64) error
//	OwnerListCount(ud int64) (int64, error)
//	ActiveProductChainStatusByOrderNo(orderNo string, hash string) error
//	PublicCreateProduct(product *models.Product) (interface{}, error)
//}
//
//func NewProductService(session g.Session) ProductService {
//	return &productService{session: session}
//}
//
//type productService struct {
//	session g.Session
//}
//
//func (p productService) PublicCreateProduct(product *models.Product) (interface{}, error) {
//	// 获取当前的创作者
//	user, err := repository.NewAccountRepository(p.session).GetByID(product.AuthorID)
//	if err != nil {
//		return nil, err
//	}
//	orderNo := strings.ReplaceAll(uuid.New().String(), "-", "")[:15] + "-1"
//	client := chain.NewChainClient()
//	resp, err := client.NewProduct(product.Image, user.BsnAddress, product.Description, conf.Instance.Server.Domain+"/api/v1/product/create/callback/"+orderNo)
//	if err != nil {
//		return nil, err
//	}
//	if !resp.Success {
//		return nil, errors.New(resp.ErrMsg)
//	}
//	product.TxID = resp.Data.TxId
//	product.TokenID = resp.Data.TokenId
//
//	// todo 暂时在上线之后删除
//	{
//		product.Status = enum.ProductPaid
//	}
//
//	// 创建作品
//	if err = repository.NewProductRepository(p.session).Create(product); err != nil {
//		return nil, err
//	}
//
//	// todo 删除这里
//	{
//		//if err = p.ActiveProduct(product.ID); err != nil {
//		//	return nil, err
//		//}
//	}
//	// 创建一条拥有者记录
//	if err = repository.NewUserProductRepository(p.session).Create(&models.UserProduct{
//		PID:        product.ID,
//		UID:        user.ID,
//		CreateTime: time.Now().Unix(),
//	}); err != nil {
//		return nil, err
//	}
//	return "ok", nil
//}
//
//func (p productService) ActiveProductChainStatusByOrderNo(orderNo string, hash string) error {
//	repo := repository.NewProductRepository(p.session)
//	product, err := repo.GetProductByOrderID(orderNo)
//	if err != nil {
//		return err
//	}
//	if err = repo.UpdateChainStatus(product.ID, enum.ChainPaid); err != nil {
//		return err
//	}
//	if err = repo.UpdateHashByID(product.ID, hash); err != nil {
//		return err
//	}
//	return nil
//}
//
//func (p productService) OwnerList(uid int64, page, pageSize int) ([]*entity.ProductList, error) {
//	limit, offset := pageSize, pageSize*(page-1)
//	pids, err := repository.NewUserProductRepository(p.session).GetOwnerPIDList(uid, limit, offset)
//	if err == sql.ErrNoRows || len(pids) == 0 {
//		return make([]*entity.ProductList, 0), nil
//	}
//	result, err := repository.NewProductRepository(p.session).QueryByIDList(uid, pids)
//	if err == sql.ErrNoRows {
//		return make([]*entity.ProductList, 0), nil
//	}
//	return result, err
//}
//
//func (p productService) OwnerListCount(ud int64) (int64, error) {
//	return repository.NewUserProductRepository(p.session).CountOwnerList(ud)
//}
//
//func (p productService) ActiveProduct(pid int64) error {
//	productRepo := repository.NewProductRepository(p.session)
//	product, err := productRepo.QueryByID(pid)
//	if err != sql.ErrNoRows {
//		return g.Error("商品不存在")
//	}
//	if err != nil {
//		return err
//	}
//	if err = productRepo.ChangeStatusByID(pid, enum.ProductPaid); err != nil {
//		return err
//	}
//	// 给当前的作者创建一笔支出记录
//	token := "#" + utils.ZeroFill(1, len(strconv.FormatInt(product.Stock, 10)))
//	income := models.AccountIncome{
//		UID:        product.AuthorID,
//		Type:       enum.Expense,
//		Amount:     1900,
//		Remark:     fmt.Sprintf("%s %s 作品铸件", token, product.Name),
//		CreateTime: time.Now().Unix(),
//	}
//	if err = repository.NewAccountIncomeRepository(p.session).Create(&income); err != nil {
//		return err
//	}
//	return nil
//}
//
//func (p productService) SalesRank(page, pageSize int) ([]*entity.SalesRank, error) {
//	limit, offset := pageSize, pageSize*(page-1)
//	result, err := repository.NewProductRepository(p.session).SalesRank(limit, offset)
//	if err == sql.ErrNoRows {
//		return make([]*entity.SalesRank, 0), nil
//	}
//	return result, err
//}
//
//// LikeList 用户点赞列表
//func (p productService) LikeList(uid int64, page, pageSize int) ([]*entity.ProductList, error) {
//	// 先查询去重复过的当前用户点赞过的pid
//	limit, offset := pageSize, (page-1)*pageSize
//	likes, err := repository.NewProductLikesRepository(p.session).QueryByUID(uid, limit, offset)
//	if err != nil {
//		return nil, err
//	}
//	idx := make([]int64, 0)
//	for _, like := range likes {
//		idx = append(idx, like.PID)
//	}
//	result, err := repository.NewProductRepository(p.session).QueryByIDList(uid, idx)
//	if err == sql.ErrNoRows {
//		return make([]*entity.ProductList, 0), nil
//	}
//	return result, err
//}
//
//func (p productService) BuyCount(uid int64) (int64, error) {
//	return repository.NewProductOrderRepository(p.session).CountByUID(uid)
//}
//
//func (p productService) BuyList(uid int64, page, pageSize int) ([]*entity.ProductList, error) {
//	// 先将订单列表分页
//	limit, offset := pageSize, pageSize*(page-1)
//	orderList, err := repository.NewProductOrderRepository(p.session).QueryProductId(uid, limit, offset)
//	if err != nil && err != sql.ErrNoRows {
//		return nil, err
//	}
//	// 再根据订单列表查询商品列表
//	result, err := repository.NewProductRepository(p.session).QueryByIDList(uid, orderList)
//	if err != nil {
//		return nil, err
//	}
//	return result, nil
//}
//
//func (p productService) PublishCount(uid int64) (int64, error) {
//	return repository.NewProductRepository(p.session).PublishCount(uid)
//}
//
//func (p productService) PublishList(uid int64, page, pageSize int) ([]*entity.ProductList, error) {
//	limit, offset := pageSize, pageSize*(page-1)
//	result, err := repository.NewProductRepository(p.session).QueryPublishList(uid, limit, offset)
//	if err == sql.ErrNoRows {
//		return make([]*entity.ProductList, 0), nil
//	}
//	return result, err
//}
//
//func (p productService) Detail(uid, pk int64) (*entity.ProductDetail, error) {
//	item, err := repository.NewProductRepository(p.session).QueryByID(pk)
//	if err == sql.ErrNoRows {
//		return nil, g.Error("没有找到该商品")
//	}
//	if err != nil {
//		return nil, err
//	}
//	var result = entity.ProductDetail{
//		ID:          pk,
//		Name:        item.Name,
//		Price:       item.Price,
//		Stock:       item.Stock,
//		Image:       item.Image,
//		AuthorID:    item.AuthorID,
//		Description: item.Description,
//		CreateTime:  item.CreateTime,
//	}
//	// 查询商品的点赞数
//	likeDao := repository.NewProductLikesRepository(p.session)
//	likeCount, err := likeDao.CountByPID(pk)
//	if err != nil {
//		return nil, err
//	}
//	result.Likes = likeCount
//	_, err = likeDao.GetByUIDAndPID(uid, pk)
//	if err == nil {
//		result.Liked = true
//	} else {
//		if err != sql.ErrNoRows {
//			return nil, err
//		}
//	}
//	accountDao := repository.NewAccountRepository(p.session)
//	author, err := accountDao.GetByID(item.AuthorID)
//	if err != nil {
//		return nil, err
//	}
//	result.AuthorName = author.Nickname
//	result.AuthorAvatar = author.Avatar
//
//	productHistoryDao := repository.NewUserProductRepository(p.session)
//	history, err := productHistoryDao.GetProductLastSellHistory(pk)
//	if err != nil && err != sql.ErrNoRows {
//		return nil, err
//	}
//	sales, err := productHistoryDao.CountByPID(pk)
//	if err != nil {
//		return nil, err
//	}
//	result.Sales = sales - 1
//	owner, err := accountDao.GetByID(history.UID)
//	if err != nil {
//		return nil, err
//	}
//	result.OwnerName = owner.Nickname
//	result.OwnerID = history.UID
//	result.OwnerAvatar = owner.Avatar
//	token := "#" + utils.ZeroFill(sales, len(strconv.FormatInt(item.Stock, 10)))
//	result.TokenID = token
//	return &result, nil
//}
//
//func (p productService) QueryByID(pk int64) (product *models.Product, err error) {
//	result, err := repository.NewProductRepository(p.session).QueryByID(pk)
//	if err == sql.ErrNoRows {
//		return nil, nil
//	}
//	return result, err
//}
//
//func (p productService) CreateProduct(product *models.Product, payType enum.PayType) (interface{}, error) {
//	// 获取当前的创作者
//	user, err := repository.NewAccountRepository(p.session).GetByID(product.AuthorID)
//	if err != nil {
//		return nil, err
//	}
//	// 上链
//	var prefix string
//	switch payType {
//	case enum.Alipay:
//		prefix = "103A"
//	case enum.Wechat:
//		prefix = "32FY"
//	}
//	orderNo := prefix + strings.ReplaceAll(uuid.New().String(), "-", "")[:15] + "-1"
//	client := chain.NewChainClient()
//	resp, err := client.NewProduct(product.Image, user.BsnAddress, product.Description, conf.Instance.Server.Domain+"/api/v1/product/create/callback/"+orderNo)
//	if err != nil {
//		return nil, err
//	}
//	if !resp.Success {
//		return nil, errors.New(resp.ErrMsg)
//	}
//	product.TxID = resp.Data.TxId
//	product.TokenID = resp.Data.TokenId
//
//	// todo 暂时在上线之后删除
//	{
//		product.Status = enum.ProductPaid
//	}
//
//	// 创建支付订单
//	payClient := pay.PayerFactory(payType)
//	if payClient == nil {
//		return nil, g.Error("支付类型错误")
//	}
//	payResp, err := payClient.Pay(orderNo, product.Price)
//	if err != nil {
//		return nil, err
//	}
//	product.OrderNo = orderNo
//
//	// 创建作品
//	if err = repository.NewProductRepository(p.session).Create(product); err != nil {
//		return nil, err
//	}
//
//	// todo 删除这里
//	{
//		//if err = p.ActiveProduct(product.ID); err != nil {
//		//	return nil, err
//		//}
//	}
//
//	// 创建一条拥有者记录
//	if err = repository.NewUserProductRepository(p.session).Create(&models.UserProduct{
//		PID:        product.ID,
//		UID:        user.ID,
//		CreateTime: time.Now().Unix(),
//		Times:      1,
//	}); err != nil {
//		return nil, err
//	}
//	return payResp, nil
//}
//
//func (p productService) Count() (int64, error) {
//	return repository.NewProductRepository(p.session).Count()
//}
//
//func (p productService) QueryList(uid int64, page, pageSize int, orderBy string) (g.ResponseWriter, error) {
//	limit, offset := pageSize, pageSize*(page-1)
//	if strings.Contains(orderBy, "index") {
//		//orderBy = "create_time desc"
//		// 去查推荐列表里面的
//		recommendDao := repository.NewRecommendRepository(p.session)
//		recommends, err := recommendDao.QueryList(limit, offset)
//		if err == sql.ErrNoRows {
//			return g.EmptyListRespWriter, nil
//		}
//		if err != nil {
//			return nil, err
//		}
//		var ids []int64
//		for _, recommend := range recommends {
//			ids = append(ids, recommend.ID)
//		}
//		// 根据推荐的id去查询详情
//		//result, err := repository.NewProductRepository(p.session).QueryByIDList(uid, ids)
//		//if err == sql.ErrNoRows {
//		//	return g.EmptyListRespWriter, nil
//		//}
//		if err != nil {
//			return nil, err
//		}
//		return result, nil
//	} else {
//		result, err := repository.NewProductRepository(p.session).QueryList(uid, limit, offset, orderBy)
//		if err == sql.ErrNoRows {
//			return make([]*entity.ProductList, 0), nil
//		}
//		return result, err
//	}
//}
//
//type ProductLikeService interface {
//	Like(uid int64, pid int64) error
//	CountByPID(pid int64) (int64, error)
//	CountByUID(uid int64) (int64, error)
//}
//
//func NewProductLikeService(session g.Session) ProductLikeService {
//	return &productLikeService{session: session}
//}
//
//type productLikeService struct {
//	session g.Session
//}
//
//func (p productLikeService) CountByUID(uid int64) (int64, error) {
//	return repository.NewProductLikesRepository(p.session).CountByUID(uid)
//}
//
//func (p productLikeService) CountByPID(pid int64) (int64, error) {
//	return repository.NewProductLikesRepository(p.session).CountByPID(pid)
//}
//
//func (p productLikeService) Like(uid int64, pid int64) error {
//	dao := repository.NewProductLikesRepository(p.session)
//	item, err := dao.GetByUIDAndPID(uid, pid)
//	if err == sql.ErrNoRows {
//		// 创建新的
//		model := models.ProductLikes{UID: uid, PID: pid, CreateTime: time.Now().Unix()}
//		return dao.Create(&model)
//	}
//	if err != nil {
//		return err
//	}
//	if item.UID != uid {
//		return g.Error("没有权限")
//	}
//	return dao.DeleteByID(item.ID)
//}
