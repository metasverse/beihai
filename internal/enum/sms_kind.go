package enum

type SMSKind int

const (
	SMSLoginMessage SMSKind = iota
	SMSUpdatePhoneMessage
)
