package notice

import "github.com/eatmoreapple/regia"

func NewRouter() *regia.BluePrint {
	app := regia.NewBluePrint()
	app.GET("/notice/type", noticeTypeList())
	app.GET("/notice/type/list", noticeByType())
	app.GET("/notice/type/all", allNoticeByType())
	app.GET("/notice", notice())
	app.GET("/notice/detail/:id", noticeDetail())
	return app
}
