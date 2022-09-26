package enum

type AccountStatus int

const (
	AccountStatusOK     AccountStatus = 0    // 正常
	AccountStatusFrozen AccountStatus = iota // 冻结
)
