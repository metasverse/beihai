package account

import (
	"github.com/eatmoreapple/regia"
	"lihood/g"
)

func NewRouter() *regia.BluePrint {
	app := regia.NewBluePrint()
	{
		controller := public{}
		app.POST("/public/createUser", controller.createUser())
		app.POST("/public/authentication", controller.authentication())
	}
	{
		controller := newSmsLoginController()
		// 发送短信验证码
		app.POST("/smsLogin/code", controller.sendLoginCode())
		// 短信验证码登录
		app.POST("/smsLogin", controller.login())
	}
	{
		controller := newAccountController()
		// 更新手机号
		app.Use(g.JWTRequired())
		app.POST("/updatePhone", controller.updatePhone())
		// 获取更新手机号的验证码
		app.POST("/updatePhone/code", controller.updatePhoneCode())
		// 获取当前用户详情
		app.GET("/info", controller.info())
		// 获取用户详情
		app.GET("/info/:id", controller.userInfo())
		// 用户实名认证
		app.POST("/authentication", controller.authentication())
		// 更新用户信息
		app.POST("/update", controller.updateAccountInfo())
		// 交易记录
		app.GET("/tradeHistory", controller.tradeHistory())
		// 个人bsn二维码
		app.GET("/qrcode", controller.qrcode())
	}
	app.GET("/artists", artistList())
	return app
}
