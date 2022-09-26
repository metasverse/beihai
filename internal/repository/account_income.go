package repository

import (
	"github.com/eatmoreapple/sqlbuilder"
	"lihood/g"
	"lihood/internal/enum"
	"lihood/internal/models"
)

type AccountIncomeRepository interface {
	Create(income *models.AccountIncome) error
	QueryByType(uid int64, t enum.IncomeType, limit, offset int) ([]*models.AccountIncome, error)
	CountByType(uid int64, t enum.IncomeType) (int64, error)
}

func NewAccountIncomeRepository(session g.Session) AccountIncomeRepository {
	return &accountIncomeRepository{session: session}
}

type accountIncomeRepository struct {
	session g.Session
}

func (a accountIncomeRepository) CountByType(uid int64, t enum.IncomeType) (int64, error) {
	builder := sqlbuilder.NewSelect("?")
	builder.Fields("count(*)")
	builder.From(models.AccountIncome{}.TableName())
	builder.Where("uid = ? ", uid)
	builder.And("`type` = ? ", t)
	var count int64
	err := a.session.QueryRow(builder.String(), builder.Args()...).Scan(&count)
	return count, err
}

func (a accountIncomeRepository) QueryByType(uid int64, t enum.IncomeType, limit, offset int) ([]*models.AccountIncome, error) {
	builder := sqlbuilder.NewSelect("?")
	builder.Fields("id", "uid", "type", "amount", "remark", "create_time")
	builder.From(models.AccountIncome{}.TableName())
	builder.Where("uid = ? ", uid)
	builder.And("`type` = ? ", t)
	builder.Limit(limit).Offset(offset)
	builder.OrderBy("id desc")
	rows, err := a.session.Query(builder.String(), builder.Args()...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var incomes = make([]*models.AccountIncome, 0)
	for rows.Next() {
		var income models.AccountIncome
		err := rows.Scan(&income.ID, &income.UID, &income.Type, &income.Amount, &income.Remark, &income.CreateTime)
		if err != nil {
			return nil, err
		}
		incomes = append(incomes, &income)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return incomes, nil
}

func (a accountIncomeRepository) Create(income *models.AccountIncome) error {
	builder := sqlbuilder.NewInserter("?")
	builder.Table(income.TableName())
	builder.Fields("uid", "type", "amount", "remark", "create_time")
	builder.Values(income.UID, income.Type, income.Amount, income.Remark, income.CreateTime)
	result, err := a.session.Exec(builder.String(), builder.Args()...)
	if err != nil {
		return err
	}
	income.ID, err = result.LastInsertId()
	return err
}

func (a accountIncomeRepository) name() {

}
