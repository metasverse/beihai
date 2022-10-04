package repository

import (
	"context"

	"github.com/eatmoreapple/sqlbuilder"

	"lihood/g"
)

type WhiteListRepository interface {
	IsWhiteList(pid, uid int64) (bool, error)
}

func NewWhiteListRepository(session g.Session) WhiteListRepository {
	return whiteListRepository{session: session}
}

type whiteListRepository struct {
	session g.Session
}

func (w whiteListRepository) IsWhiteList(pid, uid int64) (bool, error) {
	builder := sqlbuilder.NewSelect("?")
	builder.From("tbl_product_white_list").Fields("1").Where("pid = ? AND uid = ?", pid, uid).Limit(1)
	return sqlbuilder.BuilderScanner[bool](w.session, builder).One(context.Background())
}
