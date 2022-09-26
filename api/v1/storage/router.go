package storage

import (
	"github.com/eatmoreapple/regia"
	"lihood/g"
)

func NewRouter() *regia.BluePrint {
	app := regia.NewBluePrint()
	app.Use(g.JWTRequired())
	{
		controller := newController()
		// 上传文件
		app.POST("/productUpload", controller.productUpload())
		app.POST("/avatarUpload", controller.avatarUpload())
		app.POST("/idcardUpload", controller.idcardImageUpload())
	}
	return app
}
