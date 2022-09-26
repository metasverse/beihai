package repository

import (
	"context"
	"github.com/eatmoreapple/sqlbuilder"
	"lihood/g"
	"lihood/internal/models"
)

type NoticeRepository interface {
	ListByType(nid int64, limit, offset int) ([]*models.Notice, error)
	CountByType(nid int64) (int64, error)
	List(limit, offset int) ([]*models.Notice, error)
	GetByID(id int64) (*models.Notice, error)
}

type noticeRepository struct {
	session g.Session
}

func (n noticeRepository) GetByID(id int64) (*models.Notice, error) {
	engine := sqlbuilder.NewSelectEngine[*models.Notice]("?")
	engine.From(models.Notice{}.TableName()).Where("id = ?", id)
	return engine.Session(n.session).One(context.Background())
}

func (n noticeRepository) List(limit, offset int) ([]*models.Notice, error) {
	engine := sqlbuilder.NewSelectEngine[*models.Notice]("?")
	engine.From(models.Notice{}.TableName()).Where("display = 1")
	engine.OrderBy("id DESC")
	engine.Limit(limit).Offset(offset)
	return engine.Session(n.session).List(context.Background())
}

func (n noticeRepository) CountByType(nid int64) (int64, error) {
	engine := sqlbuilder.NewSelectEngine[*models.Notice]("?")
	engine.From(models.Notice{}.TableName()).Where("nid = ? AND display = 1", nid)
	return engine.Session(n.session).Count(context.Background())
}

func (n noticeRepository) ListByType(nid int64, limit, offset int) ([]*models.Notice, error) {
	engine := sqlbuilder.NewSelectEngine[*models.Notice]("?")
	engine.From(models.Notice{}.TableName()).Where("nid = ? AND display = 1", nid)
	engine.OrderBy("id DESC")
	engine.Limit(limit).Offset(offset)
	return engine.Session(n.session).List(context.Background())
}

func NewNoticeRepository(session g.Session) NoticeRepository {
	return &noticeRepository{session: session}
}
