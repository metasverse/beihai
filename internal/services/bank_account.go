package services

import (
	"database/sql"
	"lihood/g"
	"lihood/internal/entity"
	"lihood/internal/enum"
	"lihood/internal/models"
	"lihood/internal/repository"
)

type BankAccountService interface {
	Create(model *models.AccountBank) error
	GetUserBankAccountList(uid int64, page, pageSize int) ([]*entity.AccountBank, error)
	CountUserBankAccount(uid int64) (int64, error)
	UnboundUserBankAccount(uid int64, accountID int64) error
}

func NewBankAccountService(session g.Session) BankAccountService {
	return &bankAccountService{session: session}
}

type bankAccountService struct {
	session g.Session
}

func (b bankAccountService) UnboundUserBankAccount(uid int64, accountID int64) error {
	dao := repository.NewAccountBankRepository(b.session)
	account, err := dao.GetByID(accountID)
	if err == sql.ErrNoRows {
		return g.Error("账户不存在")
	}
	if err != nil {
		return err
	}
	if account.UID != uid {
		return g.ForbiddenError("账户不属于该用户")
	}
	if account.Status != enum.AccountBankStatusOK {
		return g.Error("该账户已处于解绑状态")
	}
	return dao.UnboundBankAccount(accountID)
}

func (b bankAccountService) CountUserBankAccount(uid int64) (int64, error) {
	dao := repository.NewAccountBankRepository(b.session)
	return dao.CountUserBankAccount(uid)
}

func (b bankAccountService) GetUserBankAccountList(uid int64, page, pageSize int) ([]*entity.AccountBank, error) {
	limit, offset := pageSize, (page-1)*pageSize
	result, err := repository.NewAccountBankRepository(b.session).GetUserBankAccountList(uid, limit, offset)
	if err == sql.ErrNoRows {
		return make([]*entity.AccountBank, 0), err
	}
	return result, nil
}

func (b bankAccountService) Create(model *models.AccountBank) error {
	_, err := repository.NewBankRepository(b.session).GetByID(model.BankID)
	if err == sql.ErrNoRows {
		return g.Error("银行不存在")
	}
	if err != nil {
		return err
	}
	return repository.NewAccountBankRepository(b.session).Create(model)
}
