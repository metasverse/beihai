package enum

type AccountBankStatus int

const (
	AccountBankStatusOK     AccountBankStatus = 0    // 正常
	AccountBankStatusFrozen AccountBankStatus = iota // 冻结
)
