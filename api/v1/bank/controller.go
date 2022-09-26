package bank

import (
	"github.com/eatmoreapple/regia"
	"lihood/g"
	"lihood/internal/models"
	"lihood/internal/requests"
	"lihood/internal/services"
	"time"
)

func newBankController() *bankController {
	return &bankController{}
}

type bankController struct{}

// 账户类型列表
func (c bankController) list() regia.HandleFunc {
	return g.Wrapper(func(context *regia.Context) error {
		service := services.NewBankService(g.DB)
		list, err := service.ALL()
		if err != nil {
			return err
		}
		return g.OK(context, list)
	})
}

func newBankCardController() *bankCardController {
	return &bankCardController{}
}

type bankCardController struct{}

// 添加用户银行卡
func (b bankCardController) create() regia.HandleFunc {
	return g.Wrapper(func(context *regia.Context) error {
		var req requests.BankCreateRequest
		if err := context.Data(&req); err != nil {
			return g.BadRequest(context, err.Error())
		}
		uid := g.CurrentUserID(context)
		model := &models.AccountBank{
			UID:        uid,
			BankID:     req.BankID,
			Name:       req.Name,
			BankName:   req.BankName,
			BankNum:    req.BankNum,
			CreateTime: time.Now().Unix(),
		}
		if err := services.NewBankAccountService(g.DB).Create(model); err != nil {
			return err
		}
		return g.OK(context, nil)
	})
}

// 用户银行卡列表
func (b bankCardController) list() regia.HandleFunc {
	return g.Wrapper(func(context *regia.Context) error {
		uid := g.CurrentUserID(context)
		page := g.NewQueryPagination(context)
		service := services.NewBankAccountService(g.DB)
		list, err := service.GetUserBankAccountList(uid, page.Page(), page.PageSize())
		if err != nil {
			return err
		}
		count, err := service.CountUserBankAccount(uid)
		if err != nil {
			return err
		}
		return g.Many(context, list, count)
	})
}

// 用户银行卡解绑
func (b bankCardController) unbound() regia.HandleFunc {
	return g.Wrapper(func(context *regia.Context) error {
		pk, err := context.Params.Get("id").Uint64()
		if err != nil {
			context.Engine.NotFoundHandle(context)
			return nil
		}
		uid := g.CurrentUserID(context)
		service := services.NewBankAccountService(g.DB)
		if err := service.UnboundUserBankAccount(uid, int64(pk)); err != nil {
			return err
		}
		return g.OK(context, nil)
	})
}
