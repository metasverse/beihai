package product

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"time"

	"github.com/eatmoreapple/regia"
	"github.com/google/uuid"
	"github.com/mozillazg/go-pinyin"

	"lihood/conf"
	"lihood/g"
	"lihood/internal/enum"
	"lihood/internal/models"
	"lihood/internal/repository"
	"lihood/internal/requests"
	"lihood/internal/services"
	"lihood/pkg/chain"
	"lihood/pkg/pay"
)

func newProductController() *productController {
	return &productController{}
}

type productController struct{}

// 作品列表
func (p productController) list() regia.HandleFunc {
	return g.Wrapper(func(context *regia.Context) error {
		pagination := g.NewQueryPagination(context)
		service := services.NewProductService(g.DB)
		uid := g.CurrentUserID(context)
		// 添加排序
		order := context.QueryValue("sort").Text()
		writer, err := service.QueryList(uid, pagination.Page(), pagination.PageSize(), order)
		if err != nil {
			return err
		}
		return writer.Write(context)
	})
}

// 创建作品
func (p productController) create() regia.HandleFunc {
	return g.Wrapper(func(context *regia.Context) error {
		var req requests.ProductRequest
		if err := context.Data(&req); err != nil {
			return g.BadRequest(context, err.Error())
		}
		uid := g.CurrentUserID(context)

		now := time.Now()

		model := models.Product{
			Name:        req.Name,
			Price:       req.Price,
			Image:       req.Image,
			Stock:       int64(req.Count),
			Description: req.Description,
			AuthorID:    uid,
			CreateTime:  now.Unix(),
		}

		switch req.SaleTime {
		case 0:
			model.SaleTime = now.Unix()
		case 1:
			// 明天中午12点
			model.SaleTime = time.Date(now.Year(), now.Month(), now.Day()+1, 12, 0, 0, 0, now.Location()).Unix()
		case 2:
			// 后天中午12点
			model.SaleTime = time.Date(now.Year(), now.Month(), now.Day()+2, 12, 0, 0, 0, now.Location()).Unix()
		}

		args := pinyin.NewArgs()
		items := pinyin.Pinyin(model.Name, args)
		var builder strings.Builder
		for _, item := range items {
			builder.WriteString(string(item[0][0]))
		}

		model.Cname = builder.String()

		// 查找cname是否存在
		count, err := repository.NewProductRepository(g.DB).CountByCname(model.Cname)
		if err != nil {
			return err
		}
		model.Cname = "FX" + model.Cname
		if count > 0 {
			model.Cname = fmt.Sprintf("%s%d", model.Cname, count)
		}
		payType := enum.CloudPay

		// 获取当前用户
		user, err := repository.NewAccountRepository(g.DB).GetByID(uid)
		if err != nil {
			return err
		}

		if user.IDCardNum == "" {
			return g.Error("请先实名认证")
		}
		model.OrderNo = uuid.New().String()
		cb := conf.Instance.Server.Domain + fmt.Sprintf("/api/v1/product/create/callback/%s", model.OrderNo)
		resp, err := pay.HfPay(model.Description, user.IDCardNum, user.Name, "9.99", cb, "1")
		if err != nil {
			return g.Error(err.Error())
		}

		tx, err := g.DB.Begin()
		if err != nil {
			return err
		}
		service := services.NewProductService(g.DB)
		_, err = service.CreateProduct(&model, payType)
		if err != nil {
			tx.Rollback()
			return err
		}
		// 请求勇哥地址
		if err = tx.Commit(); err != nil {
			return err
		}
		// 创建订单
		context.SetHeader("Content-Type", "application/json")
		return context.Write(resp)
	})
}

// 商品详情
func (p productController) detail() regia.HandleFunc {
	return g.Wrapper(func(context *regia.Context) error {
		pk, err := context.Params.Get("pk").Int64()
		if err != nil {
			context.Engine.NotFoundHandle(context)
			return nil
		}
		uid := g.CurrentUserID(context)
		service := services.NewProductService(g.DB)
		model, err := service.ProductDetail(pk, uid)
		if err != nil {
			return err
		}
		if model == nil {
			context.Engine.NotFoundHandle(context)
			return nil
		}
		return model.Write(context)
	})
}

