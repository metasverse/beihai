package repository

import (
	"context"
	"github.com/eatmoreapple/sqlbuilder"
	"lihood/g"
	"lihood/internal/enum"
	"lihood/internal/models"
)

type RecommendRepository interface {
	QueryList(limit, offset int) ([]*models.Recommend, error)
	Count() (int64, error)
}

func NewRecommendRepository(session g.Session) RecommendRepository {
	return &recommendRepository{session: session}
}

type recommendRepository struct {
	session g.Session
}

func (r recommendRepository) Count() (int64, error) {
	query := sqlbuilder.NewSelect("?")
	query.Fields(sqlbuilder.Count("*")).From("tbl_recommend").Where("status = ?", enum.RecommendStatusEnable)
	row := r.session.QueryRow(query.String(), query.Args()...)
	var count int64
	err := row.Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r recommendRepository) QueryList(limit, offset int) ([]*models.Recommend, error) {
	query := sqlbuilder.NewSelect("?")
	query.From(models.Recommend{}.TableName())
	query.Where("`status` = ?", enum.RecommendStatusEnable)
	query.OrderBy(sqlbuilder.Desc("`index`"))
	query.Limit(limit).Offset(offset)
	result, err := sqlbuilder.BuilderScanner[*models.Recommend](r.session, query).List(context.Background())
	return result, err
}
