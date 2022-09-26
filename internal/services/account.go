package services

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/go-redis/redis/v9"
	"lihood/g"
	"lihood/internal/models"
	"lihood/internal/repository"
	"math/rand"
	"time"
)

type AccountService interface {
	GetByID(id int64) (*models.Account, error)
	GetByPhone(phone string) (*models.Account, error)
	UpdatePhone(uid int64, phone string, code string) error
	SendUpdatePhoneCode(uid int64, phone string) error
	Authentication(id int64, name string, idCard string, positiveImage, negativeImage string) error
	UpdateAccountInfo(uid int64, nickname string, avatar string, desc string) error
}

func NewAccountService(session g.Session) AccountService {
	return &accountService{session: session}
}

type accountService struct {
	session g.Session
}

func (a accountService) UpdateAccountInfo(uid int64, nickname string, avatar string, desc string) error {
	return repository.NewAccountRepository(a.session).UpdateAccountInfo(uid, nickname, avatar, desc)
}

func (a accountService) SendUpdatePhoneCode(uid int64, phone string) error {
	// 先限流
	checkKey := a.checkKey(phone)
	ctx := context.Background()
	if err := g.Redis.SetNX(ctx, checkKey, 1, time.Minute*10).Err(); err != nil {
		return err
	}
	times, err := g.Redis.Incr(ctx, checkKey).Result()
	if err != nil {
		return err
	}
	if times > 8 {
		return g.Error("您请求的次数过多，请稍后再试")
	}
	// 发送长度为4的验证码

	key := fmt.Sprintf("sms:%d:update:phone", uid)
	// 先判断是否已经发送过验证码
	value, err := g.Redis.Get(ctx, key).Result()
	if err != nil {
		if err != redis.Nil {
			return err
		}
		value = fmt.Sprintf("%04d", rand.Intn(9999))
	}
	value = "1234"
	if err = g.Redis.Set(ctx, key, value, time.Minute*10).Err(); err != nil {
		return err
	}
	// TODO 发送验证码
	return nil
}

func (a accountService) GetByID(id int64) (*models.Account, error) {
	result, err := repository.NewAccountRepository(a.session).GetByID(id)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (a accountService) GetByPhone(phone string) (*models.Account, error) {
	result, err := repository.NewAccountRepository(a.session).GetByPhone(phone)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (a accountService) UpdatePhone(id int64, phone string, code string) error {
	// 先限流
	checkKey := a.checkKey(phone)
	ctx := context.Background()
	if err := g.Redis.SetNX(ctx, checkKey, 1, time.Minute*10).Err(); err != nil {
		return err
	}
	times, err := g.Redis.Incr(ctx, checkKey).Result()
	if err != nil {
		return err
	}
	if times > 8 {
		return g.Error("您请求的次数过多，请稍后再试")
	}
	// 更新手机号
	key := fmt.Sprintf("sms:%d:update:phone", id)
	value, err := g.Redis.Get(ctx, key).Result()
	if err != nil {
		return err
	}
	if value != code {
		return g.Error("验证码不正确")
	}
	// 更新手机号
	user, err := NewAccountService(g.DB).GetByPhone(phone)
	if err != nil {
		if err == redis.Nil {
			return g.Error("请先获取验证码")
		}
		return err
	}
	if user != nil {
		if user.ID == id {
			return g.Error("请不要重复绑定自己的手机号")
		}
		return g.Error("该手机号已经被绑定")
	}
	return repository.NewAccountRepository(a.session).UpdatePhoneById(id, phone)
}

func (a accountService) checkKey(phone string) string {
	return fmt.Sprintf("ck:%s", phone)
}

// Authentication 用户身份认证
func (a accountService) Authentication(id int64, name string, idCard string, positiveImage, negativeImage string) error {
	return repository.NewAccountRepository(a.session).UpdateAccountIdentity(id, name, idCard, positiveImage, negativeImage)
}
