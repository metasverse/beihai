package enum

type Role int

const (
	RoleAdmin  Role = iota // 管理员
	RoleAuthor             // 作者
)
