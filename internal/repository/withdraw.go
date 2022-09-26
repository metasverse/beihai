package repository

import (
	"github.com/eatmoreapple/sqlbuilder"
	"lihood/g"
	"lihood/internal/models"
)

type WithdrawRepository interface {
	Create(withdraw *models.Withdraw) error
}

func NewWithdrawRepository(session g.Session) WithdrawRepository {
	return &withdrawRepository{session: session}
}

type withdrawRepository struct {
	session g.Session
}

func (w withdrawRepository) Create(withdraw *models.Withdraw) error {
	builder := sqlbuilder.NewInserter("?")
	builder.Table(withdraw.TableName())
	builder.Fields("uid", "bank_id", "amount", "status", "create_time", "withdraw_time")
	builder.Values(withdraw.UID, withdraw.BankID, withdraw.Amount, withdraw.Status, withdraw.CreateTime, withdraw.WithdrawTime)
	result, err := w.session.Exec(builder.String(), builder.Args()...)
	if err != nil {
		return err
	}
	withdraw.ID, err = result.LastInsertId()
	return err
}
