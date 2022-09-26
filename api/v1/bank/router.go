package bank

import (
	"github.com/eatmoreapple/regia"
	"lihood/g"
)

func NewRouter() *regia.BluePrint {
	app := regia.NewBluePrint()
	app.Use(g.JWTRequired())
	{
		controller := newBankController()
		// 银行列表
		app.GET("/list", controller.list())
	}
	{
		controller := newBankCardController()
		// 创建账户
		app.POST("/account/create", controller.create())
		// 账户列表
		app.GET("/account/list", controller.list())
		// 银行卡解绑
		app.POST("/account/unbound/:id", controller.unbound())
	}
	return app
}
