package repository

import (
	"github.com/eatmoreapple/sqlbuilder"
	"lihood/g"
	"lihood/internal/enum"
	"lihood/internal/models"
)

type ProductOrderRepository interface {
	QueryProductId(uid int64, limit, offset int) ([]int64, error)
	CountByUID(uid int64) (int64, error)
	CountByPID(pid int64) (int64, error)
	Create(item *models.ProductOrder) error
	QueryByOID(oid string) (*models.ProductOrder, error)
	UpdateStatusByID(id int64, status enum.ProductStatus) error
}

func NewProductOrderRepository(session g.Session) ProductOrderRepository {
	return &productOrderRepository{session: session}
}

type productOrderRepository struct {
	session g.Session
}

func (p productOrderRepository) UpdateStatusByID(id int64, status enum.ProductStatus) error {
	builder := sqlbuilder.NewUpdater("?")
	builder.Table(models.ProductOrder{}.TableName())
	builder.Set("status = ?", status)
	builder.Where("id = ?", id)
	_, err := p.session.Exec(builder.String(), builder.Args()...)
	return err
}

func (p productOrderRepository) QueryByOID(oid string) (*models.ProductOrder, error) {
	builder := sqlbuilder.NewSelect("?")
	builder.From(models.ProductOrder{}.TableName())
	builder.Fields("id", "oid", "pay_type", "pid", "uid", "create_time")
	builder.Where("oid = ?", oid)
	builder.Limit(1)
	row := p.session.QueryRow(builder.String(), builder.Args()...)
	item := &models.ProductOrder{}
	err := row.Scan(&item.ID, &item.OID, &item.PayType, &item.PID, &item.UID, &item.CreateTime)
	return item, err
}

func (p productOrderRepository) Create(item *models.ProductOrder) error {
	builder := sqlbuilder.NewInserter("?")
	builder.Table(item.TableName())
	builder.Fields("oid", "pay_type", "pid", "uid", "create_time")
	builder.Values(item.OID, item.PayType, item.PID, item.UID, item.CreateTime)
	result, err := p.session.Exec(builder.String(), builder.Args()...)
	if err != nil {
		return err
	}
	item.ID, err = result.LastInsertId()
	return err
}

func (p productOrderRepository) CountByPID(pid int64) (int64, error) {
	builder := sqlbuilder.NewSelect("?")
	builder.From(models.ProductOrder{}.TableName())
	builder.Fields("COUNT(*)")
	builder.Where("pid = ?", pid)
	row := p.session.QueryRow(builder.String(), builder.Args()...)
	var count int64
	err := row.Scan(&count)
	return count, err
}

func (p productOrderRepository) QueryProductId(uid int64, limit, offset int) ([]int64, error) {
	builder := sqlbuilder.NewSelect("?")
	builder.From(models.ProductOrder{}.TableName())
	builder.Fields("DISTINCT pid", "id")
	builder.Where("uid = ?", uid).And("status = ?", 1)
	builder.Limit(limit)
	builder.Offset(offset)
	builder.OrderBy("id DESC")
	rows, err := p.session.Query(builder.String(), builder.Args()...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items = make([]int64, 0)
	for rows.Next() {
		var item int64
		err := rows.Scan(&item, new(interface{}))
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

func (p productOrderRepository) CountByUID(uid int64) (int64, error) {
	builder := sqlbuilder.NewSelect("?")
	builder.From(models.ProductOrder{}.TableName())
	builder.Fields("COUNT(DISTINCT pid)")
	builder.Where("uid = ?", uid).And("status = ?", 1)
	row := p.session.QueryRow(builder.String(), builder.Args()...)
	var count int64
	err := row.Scan(&count)
	return count, err
}
