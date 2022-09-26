package withdraw

import (
	"github.com/eatmoreapple/regia"
	"lihood/g"
)

func NewRouter() *regia.BluePrint {
	app := regia.NewBluePrint()
	app.Use(g.JWTRequired())
	{
		controller := withdrawController{}
		app.POST("/create", controller.withdraw())
		app.POST("/code", controller.getWithdrawCode())
	}
	return app
}
