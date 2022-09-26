package banner

import (
	"github.com/eatmoreapple/regia"
	"lihood/g"
	"lihood/internal/services"
)

func newBannerController() *indexBannerController {
	return &indexBannerController{}
}

type indexBannerController struct{}

func (i indexBannerController) bannerList() regia.HandleFunc {
	type Resp struct {
		Id    int64  `json:"id"`
		Pid   int64  `json:"pid"`
		Name  string `json:"name"`
		Link  string `json:"link"`
		Image string `json:"image"`
	}
	return g.Wrapper(func(context *regia.Context) error {
		service := services.NewBannerService(g.DB)
		// 固定取前5个
		list, err := service.QueryList(1, 5)
		if err != nil {
			return err
		}
		var resp = make([]*Resp, 0)
		for _, v := range list {
			item := &Resp{
				Id:    v.ID,
				Name:  v.Name,
				Pid:   v.ProductID,
				Link:  v.Link,
				Image: v.Image,
			}
			resp = append(resp, item)
		}
		return g.OK(context, resp)
	})
}
