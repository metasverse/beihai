package notice

import (
	"fmt"
	"github.com/eatmoreapple/regia"
	"lihood/g"
	"lihood/internal/models"
	"lihood/internal/repository"
)

// 消息类型列表
func noticeTypeList() regia.HandleFunc {
	return g.Wrapper(func(context *regia.Context) error {
		repo := repository.NewNoticeTypeRepository(g.DB)
		list, err := repo.List()
		if err != nil {
			return err
		}
		return g.OK(context, list)
	})
}

// 具体消息类型消息列表
func noticeByType() regia.HandleFunc {
	return g.Wrapper(func(context *regia.Context) error {
		nid, err := context.QueryValue("nid").Int64()
		if err != nil {
			return err
		}
		pagination := g.NewQueryPagination(context)
		repo := repository.NewNoticeRepository(g.DB)
		fmt.Println(nid)
		list, err := repo.ListByType(nid, pagination.Limit(), pagination.Offset())
		if err != nil {
			return err
		}
		count, err := repo.CountByType(nid)
		if err != nil {
			return err
		}
		return g.Many(context, list, count)
	})
}

// 综合消息列表
func allNoticeByType() regia.HandleFunc {
	return g.Wrapper(func(context *regia.Context) error {
		types, err := repository.NewNoticeTypeRepository(g.DB).List()
		if err != nil {
			return err
		}
		repo := repository.NewNoticeRepository(g.DB)
		var items = make(map[string][]*models.Notice)
		for _, t := range types {
			list, err := repo.ListByType(t.ID, 3, 0)
			if err != nil {
				return err
			}
			items[t.Name] = list
		}
		return g.OK(context, items)
	})
}

// 小喇叭消息列表
func notice() regia.HandleFunc {
	return g.Wrapper(func(context *regia.Context) error {
		pagination := g.NewQueryPagination(context)
		repo := repository.NewNoticeRepository(g.DB)
		list, err := repo.List(pagination.Limit(), pagination.Offset())
		if err != nil {
			return err
		}
		return g.OK(context, list)
	})
}

// 小喇叭详情
func noticeDetail() regia.HandleFunc {
	return g.Wrapper(func(context *regia.Context) error {
		id := context.Params.Get("id").MustInt64()
		repo := repository.NewNoticeRepository(g.DB)
		result, err := repo.GetByID(id)
		if err != nil {
			return err
		}
		return g.OK(context, result)
	})
}
