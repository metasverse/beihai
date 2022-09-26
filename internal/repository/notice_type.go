package repository

import (
	"context"
	"github.com/eatmoreapple/sqlbuilder"
	"lihood/g"
	"lihood/internal/models"
)

type NoticeTypeRepository interface {
	GetByID(id int64) (*models.NoticeType, error)
	List() ([]*models.NoticeType, error)
}

type noticeTypeRepository struct {
	session g.Session
}

func (n noticeTypeRepository) GetByID(id int64) (*models.NoticeType, error) {
	engine := sqlbuilder.NewSelectEngine[*models.NoticeType]("?")
	engine.From(models.NoticeType{}.TableName())
	engine.Where("id = ?", id)
	engine.Limit(1)
	return engine.One(context.Background())
}

func (n noticeTypeRepository) List() ([]*models.NoticeType, error) {
	engine := sqlbuilder.NewSelectEngine[*models.NoticeType]("?")
	engine.From(models.NoticeType{}.TableName()).Where("display = 1")
	return engine.Session(n.session).List(context.Background())
}

func NewNoticeTypeRepository(session g.Session) NoticeTypeRepository {
	return &noticeTypeRepository{session: session}
}
