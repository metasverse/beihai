package services

import (
	"database/sql"
	"lihood/g"
	"lihood/internal/models"
	"lihood/internal/repository"
)

type BankService interface {
	ALL() ([]*models.Bank, error)
}

func NewBankService(session g.Session) BankService {
	return &bankService{session: session}
}

type bankService struct {
	session g.Session
}

func (b bankService) ALL() ([]*models.Bank, error) {
	result, err := repository.NewBankRepository(b.session).ALL()
	if err == sql.ErrNoRows {
		return make([]*models.Bank, 0), err
	}
	return result, nil
}
