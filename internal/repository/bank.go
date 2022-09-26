package repository

import (
	"github.com/eatmoreapple/sqlbuilder"
	"lihood/g"
	"lihood/internal/models"
)

type BankRepository interface {
	GetByID(id int64) (*models.Bank, error)
	ALL() ([]*models.Bank, error)
}

func NewBankRepository(session g.Session) BankRepository {
	return &bankRepository{session: session}
}

type bankRepository struct {
	session g.Session
}

func (b bankRepository) GetByID(id int64) (*models.Bank, error) {
	selector := sqlbuilder.NewSelect("?")
	selector.Fields("id", "name").From(models.Bank{}.TableName()).Where("id = ?", id).Limit(1)
	var result models.Bank
	if err := b.session.QueryRow(selector.String(), selector.Args()...).Scan(
		&result.ID,
		&result.Name,
	); err != nil {
		return nil, err
	}
	return &result, nil
}

func (b bankRepository) ALL() ([]*models.Bank, error) {
	selector := sqlbuilder.NewSelect("?")
	selector.Fields("id", "name").From(models.Bank{}.TableName()).OrderBy("`id` desc")
	var result = make([]*models.Bank, 0)
	rows, err := b.session.Query(selector.String(), selector.Args()...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var item models.Bank
		if err := rows.Scan(
			&item.ID,
			&item.Name,
		); err != nil {
			return nil, err
		}
		result = append(result, &item)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return result, nil
}
