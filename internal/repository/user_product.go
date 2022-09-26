package repository

import (
	"context"
	"fmt"
	"github.com/eatmoreapple/sqlbuilder"
	"lihood/g"
	"lihood/internal/entity"
	"lihood/internal/models"
)

type UserProductRepository interface {
	CountByPID(pid int64) (int64, error)
	Create(item *models.UserProduct) error
	QueryProductListByIDs(uid int64, ids []int64) ([]*entity.ProductList, error)
	GetByID(id int64) (*models.UserProduct, error)
	Count() (int64, error)
	QueryFirstRecordByUID(uid int64, limit, offset int) ([]*models.UserProduct, error)
	QueryFirstRecordCountByUID(uid int64) (int64, error)
	QueryNotFirstRecordByUID(uid int64, limit, offset int) ([]*models.UserProduct, error)
	QueryNotFirstRecordCountByUID(uid int64) (int64, error)
	GetByCID(cid string) (*models.UserProduct, error)
	UpdateStatusByID(id int64, status int) error
	UpdateHashByID(id int64, hash string) error
}

func NewUserProductRepository(session g.Session) UserProductRepository {
	return &userProductRepository{session: session}
}

type userProductRepository struct {
	session g.Session
}

func (p userProductRepository) UpdateHashByID(id int64, hash string) error {
	engine := sqlbuilder.NewUpdateEngine("?")
	engine.Table(models.UserProduct{}.TableName())
	engine.Set("hash = ?", hash)
	engine.Where("id = ?", id)
	engine.Limit(1)
	_, err := engine.Session(p.session).ExecUpdate(context.Background())
	return err
}

func (p userProductRepository) UpdateStatusByID(id int64, status int) error {
	_, err := p.session.Exec("UPDATE tbl_user_product SET `status` = ? WHERE id = ? LIMIT 1", status, id)
	return err
}

func (p userProductRepository) GetByCID(cid string) (*models.UserProduct, error) {
	builder := sqlbuilder.NewSelect("?")
	builder.From(models.UserProduct{}.TableName())
	builder.Where("cid = ?", cid)
	builder.Limit(1)
	return sqlbuilder.BuilderScanner[*models.UserProduct](p.session, builder).One(context.Background())
}

func (p userProductRepository) QueryNotFirstRecordByUID(uid int64, limit, offset int) ([]*models.UserProduct, error) {
	builder := sqlbuilder.NewSelect("?")
	builder.From(models.UserProduct{}.TableName())
	builder.Where("uid = ?", uid).And("times > 1")
	builder.OrderBy(sqlbuilder.Desc("id"))
	builder.Limit(limit).Offset(offset)
	return sqlbuilder.BuilderScanner[*models.UserProduct](p.session, builder).List(context.Background())
}

func (p userProductRepository) QueryNotFirstRecordCountByUID(uid int64) (int64, error) {
	builder := sqlbuilder.NewSelect("?")
	builder.From(models.UserProduct{}.TableName())
	builder.Fields(sqlbuilder.Count("*"))
	builder.Where("uid = ?", uid).And("times > 1")
	return sqlbuilder.BuilderScanner[int64](p.session, builder).One(context.Background())
}

func (p userProductRepository) QueryFirstRecordCountByUID(uid int64) (int64, error) {
	builder := sqlbuilder.NewSelect("?")
	builder.From(models.UserProduct{}.TableName())
	builder.Fields(sqlbuilder.Count("*"))
	builder.Where("uid = ?", uid).And("times = 1")
	return sqlbuilder.BuilderScanner[int64](p.session, builder).One(context.Background())
}

func (p userProductRepository) QueryFirstRecordByUID(uid int64, limit, offset int) ([]*models.UserProduct, error) {
	builder := sqlbuilder.NewSelect("?")
	builder.From(models.UserProduct{}.TableName())
	builder.Where("uid = ?", uid).And("times = 1")
	builder.OrderBy(sqlbuilder.Desc("id"))
	builder.Limit(limit).Offset(offset)
	return sqlbuilder.BuilderScanner[*models.UserProduct](p.session, builder).List(context.Background())
}

