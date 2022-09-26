package product

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/eatmoreapple/regia"
	"github.com/google/uuid"
	"io/ioutil"
	"lihood/g"
	"lihood/internal/enum"
	"lihood/internal/models"
	"lihood/internal/repository"
	"lihood/pkg/chain"
	"log"
	"time"
)

func NewRouter() *regia.BluePrint {
	app := regia.NewBluePrint()

	app.POST("/create/public", newProductController().publicProductCreate())
	app.POST("/airdrop", newProductController().airDrop())
	app.POST("/air-drop-callback/:cid", newProductController().airDropCallback())
	//app.POST("/preorder", newProductController().preOrder())
	//app.POST("/transfer", newProductController().transfer())

	app.POST("/create/callback/:orderId", func(context *regia.Context) {
		// 支付成功回调

		fmt.Println("支付回调了")
		context.String("success")
		orderId := context.Params.Get("orderId").Text()
		// 根据orderID 查找
		productRepo := repository.NewProductRepository(g.DB)
		product, err := productRepo.GetProductByOrderID(orderId)
		if err != nil {
			log.Println("作品查找失败")
			return
		}
		if product.Status == enum.ProductPaid {
			return
		}

		author, err := repository.NewAccountRepository(g.DB).GetByID(product.AuthorID)
		if err != nil {
			log.Println("查找作者失败")
		}
		tx, err := g.DB.Begin()
		if err != nil {
			log.Println("事务开启失败")
			return
		}
		// 将状态改为已支付
		if err = productRepo.ChangeStatusByID(product.ID, enum.ProductPaid); err != nil {
			log.Println("修改支付状态失败")
			tx.Rollback()
			return
		}

		// NOTE: 在订单支付之后将该作品上链
		// 创建一条拥有者记录
		history := &models.UserProduct{
			PID:        product.ID,
			UID:        product.AuthorID,
			CreateTime: time.Now().Unix(),
			Times:      1,
			CID:        uuid.New().String(),
			Display:    true,
			CName:      fmt.Sprintf("%s #00%d", product.Cname, 1),
			SaleTime:   product.SaleTime,
		}

		client := chain.NewChainClient()

		resp, err := client.NewProduct(history.CID, product.Image, author.BsnAddress, product.Description, g.ChainCallback(history.CID))
		if err != nil {
			log.Println("上链失败")
			tx.Rollback()
			return
		}
		if !resp.Success {
			log.Println("上链失败", resp.ErrMsg)
			tx.Rollback()
			return
		}
		history.TxID = resp.Data.TxId
		history.Hash = resp.Data.Hash

		if err = repository.NewUserProductRepository(tx).Create(history); err != nil {
			log.Println("创建失败")
			tx.Rollback()
			return
		}
		tx.Commit()
		// 进行上链
		//data, err := ioutil.ReadAll(context.Request.Body)
		//if err != nil {
		//	fmt.Println(err)
		//	return
		//}
		//crypted, err := base64.StdEncoding.DecodeString(string(data))
		//if err != nil {
		//	fmt.Println(err)
		//	return
		//}
		//body, err := chain.AesDecrypt(crypted, []byte("961b0713933818bcec4b5a06ff341b1b"))
		//if err != nil {
		//	fmt.Println(err)
		//	return
		//}
		//var callback chain.Callback
		//if err = json.Unmarshal(body, &callback); err != nil {
		//	fmt.Println(err)
		//	return
		//}
		//if callback.Appid != "mssc002" {
		//	return
		//}
		//// 将作品的状态更新为已经上链
		//orderId := context.Params.Get("orderId").Text()
		//if err := services.NewProductService(g.DB).ActiveProductChainStatusByOrderNo(orderId, callback.Hash); err != nil {
		//	fmt.Println(err)
		//	return
		//}
	})

	app.POST("/chain/callback/:cid", func(context *regia.Context) {
		fmt.Println("上链回调")
		defer context.String("success")
		cid := context.Params.Get("cid").Text()
		repo := repository.NewUserProductRepository(g.DB)
		item, err := repo.GetByCID(cid)
		if err != nil {
			log.Println("更新状态失败", err.Error())
			return
		}
		if err = repo.UpdateStatusByID(item.ID, 1); err != nil {
			log.Println("更新状态失败", err.Error())
			return
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
	})

	app.Use(g.JWTRequired())
	{
		controller := newProductController()
		app.GET("/list", controller.list())                     // ok
		app.GET("/saleRankList", controller.salesRankList())    // ok
		app.POST("/create", controller.create())                // ok
		app.GET("/detail/:pk", controller.detail())             // ok
		app.POST("/like/:pk", controller.like())                // ok
		app.GET("/publishList", controller.publishList())       // ok
		app.GET("/buyList", controller.buyList())               // ok
		app.GET("/likeList", controller.likeList())             // ok
		app.GET("/saleList", controller.saleList())             // ok
		app.GET("/userSaleList/:id", controller.userSaleList()) // ok
		// 用户发售的作品
		app.GET("/salesHistory/:id", controller.salesHistory()) // ok
		app.GET("/productList/:id", controller.productList())   // ok
	}
	return app
}
