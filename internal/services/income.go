package services

import (
	"database/sql"
	"lihood/g"
	"lihood/internal/enum"
	"lihood/internal/models"
	"lihood/internal/repository"
)

type IncomeService interface {
	QueryByType(uid int64, t enum.IncomeType, page, pageSize int) ([]*models.AccountIncome, error)
	CountByType(uid int64, t enum.IncomeType) (int64, error)
}

func NewIncomeService(session g.Session) IncomeService {
	return &incomeService{session: session}
}

type incomeService struct {
	session g.Session
}

func (i incomeService) CountByType(uid int64, t enum.IncomeType) (int64, error) {
	return repository.NewAccountIncomeRepository(i.session).CountByType(uid, t)
}

func (i incomeService) QueryByType(uid int64, t enum.IncomeType, page, pageSize int) ([]*models.AccountIncome, error) {
	limit, offset := pageSize, (page-1)*pageSize
	result, err := repository.NewAccountIncomeRepository(i.session).QueryByType(uid, t, limit, offset)
	if err == sql.ErrNoRows {
		return make([]*models.AccountIncome, 0), nil
	}
	return result, err
}
