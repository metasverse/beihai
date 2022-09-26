package services

import (
	"database/sql"
	"lihood/g"
	"lihood/internal/models"
	"lihood/internal/repository"
)

type BannerService interface {
	QueryList(page, pageSize int) ([]*models.Banner, error)
}

func NewBannerService(session g.Session) BannerService {
	return &bannerService{session: session}
}

type bannerService struct {
	session g.Session
}

// QueryList 查询banner列表
func (s bannerService) QueryList(page, pageSize int) (list []*models.Banner, err error) {
	limit, offset := pageSize, (page-1)*pageSize
	list, err = repository.NewBannerRepository(s.session).QueryList(limit, offset)
	if err == sql.ErrNoRows {
		return make([]*models.Banner, 0), nil
	}
	return list, err
}
