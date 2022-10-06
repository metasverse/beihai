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
	query := "SELECT 1 FROM tbl_product_white_list WHERE pid = ? AND uid = ? LIMIT 1"
	return sqlbuilder.NewQueryScanner[bool](w.session, query, pid, uid).One(context.Background())
}
