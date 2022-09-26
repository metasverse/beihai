package enum

type UpChainStatus int

const (
	UpChainStatusUnchain  UpChainStatus = iota // 未上链
	UpChainStatusChaining                      // 上链中
	UpChainStatusChained                       // 已上链
)