// 用户点赞
func (p productController) like() regia.HandleFunc {
	return g.Wrapper(func(context *regia.Context) error {
		pk, err := context.Params.Get("pk").Int64()
		if err != nil {
			context.Engine.NotFoundHandle(context)
			return nil
		}
		service := services.NewProductService(g.DB)
		uid := g.CurrentUserID(context)
		resp, err := service.ProductLike(pk, uid)
		if err != nil {
			return err
		}
		return resp.Write(context)
	})
}

// 我发布的作品
func (p productController) publishList() regia.HandleFunc {
	return g.Wrapper(func(context *regia.Context) error {
		pagination := g.NewQueryPagination(context)
		uid := g.CurrentUserID(context)
		service := services.NewProductService(g.DB)
		writer, err := service.PublishList(uid, pagination.Page(), pagination.PageSize())
		if err != nil {
			return err
		}
		return writer.Write(context)
	})
}

// 我购买的作品
func (p productController) buyList() regia.HandleFunc {
	return g.Wrapper(func(context *regia.Context) error {
		pagination := g.NewQueryPagination(context)
		uid := g.CurrentUserID(context)
		service := services.NewProductService(g.DB)
		writer, err := service.BuyList(uid, pagination.Page(), pagination.PageSize())
		if err != nil {
			return err
		}
		return writer.Write(context)
	})
}

// 我收藏的列表
func (p productController) likeList() regia.HandleFunc {
	return g.Wrapper(func(context *regia.Context) error {
		pagination := g.NewQueryPagination(context)
		uid := g.CurrentUserID(context)
		service := services.NewProductService(g.DB)
		writer, err := service.LikeList(uid, pagination.Page(), pagination.PageSize())
		if err != nil {
			return err
		}
		return writer.Write(context)
	})
}

// 作品的销售额排行榜
func (p productController) salesRankList() regia.HandleFunc {
	return g.Wrapper(func(context *regia.Context) error {
		//pagination := g.NewQueryPagination(context)
		service := services.NewProductService(g.DB)
		result, err := service.SalesRank(1, 10)
		if err != nil {
			return err
		}
		return result.Write(context)
	})
}

// 查找制定用户的发售的作品
func (p productController) salesHistory() regia.HandleFunc {
	return g.Wrapper(func(context *regia.Context) error {
		uid, err := context.Params.Get("id").Int64()
		if err != nil {
			context.Engine.NotFoundHandle(context)
			return nil
		}
		pagination := g.NewQueryPagination(context)
		service := services.NewProductService(g.DB)
		writer, err := service.PublishList(uid, pagination.Page(), pagination.PageSize())
		if err != nil {
			return err
		}
		return writer.Write(context)
	})
}

func (p productController) productList() regia.HandleFunc {
	return g.Wrapper(func(context *regia.Context) error {
		uid, err := context.Params.Get("id").Int64()
		if err != nil {
			context.Engine.NotFoundHandle(context)
			return nil
		}
		pagination := g.NewQueryPagination(context)
		service := services.NewProductService(g.DB)
		writer, err := service.BuyList(uid, pagination.Page(), pagination.PageSize())
		if err != nil {
			return err
		}
		return writer.Write(context)
	})
}

func (p productController) saleList() regia.HandleFunc {
	return g.Wrapper(func(context *regia.Context) error {
		uid := g.CurrentUserID(context)
		pagination := g.NewQueryPagination(context)
		service := services.NewProductService(g.DB)
		writer, err := service.SaleList(uid, pagination.Page(), pagination.PageSize())
		if err != nil {
			return err
		}
		return writer.Write(context)
	})
}

func (p productController) userSaleList() regia.HandleFunc {
	return g.Wrapper(func(context *regia.Context) error {
		uid := context.Params.Get("id").MustInt64()
		pagination := g.NewQueryPagination(context)
		service := services.NewProductService(g.DB)
		writer, err := service.SaleList(uid, pagination.Page(), pagination.PageSize())
		if err != nil {
			return err
		}
		return writer.Write(context)
	})
}

