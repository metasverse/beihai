package repository

import (
	"context"
	"github.com/eatmoreapple/sqlbuilder"
	"lihood/g"
	"lihood/internal/models"
	"time"
)

type BannerRepository interface {
	QueryList(limit, offset int) ([]*models.Banner, error)
}

type bannerRepository struct {
	session g.Session
}

func (b bannerRepository) QueryList(limit, offset int) ([]*models.Banner, error) {
	engine := sqlbuilder.NewSelectEngine[*models.Banner]("?")
	engine.From(models.Banner{}.TableName())
	engine.Where("status = ?", 1)
	engine.And("del_time = 0")
	now := time.Now().Unix()
	engine.And("start_time <= ?", now).And("end_time >= ?", now)
	engine.Limit(limit).Offset(offset)
	return engine.Session(b.session).List(context.Background())
}

func NewBannerRepository(session g.Session) BannerRepository {
	return &bannerRepository{session: session}
}
