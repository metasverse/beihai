package services

import (
	"lihood/g"
	"lihood/internal/enum"
	"lihood/internal/models"
	"lihood/internal/repository"
	"time"
)

type WithdrawService interface {
	Create(withdraw *models.Withdraw) error
}

type withdrawService struct {
	session g.Session
}

func (w withdrawService) Create(withdraw *models.Withdraw) error {
	// 先判断金额的正确性
	if withdraw.Amount <= 0 {
		return g.Error("金额不合法")
	}
	// 再判断余额是否足够
	repo := repository.NewAccountRepository(w.session)
	account, err := repo.GetByID(withdraw.UID)
	if err != nil {
		return err
	}
	if account.Amount < withdraw.Amount {
		return g.Error("余额不足")
	}
	// 查询bank_id是否属于当前用户
	bank, err := repository.NewAccountBankRepository(w.session).GetByID(withdraw.BankID)
	if err != nil {
		return err
	}
	if bank.UID != withdraw.UID {
		return g.Error("请选择自己的银行卡")
	}
	// 更新账户余额
	if err = repo.UpdateAmountById(withdraw.UID, -withdraw.Amount); err != nil {
		return err
	}
	// 创建提现记录
	if err = repository.NewWithdrawRepository(w.session).Create(withdraw); err != nil {
		return err
	}
	// 创建余额变动记录
	if err = repository.NewAccountIncomeRepository(w.session).Create(&models.AccountIncome{
		UID:        withdraw.UID,
		Type:       enum.Expense,
		Amount:     withdraw.Amount,
		Remark:     "提现",
		CreateTime: time.Now().Unix(),
	}); err != nil {
		return err
	}
	return nil
}

func NewWithdrawService(session g.Session) WithdrawService {
	return &withdrawService{session: session}
}
