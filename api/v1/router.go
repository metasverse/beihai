package v1

import (
	"github.com/eatmoreapple/regia"
	"lihood/api/v1/account"
	"lihood/api/v1/bank"
	"lihood/api/v1/banner"
	"lihood/api/v1/notice"
	"lihood/api/v1/order"
	"lihood/api/v1/product"
	"lihood/api/v1/storage"
	"lihood/api/v1/withdraw"
	"lihood/g"
	"lihood/internal/repository"
)

// Router is the router for the API.
func Router() *regia.BluePrint {
	app := regia.NewBluePrint()
	// global middleware
	app.Use(g.Recover)
	// register routers
	app.Include("/account", account.NewRouter())
	app.Include("/banner", banner.NewRouter())
	app.Include("/product", product.NewRouter())
	app.Include("/bank", bank.NewRouter())
	app.Include("/storage", storage.NewRouter())
	app.Include("/withdraw", withdraw.NewRouter())
	app.Include("/order", order.NewRouter())
	app.Include("", notice.NewRouter())
	app.GET("/paytype", func(context *regia.Context) {
		list, err := repository.NewPayTypeRepository(g.DB).ALL()
		if err != nil {
			return
		}
		g.OK(context, list)
	})
	return app
}
