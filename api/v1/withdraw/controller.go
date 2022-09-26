package withdraw

import (
	"github.com/eatmoreapple/regia"
	"lihood/g"
	"lihood/internal/models"
	"lihood/internal/services"
	"time"
)

type withdrawController struct{}

func (w withdrawController) withdraw() regia.HandleFunc {
	type Request struct {
		Amount int64  `json:"amount" validate:"gt(m=金额不合法,v=0)"`
		Code   string `json:"code" validate:"required(m=短信验证码不能为空)"`
		BankID int64  `json:"bank_id" validate:"gt(m=银行卡不合法,v=0)"`
	}
	return g.Wrapper(func(context *regia.Context) error {
		var request Request
		if err := context.Data(&request); err != nil {
			return g.BadRequest(context, err.Error())
		}
		uid := g.CurrentUserID(context)
		user, err := services.NewAccountService(g.DB).GetByID(uid)
		if err != nil {
			return err
		}
		if err := services.NewSmsService().CheckWithDrawCode(user.Phone, request.Code); err != nil {
			return err
		}
		// 校验短信验证码
		tx, err := g.DB.Begin()
		if err != nil {
			return err
		}
		service := services.NewWithdrawService(tx)
		withdraw := &models.Withdraw{
			UID:        uid,
			Amount:     request.Amount,
			BankID:     request.BankID,
			Status:     0,
			CreateTime: time.Now().Unix(),
		}
		if err := service.Create(withdraw); err != nil {
			tx.Rollback()
			return err
		}
		tx.Commit()
		return g.OK(context, nil)
	})
}

func (w withdrawController) getWithdrawCode() regia.HandleFunc {
	return g.Wrapper(func(context *regia.Context) error {
		uid := g.CurrentUserID(context)
		user, err := services.NewAccountService(g.DB).GetByID(uid)
		if err != nil {
			return err
		}
		if err = services.NewSmsService().SendWithDrawMessage(user.Phone); err != nil {
			return err
		}
		return g.OK(context, nil)
	})
}
