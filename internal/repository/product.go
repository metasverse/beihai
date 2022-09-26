package repository

import (
	"context"
	"fmt"
	"github.com/eatmoreapple/sqlbuilder"
	"lihood/g"
	"lihood/internal/entity"
	"lihood/internal/enum"
	"lihood/internal/models"
)

type ProductRepository interface {
	QueryList(uid int64, limit, offset int, orderBy string) ([]*entity.ProductList, error)
	Create(product *models.Product) error
	QueryByID(id int64) (*models.Product, error)
	Count() (int64, error)
	PublishCount(uid int64) (int64, error)
	QueryByIDList(uid int64, ids []int64) ([]*entity.ProductList, error)
	SalesRank(limit, offset int) ([]*entity.SalesRank, error)
	ChangeStatusByID(id int64, status enum.ProductStatus) error
	GetProductByOrderID(orderId string) (*models.Product, error)
	UpdateHashByID(id int64, hash string) error
	QueryProductByAuthorID(authorID int64, limit, offset int) ([]*models.Product, error)
	CountByAuthorID(authorID int64) (int64, error)
	QueryProductListByIDs(uid int64, idList []int64) ([]*entity.ProductList, error)
	CountByCname(cname string) (int64, error)
}

func NewProductRepository(session g.Session) ProductRepository {
	return &productRepository{session: session}
}

type productRepository struct {
	session g.Session
}

func (p productRepository) CountByCname(cname string) (int64, error) {
	engine := sqlbuilder.NewSelectEngine[int64]("?")
	engine.Table(models.Product{}.TableName())
	engine.Where("c_name = ?", cname)
	return engine.Session(p.session).Count(context.Background())
}

func (p productRepository) CountByAuthorID(authorID int64) (int64, error) {
	builder := sqlbuilder.NewSelect("?")
	builder.From(models.Product{}.TableName())
	builder.Where("author_id = ?", authorID)
	builder.Fields(sqlbuilder.Count("*"))
	return sqlbuilder.BuilderScanner[int64](p.session, builder).One(context.Background())
}

