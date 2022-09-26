package banner

import (
	"github.com/eatmoreapple/regia"
	"lihood/g"
)

func NewRouter() *regia.BluePrint {
	app := regia.NewBluePrint()
	app.Use(g.JWTRequired())
	{
		controller := newBannerController()
		app.GET("/list", controller.bannerList())
	}
	return app
}
