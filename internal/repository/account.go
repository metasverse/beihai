package repository

import (
	"context"
	"fmt"
	"github.com/eatmoreapple/sqlbuilder"
	"lihood/g"
	"lihood/internal/entity"
	"lihood/internal/models"
	"time"
)

type AccountRepository interface {
	Create(account *models.Account) error
	GetByID(id int64) (*models.Account, error)
	GetByPhone(phone string) (*models.Account, error)
	UpdateAccountIdentity(id int64, name, idCard, positiveImage, negativeImage string) error
	UpdatePhoneById(id int64, phone string) error
	UpdateAccountInfo(id int64, nickname string, avatar string, desc string) error
	GetByInviteCode(code string) (*models.Account, error)
	UpdateAmountById(id int64, amount int64) error
	AuthorList(limit, offset int) ([]*entity.AuthorEntity, error)
	AuthorCount() (int64, error)
}

func NewAccountRepository(session g.Session) AccountRepository {
	return &accountRepository{session: session}
}

type accountRepository struct {
	session g.Session
}

func (a accountRepository) AuthorCount() (int64, error) {
	engine := sqlbuilder.NewSelectEngine[*entity.AuthorEntity]("?")
	engine.From(models.Account{}.TableName(), "a")
	engine.LeftJoin(models.Product{}.TableName(), "b.author_id = a.id", "b")
	engine.Fields(
		sqlbuilder.As("a.id", "id"),
		sqlbuilder.As("a.nickname", "nickname"),
		sqlbuilder.As("a.avatar", "avatar"),
		sqlbuilder.Count("b.id", "works_count"),
	)
	engine.OrderBy("works_count desc")
	engine.Having("works_count > 0")
	engine.GroupBy("a.id")
	return sqlbuilder.NewQueryScanner[int64](a.session, fmt.Sprintf("select count(*) from (%s) as t", engine.Select().String())).One(context.Background())
}

func (a accountRepository) AuthorList(limit, offset int) ([]*entity.AuthorEntity, error) {
	engine := sqlbuilder.NewSelectEngine[*entity.AuthorEntity]("?")
	engine.From(models.Account{}.TableName(), "a")
	engine.LeftJoin(models.Product{}.TableName(), "b.author_id = a.id", "b")
	engine.Fields(
		sqlbuilder.As("a.id", "id"),
		sqlbuilder.As("a.nickname", "nickname"),
		sqlbuilder.As("a.avatar", "avatar"),
		sqlbuilder.Count("b.id", "works_count"),
	)
	engine.OrderBy("works_count desc")
	engine.Having("works_count > 0")
	engine.GroupBy("a.id")
	engine.Limit(limit).Offset(offset)
	return engine.Session(a.session).List(context.Background())
}

func (a accountRepository) UpdateAmountById(id int64, amount int64) error {
	builder := sqlbuilder.NewUpdater("?")
	builder.Table(models.Account{}.TableName())
	builder.Set("amount = amount + ?", amount)
	builder.Where("id = ?", id).Limit(1)
	_, err := a.session.Exec(builder.String(), builder.Args()...)
	return err
}

func (a accountRepository) GetByInviteCode(code string) (*models.Account, error) {
	builder := sqlbuilder.NewSelect("?")
	builder.From(models.Account{}.TableName())
	builder.Where("uid = ?", code)
	builder.Limit(1)
	builder.Fields("uid", "nickname", "avatar", "status", "phone", "bsn_address", "name", "id_card_num",
		"id_card_positive_image_url", "id_card_negative_image_url", "role", "create_time")
	row := a.session.QueryRow(builder.String(), builder.Args()...)
	account := &models.Account{}
	err := row.Scan(&account.UID, &account.Nickname, &account.Avatar, &account.Status, &account.Phone,
		&account.BsnAddress, &account.Name, &account.IDCardNum, &account.IdCardPositiveImageUrl,
		&account.IdCardPositiveImageUrl, &account.Role, &account.CreateTime)
	if err != nil {
		return nil, err
	}
	return account, nil
}

func (a accountRepository) UpdateAccountInfo(id int64, nickname string, avatar string, desc string) error {
	builder := sqlbuilder.NewUpdater("?")
	builder.Table(models.Account{}.TableName())
	builder.Set("nickname = ?, avatar = ?, description = ?", nickname, avatar, desc)
	builder.Where("id = ?", id)
	_, err := a.session.Exec(builder.String(), builder.Args()...)
	return err
}

func (a accountRepository) UpdatePhoneById(id int64, phone string) error {
	builder := sqlbuilder.NewUpdater("?")
	builder.Table(models.Account{}.TableName())
	builder.Set("phone = ?", phone)
	builder.Where("id = ?", id)
	_, err := a.session.Exec(builder.String(), builder.Args()...)
	return err
}

func (a accountRepository) UpdateAccountIdentity(id int64, name, idCard, positiveImage, negativeImage string) error {
	update := sqlbuilder.NewUpdater("?")
	update.Table(models.Account{}.TableName())
	update.Set("name = ?, id_card_num = ?, id_card_positive_image_url = ?, id_card_negative_image_url = ?, update_time = ?",
		name, idCard, positiveImage, negativeImage, time.Now().Unix())
	update.Where("id = ?", id)
	_, err := a.session.Exec(update.String(), update.Args()...)
	return err
}

func (a accountRepository) Create(account *models.Account) error {
	insert := sqlbuilder.NewInserter("?")
	insert.Table(models.Account{}.TableName())
	insert.Fields(
		"uid", "pid", "nickname", "avatar", "status", "phone", "bsn_address", "name", "id_card_num",
		"id_card_positive_image_url", "id_card_negative_image_url", "role", "create_time", "bsn_username")
	insert.Values(account.UID, account.PID, account.Nickname, account.Avatar, account.Status, account.Phone, account.BsnAddress, account.Name, account.IDCardNum,
		account.IdCardPositiveImageUrl, account.IdCardNegativeImageUrl, account.Role, account.CreateTime, account.BsnUsername)
	result, err := a.session.Exec(insert.String(), insert.Args()...)
	if err != nil {
		return err
	}
	account.ID, err = result.LastInsertId()
	return err
}

func (a accountRepository) GetByID(id int64) (*models.Account, error) {
	builder := sqlbuilder.NewSelect("?")
	builder.From(models.Account{}.TableName())
	builder.Where("id = ?", id)
	builder.Limit(1)
	return sqlbuilder.BuilderScanner[*models.Account](a.session, builder).One(context.Background())
}

func (a accountRepository) GetByPhone(phone string) (*models.Account, error) {
	selector := sqlbuilder.NewSelect("?")
	selector.Fields("id", "uid", "nickname", "avatar", "status", "phone", "bsn_address", "name", "id_card_num",
		"id_card_positive_image_url", "id_card_negative_image_url", "role", "create_time", "update_time", "del_time")
	selector.From(models.Account{}.TableName())
	selector.Where("phone = ?", phone)
	selector.Limit(1)
	row := a.session.QueryRow(selector.String(), selector.Args()...)
	account := &models.Account{}
	err := row.Scan(&account.ID, &account.UID, &account.Nickname, &account.Avatar, &account.Status, &account.Phone,
		&account.BsnAddress, &account.Name, &account.IDCardNum, &account.IdCardPositiveImageUrl,
		&account.IdCardNegativeImageUrl, &account.Role, &account.CreateTime, &account.UpdateTime, &account.DelTime)
	if err != nil {
		return nil, err
	}
	return account, nil
}