func (p productRepository) QueryProductListByIDs(uid int64, idList []int64) ([]*entity.ProductList, error) {
	builder := sqlbuilder.NewSelect("?")
	builder.From(models.Product{}.TableName(), "a")
	builder.LeftJoin(models.Account{}.TableName(), "a.author_id = b.id", "b")
	builder.LeftJoin(models.ProductLikes{}.TableName(), "a.id = c.pid", "c")
	builder.Fields("a.id", "a.name", "a.image", "a.author_id", "b.nickname", "a.price",
		sqlbuilder.Count("c.id", "likes"), fmt.Sprintf("EXISTS(SELECT 1 FROM tbl_product_likes WHERE pid = a.id AND uid = %d LIMIT 1) is_liked", uid))
	builder.IN("a.id", sqlbuilder.In[int64](idList)...)
	rows, err := p.session.Query(builder.String(), builder.Args()...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items = make([]*entity.ProductList, 0)
	for rows.Next() {
		var item entity.ProductList
		if err = rows.Scan(&item.ID, &item.Name, &item.Image, &item.AuthorID, &item.AuthorName, &item.Price, &item.Likes, item.Liked); err != nil {
			return nil, err
		}
		items = append(items, &item)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

func (p productRepository) QueryProductByAuthorID(authorID int64, limit, offset int) ([]*models.Product, error) {
	builder := sqlbuilder.NewSelect("?")
	builder.From(models.Product{}.TableName())
	builder.Where("author_id = ?", authorID)
	builder.Limit(limit).Offset(offset)
	builder.OrderBy(sqlbuilder.Desc("id"))
	result, err := sqlbuilder.BuilderScanner[*models.Product](p.session, builder).List(context.Background())
	return result, err
}

func (p productRepository) UpdateHashByID(id int64, hash string) error {
	builder := sqlbuilder.NewUpdater("?")
	builder.Table(models.Product{}.TableName())
	builder.Set("hash = ?", hash)
	builder.Where("id = ?", id)
	_, err := p.session.Exec(builder.String(), builder.Args()...)
	return err
}

func (p productRepository) GetProductByOrderID(orderId string) (*models.Product, error) {
	builder := sqlbuilder.NewSelect("?").From(models.Product{}.TableName()).Where("order_no = ?", orderId).Limit(1)
	return sqlbuilder.BuilderScanner[*models.Product](p.session, builder).One(context.Background())
}

func (p productRepository) ChangeStatusByID(id int64, status enum.ProductStatus) error {
	builder := sqlbuilder.NewUpdater("?")
	builder.Table(models.Product{}.TableName())
	builder.Set("status = ?", status)
	builder.Where("id = ?", id)
	builder.Limit(1)
	_, err := p.session.Exec(builder.String(), builder.Args()...)
	return err
}

func (p productRepository) SalesRank(limit, offset int) ([]*entity.SalesRank, error) {
	builder := sqlbuilder.NewSelect("?")
	builder.From(models.UserProduct{}.TableName(), "a")
	builder.LeftJoin(models.Product{}.TableName(), "a.pid = b.id", "b")
	builder.Fields("a.id", "b.name", "b.image", "b.price * COUNT(a.id) sales")
	builder.Where("a.status = ?", enum.ProductPaid)
	builder.GroupBy("b.pid")
	builder.OrderBy("sales DESC")
	builder.Limit(limit).Offset(offset)
	rows, err := p.session.Query(builder.String(), builder.Args()...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items = make([]*entity.SalesRank, 0)
	for rows.Next() {
		var item = &entity.SalesRank{}
		err := rows.Scan(&item.ID, &item.Name, &item.Image, &item.Amount)
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

func (p productRepository) QueryByIDList(uid int64, ids []int64) ([]*entity.ProductList, error) {
	if len(ids) == 0 {
		return make([]*entity.ProductList, 0), nil
	}
	builder := sqlbuilder.NewSelect("?")
	builder.From(models.Product{}.TableName(), "a")
	builder.LeftJoin(models.Account{}.TableName(), "a.author_id = b.id", "b")
	builder.LeftJoin(models.ProductLikes{}.TableName(), "a.id = c.pid", "c")
	builder.LeftJoin(models.ProductOrder{}.TableName(), "a.id = d.pid", "d")
	builder.Fields("a.id", "a.name", "a.image", "a.author_id", "b.nickname", "a.price", "COUNT(c.id) likes",
		"COUNT(d.id) sales", fmt.Sprintf("EXISTS(SELECT 1 FROM tbl_product_likes WHERE pid = a.id AND uid = %d LIMIT 1) is_liked", uid))
	builder.Where("a.del_time = ?", 0)
	builder.IN("a.id", sqlbuilder.In(ids)...)
	builder.And("a.status = ?", enum.ProductPaid)
	builder.GroupBy("a.id")
	builder.OrderBy("a.index", "sales DESC")
	fmt.Println(builder.String(), builder.Args())
	rows, err := p.session.Query(builder.String(), builder.Args()...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items = make([]*entity.ProductList, 0)
	for rows.Next() {
		var item = &entity.ProductList{}
		err := rows.Scan(&item.ID, &item.Name, &item.Image, &item.AuthorID, &item.AuthorName, &item.Price,
			&item.Likes, &item.Sales, &item.Liked)
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

func (p productRepository) PublishCount(uid int64) (int64, error) {
	builder := sqlbuilder.NewSelect("?")
	builder.Fields("COUNT(*)")
	builder.From(models.Product{}.TableName())
	builder.Where("del_time = ?", 0)
	builder.Where("`status` = ?", enum.ProductPaid)
	builder.Where("author_id = ?", uid)
	rows := p.session.QueryRow(builder.String(), builder.Args()...)
	var count int64
	err := rows.Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (p productRepository) Count() (int64, error) {
	builder := sqlbuilder.NewSelect("?")
	builder.Fields("COUNT(*)")
	builder.From(models.Product{}.TableName())
	builder.Where("del_time = ?", 0)
	builder.Where("`status` = ?", enum.ProductPaid)
	rows := p.session.QueryRow(builder.String(), builder.Args()...)
	var count int64
	err := rows.Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (p productRepository) QueryByID(id int64) (*models.Product, error) {
	builder := sqlbuilder.NewSelect("?")
	builder.From(models.Product{}.TableName())
	builder.Where("id = ?", id)
	builder.Limit(1)
	return sqlbuilder.BuilderScanner[*models.Product](p.session, builder).One(context.Background())
}

func (p productRepository) Create(product *models.Product) error {
	insert := sqlbuilder.NewInserter("?")
	insert.Table(models.Product{}.TableName())
	insert.Fields("name", "image", "author_id", "price", "description", "stock", "create_time", "update_time",
		"order_no", "tx_id", "token_id", "status", "pay_type", "c_name", "sale_time")
	insert.Values(product.Name, product.Image, product.AuthorID, product.Price, product.Description, product.Stock,
		product.CreateTime, product.UpdateTime, product.OrderNo, product.TxID, product.TokenID, product.Status,
		product.PayType, product.Cname, product.SaleTime)
	result, err := p.session.Exec(insert.String(), insert.Args()...)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	product.ID = id
	return nil
}

func (p productRepository) QueryList(uid int64, limit, offset int, orderBy string) ([]*entity.ProductList, error) {
	builder := sqlbuilder.NewSelect("?")
	builder.From(models.Product{}.TableName(), "a")
	builder.LeftJoin(models.Account{}.TableName(), "a.author_id = b.id", "b")
	builder.LeftJoin(models.ProductLikes{}.TableName(), "a.id = c.pid", "c")
	builder.LeftJoin(models.ProductOrder{}.TableName(), "a.id = d.pid", "d")
	builder.Fields("a.id", "a.name", "a.image", "a.author_id", "b.nickname", "a.price", "COUNT(c.id) likes",
		"COUNT(d.id) sales", fmt.Sprintf("EXISTS(SELECT 1 FROM tbl_product_likes WHERE pid = a.id AND uid = %d LIMIT 1) is_liked", uid))
	builder.Where("a.del_time = ?", 0)
	builder.And("a.status = ?", enum.ProductPaid)
	builder.GroupBy("a.id")
	builder.OrderBy(orderBy)
	builder.Limit(limit).Offset(offset)
	fmt.Println(builder.String(), builder.Args())
	rows, err := p.session.Query(builder.String(), builder.Args()...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items = make([]*entity.ProductList, 0)
	for rows.Next() {
		var item = &entity.ProductList{}
		err := rows.Scan(&item.ID, &item.Name, &item.Image, &item.AuthorID, &item.AuthorName, &item.Price,
			&item.Likes, &item.Sales, &item.Liked)
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
