package order

import (
	"github.com/eatmoreapple/regia"
	"lihood/g"
)

func NewRouter() *regia.BluePrint {
	app := regia.NewBluePrint()
	co := controller{}
	app.POST("/commit", g.JWTRequired(), co.commit())
	app.POST("/callback/:orderId", co.callback())
	return app
}