func (p productController) publicProductCreate() regia.HandleFunc {
	return g.Wrapper(func(context *regia.Context) error {
		var req requests.PublicProductRequest
		if err := context.Data(&req); err != nil {
			return g.BadRequest(context, err.Error())
		}
		if req.Password != "PublicProductRequest" {
			return g.Error("密码错误")
		}
		user, err := repository.NewAccountRepository(g.DB).GetByID(req.Uid)
		if err != nil {
			return err
		}
		if user == nil {
			return g.BadRequest(context, "用户不存在")
		}

		model := models.Product{
			Name:        req.Name,
			Price:       req.Price,
			Image:       req.Image,
			Stock:       int64(req.Count),
			Description: req.Description,
			AuthorID:    user.ID,
			CreateTime:  time.Now().Unix(),
		}

		args := pinyin.NewArgs()
		items := pinyin.Pinyin(model.Name, args)
		var builder strings.Builder
		for _, item := range items {
			builder.WriteString(string(item[0][0]))
		}

		model.Cname = builder.String()

		// 查找cname是否存在
		count, err := repository.NewProductRepository(g.DB).CountByCname(model.Cname)
		if err != nil {
			return err
		}
		if count > 0 {
			model.Cname = fmt.Sprintf("%s%d", model.Cname, count)
		}

		var prefix = "FX"
		if req.IsAirDrop {
			prefix = "KT"
		}
		model.Cname = prefix + model.Cname

		tx, err := g.DB.Begin()
		if err != nil {
			return err
		}
		service := services.NewProductService(g.DB)
		writer, err := service.PublicCreateProduct(uuid.New().String(), &model)
		if err != nil {
			tx.Rollback()
			return err
		}
		if err = tx.Commit(); err != nil {
			return err
		}
		return writer.Write(context)
	})
}

func (p productController) airDrop() regia.HandleFunc {
	type Request struct {
		Pid    int64    `json:"pid"`
		Phones []string `json:"phones"`
	}
	return g.Wrapper(func(context *regia.Context) error {
		var req Request
		if err := context.Data(&req); err != nil {
			return g.BadRequest(context, err.Error())
		}
		userProductRepo := repository.NewUserProductRepository(g.DB)
		userProduct, err := userProductRepo.GetByID(req.Pid)
		if err != nil {
			return err
		}
		product, err := repository.NewProductRepository(g.DB).QueryByID(userProduct.PID)
		if err != nil {
			return err
		}

		// 判断当前作品还能出手几次
		times, err := repository.NewUserProductRepository(g.DB).CountByPID(userProduct.PID)
		if err != nil {
			return err
		}
		leftTimes := product.Stock - times

		if leftTimes <= 0 {
			return g.Error("空投已达上限")
		}
		if leftTimes < int64(len(req.Phones)) {
			return g.Error(fmt.Sprintf("此次空投%d个，超出剩余空投次数%d次", len(req.Phones), leftTimes))
		}

		var accounts []*models.Account
		repo := repository.NewAccountRepository(g.DB)
		for _, phone := range req.Phones {
			user, err := repo.GetByPhone(phone)
			if err == nil && user != nil {
				accounts = append(accounts, user)
			}
		}

		client := chain.NewChainClient()
		for i, account := range accounts {
			// 先创建一条记录
			id := uuid.New().String()

			item := models.UserProduct{
				PID:       product.ID,
				UID:       account.ID,
				Times:     times + int64(i) + 1,
				CID:       id,
				IsAirDrop: true,
				Display:   false,
				CName:     fmt.Sprintf("%s #%d", product.Cname, times+int64(i)+1),
			}

			resp, err := client.NewProduct(id, product.Image, account.BsnAddress, product.Description, g.ChainCallback(id))
			if err != nil {
				// 创建一条用户作品
				log.Println(err)
				item.Reason = err.Error()
			}
			if !resp.Success {
				log.Println(err)
				// 创建一条用户作品
				item.Reason = resp.ErrMsg
			} else {
				item.TxID = resp.Data.TxId
				item.TokenID = resp.Data.TokenId
				item.Hash = resp.Data.Hash
			}

			if err = userProductRepo.Create(&item); err != nil {
				log.Println(err)
			}
		}
		return g.OK(context, "ok")
	})
}

