package repository

import (
	"github.com/eatmoreapple/sqlbuilder"
	"lihood/g"
	"lihood/internal/models"
)

type ProductLikesRepository interface {
	Create(product *models.ProductLikes) error
	DeleteByID(id int64) error
	GetByID(id int64) (*models.ProductLikes, error)
	GetByUIDAndPID(uid int64, pid int64) (*models.ProductLikes, error)
	CountByPID(pid int64) (int64, error)
	QueryByUID(uid int64, limit, offset int) ([]*models.ProductLikes, error)
	CountByUID(uid int64) (int64, error)
}

func NewProductLikesRepository(session g.Session) ProductLikesRepository {
	return &productLikesRepository{session: session}
}

type productLikesRepository struct {
	session g.Session
}

func (p productLikesRepository) CountByUID(uid int64) (int64, error) {
	builder := sqlbuilder.NewSelect("?")
	builder.From(models.ProductLikes{}.TableName())
	builder.Fields("COUNT(*)")
	builder.Where("uid = ?", uid)
	row := p.session.QueryRow(builder.String(), builder.Args()...)
	var result int64
	err := row.Scan(&result)
	return result, err
}

func (p productLikesRepository) QueryByUID(uid int64, limit, offset int) ([]*models.ProductLikes, error) {
	builder := sqlbuilder.NewSelect("?")
	builder.From(models.ProductLikes{}.TableName())
	builder.Fields("id", "pid", "uid", "create_time")
	builder.Where("uid = ?", uid)
	builder.OrderBy(sqlbuilder.Desc("id"))
	builder.Limit(limit).Offset(offset)
	rows, err := p.session.Query(builder.String(), builder.Args()...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result = make([]*models.ProductLikes, 0)
	for rows.Next() {
		var item models.ProductLikes
		if err = rows.Scan(&item.ID, &item.PID, &item.UID, &item.CreateTime); err != nil {
			return nil, err
		}
		result = append(result, &item)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return result, nil
}

func (p productLikesRepository) CountByPID(pid int64) (int64, error) {
	builder := sqlbuilder.NewSelect("?")
	builder.Fields("count(*)")
	builder.From(models.ProductLikes{}.TableName())
	builder.Where("pid = ?", pid)
	rows := p.session.QueryRow(builder.String(), builder.Args()...)
	var count int64
	err := rows.Scan(&count)
	return count, err
}

func (p productLikesRepository) GetByUIDAndPID(uid int64, pid int64) (*models.ProductLikes, error) {
	builder := sqlbuilder.NewSelect("?")
	builder.Fields("id", "pid", "uid", "create_time")
	builder.From(models.ProductLikes{}.TableName())
	builder.Where("uid = ?", uid)
	builder.And("pid = ?", pid)
	builder.Limit(1)
	rows := p.session.QueryRow(builder.String(), builder.Args()...)
	var item = &models.ProductLikes{}
	if err := rows.Scan(&item.ID, &item.PID, &item.UID, &item.CreateTime); err != nil {
		return nil, err
	}
	return item, nil
}

func (p productLikesRepository) DeleteByID(id int64) error {
	builder := sqlbuilder.NewDeleter("?")
	builder.Table(models.ProductLikes{}.TableName())
	builder.Where("id = ?", id)
	//builder.Limit(1)
	_, err := p.session.Exec(builder.String(), builder.Args()...)
	return err
}

func (p productLikesRepository) GetByID(id int64) (*models.ProductLikes, error) {
	builder := sqlbuilder.NewSelect("?")
	builder.Fields("id", "uid", "pid", "create_time")
	builder.From(models.ProductLikes{}.TableName()).Where("id = ?", id)
	rows := p.session.QueryRow(builder.String(), builder.Args()...)
	var item = &models.ProductLikes{}
	err := rows.Scan(&item.ID, &item.UID, &item.ID, &item.CreateTime)
	if err != nil {
		return nil, err
	}
	return item, nil
}

func (p productLikesRepository) Create(product *models.ProductLikes) error {
	builder := sqlbuilder.NewInserter("?")
	builder.Table(models.ProductLikes{}.TableName())
	builder.Fields("uid", "pid", "create_time")
	builder.Values(product.UID, product.PID, product.CreateTime)
	result, err := p.session.Exec(builder.String(), builder.Args()...)
	if err != nil {
		return err
	}
	product.ID, err = result.LastInsertId()
	return err
}