func (p userProductRepository) Count() (int64, error) {
	builder := sqlbuilder.NewSelect("?")
	builder.From(models.UserProduct{}.TableName())
	builder.Where("status = ?", 1).And("display = ?", 1)
	builder.Fields(sqlbuilder.Count("*"))
	return sqlbuilder.BuilderScanner[int64](p.session, builder).One(context.Background())
}

func (p userProductRepository) GetByID(id int64) (*models.UserProduct, error) {
	builder := sqlbuilder.NewSelect("?")
	builder.From(models.UserProduct{}.TableName())
	builder.Where("id = ?", id)
	builder.Limit(1)
	result, err := sqlbuilder.BuilderScanner[*models.UserProduct](p.session, builder).One(context.Background())
	return result, err
}

func (p userProductRepository) QueryProductListByIDs(uid int64, ids []int64) ([]*entity.ProductList, error) {
	builder := sqlbuilder.NewSelect("?")
	builder.From(models.UserProduct{}.TableName(), "a")
	builder.LeftJoin(models.Product{}.TableName(), "a.pid = b.id", "b")
	builder.LeftJoin(models.Account{}.TableName(), "a.uid = c.id", "c")
	builder.LeftJoin(models.ProductLikes{}.TableName(), "a.id = d.pid", "d")
	builder.Fields(
		sqlbuilder.As("a.id", "id"),
		sqlbuilder.As("b.name", "name"),
		sqlbuilder.As("b.image", "image"),
		sqlbuilder.As("a.uid", "author_id"),
		sqlbuilder.As(`IFNULL(c.nickname, "")`, "author_name"),
		sqlbuilder.As(`IFNULL(c.avatar, "")`, "author_avatar"),
		sqlbuilder.As("b.price", "price"),
		sqlbuilder.Count("d.id", "likes"),
		sqlbuilder.As(fmt.Sprintf("EXISTS(SELECT 1 FROM tbl_product_likes WHERE pid = a.id AND uid = %d LIMIT 1) ", uid), "is_liked"),
		sqlbuilder.As("a.times", "times"),
		sqlbuilder.As("a.is_air_drop", "is_air_drop"),
	)
	builder.IN("a.id", sqlbuilder.In[int64](ids)...)
	builder.GroupBy("a.id")
	result, err := sqlbuilder.BuilderScanner[*entity.ProductList](p.session, builder).List(context.Background())
	if err != nil {
		return nil, err
	}
	for _, item := range result {
		item.Name = fmt.Sprintf("#%d %s", item.Times, item.Name)
		item.Image = item.Image + "?imageMogr2/thumbnail/!50p"
	}

	newResult := make([]*entity.ProductList, len(result))
	for index, id := range ids {
		for _, item := range result {
			if int64(item.ID) == id {
				newResult[index] = item
				break
			}
		}
	}

	return newResult, nil
}

func (p userProductRepository) Create(item *models.UserProduct) error {
	builder := sqlbuilder.NewInserter("?")
	builder.Table(item.TableName())
	builder.Fields("pid", "uid", "times", "create_time", "tx_id", "token_id", "hash", "cid", "is_air_drop",
		"display", "reason", "c_name", "sale_time")
	builder.Values(item.PID, item.UID, item.Times, item.CreateTime, item.TxID, item.TokenID, item.Hash, item.CID,
		item.IsAirDrop, item.Display, item.Reason, item.CName, item.SaleTime)
	result, err := p.session.Exec(builder.String(), builder.Args()...)
	if err != nil {
		return err
	}
	item.ID, err = result.LastInsertId()
	return err
}

func (p userProductRepository) CountByPID(pid int64) (int64, error) {
	builder := sqlbuilder.NewSelect("?")
	builder.Fields("count(*)")
	//builder.Where("status = 1")
	builder.From(models.UserProduct{}.TableName())
	builder.Where("pid = ?", pid)
	rows := p.session.QueryRow(builder.String(), builder.Args()...)
	var count int64
	err := rows.Scan(&count)
	return count, err
}