//func (p productController) preOrder() regia.HandleFunc {
//	return g.Wrapper(func(context *regia.Context) error {
//		log.Println("preOrder回调")
//		defer context.String("success")
//		pid := context.Query().Get("key")
//		fromaddr := context.Query().Get("fromaddr")
//		toaddr := context.Query().Get("toaddr")
//		uid := context.Query().Get("uid")
//		// 链上下单
//		log.Println(context.Query())
//		token, err := g.Redis.Get(context2.Background(), pid).Result()
//		if err != nil {
//			log.Println(err)
//			return nil
//		}
//		prepayId := uuid.New().String()
//		params := fmt.Sprintf("?uid=%s&pid=%s&prepay_id=%s", uid, context.QueryValue("pid").Text(), prepayId)
//		client := chain.NewChainClient()
//		orderResp, err := client.PreOrder(pid, token, fromaddr, toaddr, conf.Instance.Server.Domain+"/api/v1/product/transfer"+params)
//		if err != nil {
//			log.Println(err)
//			return nil
//		}
//		if !orderResp.Success {
//			log.Println(orderResp.ErrMsg)
//			return nil
//		}
//		if err = g.Redis.Set(context2.Background(), prepayId, orderResp.Data.PrepayID, time.Hour*24).Err(); err != nil {
//			log.Println(err)
//		}
//		log.Println("preOrder回调成功")
//		return nil
//	})
//}
//
//func (p productController) transfer() regia.HandleFunc {
//	return g.Wrapper(func(context *regia.Context) error {
//		defer context.String("success")
//		log.Println("转移回调")
//		log.Println(context.Query())
//		prepayID := context.QueryValue("prepay_id").Text()
//		pid, _ := context.QueryValue("pid").Int64()
//		uid, _ := context.QueryValue("uid").Int64()
//		client := chain.NewChainClient()
//		token, err := g.Redis.Get(context2.Background(), prepayID).Result()
//		if err != nil {
//			log.Println(err)
//			return nil
//		}
//		transResp, err := client.Transfer(token)
//		if err != nil {
//			log.Println(err)
//			return nil
//		}
//		if !transResp.Success {
//			log.Println(transResp.Errmsg)
//			return nil
//		}
//		// 插入一条拥有者记录
//		history := repository.NewUserProductRepository(g.DB)
//		if err = history.Create(&models.UserProduct{
//			PID:        pid,
//			UID:        uid,
//			CreateTime: time.Now().Unix(),
//		}); err != nil {
//			log.Println(err)
//		}
//		log.Println("转移回调成功")
//		return nil
//	})
//}

func (p productController) airDropCallback() regia.HandleFunc {
	return g.Wrapper(func(context *regia.Context) error {
		defer context.String("success")
		cid := context.Params.Get("cid").Text()
		repo := repository.NewUserProductRepository(g.DB)
		item, err := repo.GetByCID(cid)
		if err != nil {
			log.Println("更新状态失败")
			return nil
		}
		if err = repo.UpdateStatusByID(item.ID, 1); err != nil {
			log.Println("更新状态失败")
			return nil
		}
		data, _ := ioutil.ReadAll(context.Request.Body)
		fmt.Println("airDropCallback", string(data))
		go func() {
			crypted, err := base64.StdEncoding.DecodeString(string(data))
			if err != nil {
				fmt.Println(err)
				return
			}
			body, err := chain.AesDecrypt(crypted, []byte("961b0713933818bcec4b5a06ff341b1b"))
			if err != nil {
				fmt.Println(err)
				return
			}
			var callback chain.Callback
			if err = json.Unmarshal(body, &callback); err != nil {
				fmt.Println(err)
				return
			}
			if callback.Appid != "mssc002" {
				return
			}
			if err = repo.UpdateHashByID(item.ID, callback.Hash); err != nil {
				fmt.Println(err)
				return
			}
		}()
		return nil
	})
}
