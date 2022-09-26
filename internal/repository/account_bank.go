package repository

import (
	"fmt"
	"github.com/eatmoreapple/sqlbuilder"
	"lihood/g"
	"lihood/internal/entity"
	"lihood/internal/enum"
	"lihood/internal/models"
	"time"
)

type AccountBankRepository interface {
	Create(model *models.AccountBank) error
	GetUserBankAccountList(uid int64, limit, offset int) ([]*entity.AccountBank, error)
	GetByID(id int64) (*models.AccountBank, error)
	UnboundBankAccount(id int64) error
	CountUserBankAccount(uid int64) (int64, error)
}

type accountBankRepository struct {
	session g.Session
}

func (a accountBankRepository) CountUserBankAccount(uid int64) (int64, error) {
	selector := sqlbuilder.NewSelect("?")
	selector.Fields("count(*)").From(models.AccountBank{}.TableName()).Where("uid = ?", uid).And("status = ?", enum.AccountBankStatusOK)
	var count int64
	if err := a.session.QueryRow(selector.String(), selector.Args()...).Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}

func (a accountBankRepository) UnboundBankAccount(id int64) error {
	update := sqlbuilder.NewUpdater("?")
	update.Table(models.AccountBank{}.TableName())
	update.Set("`status` = ?, `update_time` = ?", enum.AccountBankStatusFrozen, time.Now().Unix())
	update.Where("id = ?", id).Limit(1)
	_, err := a.session.Exec(update.String(), update.Args()...)
	return err
}

func (a accountBankRepository) GetByID(id int64) (*models.AccountBank, error) {
	builder := sqlbuilder.NewSelect("?")
	builder.Fields("id", "uid", "bank_id", "name", "bank_name", "bank_num", "status", "create_time").
		From(models.AccountBank{}.TableName()).Where("id = ?", id).Limit(1)
	var result models.AccountBank
	if err := a.session.QueryRow(builder.String(), builder.Args()...).Scan(
		&result.ID,
		&result.UID,
		&result.BankID,
		&result.Name,
		&result.BankName,
		&result.BankNum,
		&result.Status,
		&result.CreateTime,
	); err != nil {
		return nil, err
	}
	return &result, nil
}

func (a accountBankRepository) GetUserBankAccountList(uid int64, limit, offset int) ([]*entity.AccountBank, error) {
	builder := sqlbuilder.NewSelect("?")
	builder.Fields("a.id", "a.uid", "a.bank_id", "a.name", "a.bank_name", "a.bank_num", "a.status", "a.create_time", "b.name")
	builder.From(models.AccountBank{}.TableName(), "a").Where("uid = ?", uid)
	builder.LeftJoin(models.Bank{}.TableName(), "bank_id = b.id", "b")
	builder.And("status = ?", enum.AccountBankStatusOK).Limit(limit).Offset(offset).OrderBy("`id` desc")
	var result = make([]*entity.AccountBank, 0)
	rows, err := a.session.Query(builder.String(), builder.Args()...)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var item entity.AccountBank
		if err := rows.Scan(
			&item.ID,
			&item.UID,
			&item.BankID,
			&item.Name,
			&item.BankName,
			&item.BankNum,
			&item.Status,
			&item.CreateTime,
			&item.BankOfficeName,
		); err != nil {
			return nil, err
		}
		result = append(result, &item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return result, nil
}

func (a accountBankRepository) Create(model *models.AccountBank) error {
	builder := sqlbuilder.NewInserter("?")
	builder.Table(models.AccountBank{}.TableName())
	builder.Fields("uid", "bank_id", "name", "bank_name", "bank_num", "status", "create_time").
		Values(model.UID, model.BankID, model.Name, model.BankName, model.BankNum, model.Status, model.CreateTime)
	fmt.Println(builder.String(), builder.Args())
	result, err := a.session.Exec(builder.String(), builder.Args()...)
	if err != nil {
		return err
	}
	model.ID, err = result.LastInsertId()
	return err
}

func NewAccountBankRepository(session g.Session) AccountBankRepository {
	return accountBankRepository{session: session}
}
