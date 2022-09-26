package repository

import (
	"github.com/eatmoreapple/sqlbuilder"
	"lihood/g"
	"lihood/internal/models"
)

type PayTypeRepository interface {
	ALL() ([]*models.PayType, error)
}

func NewPayTypeRepository(session g.Session) PayTypeRepository {
	return &payTypeRepository{session: session}
}

type payTypeRepository struct {
	session g.Session
}

func (p payTypeRepository) ALL() ([]*models.PayType, error) {
	builder := sqlbuilder.NewSelect("?")
	builder.From(models.PayType{}.TableName())
	builder.Fields("id", "name", "status")
	rows, err := p.session.Query(builder.String(), builder.Args()...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var payTypes = make([]*models.PayType, 0)
	for rows.Next() {
		var payType models.PayType
		if err := rows.Scan(&payType.ID, &payType.Name, &payType.Status); err != nil {
			return nil, err
		}
		payTypes = append(payTypes, &payType)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return payTypes, nil
}
