package main

import (
	"log"

	"github.com/eatmoreapple/regia"

	v1 "lihood/api/v1"
	"lihood/boot"
	"lihood/conf"
)

var engine *regia.Engine

func main() {
	log.SetFlags(log.Lshortfile)
	log.Fatal(engine.Run(conf.Instance.Server.Addr))
}

func init() {
	boot.MustSetup()
	engine = regia.New()
	engine.NotFoundHandle = func(context *regia.Context) {
		context.String("success")
	}
	engine.Include("/api/v1", v1.Router())
	engine.Static("/static/", "./h5")
}
