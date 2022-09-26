package enum

type StatusCode int

const (
	StatusCodeOK           StatusCode = 0    // ok
	StatusCodeBadRequest   StatusCode = 4000 // 参数错误
	StatusCodeUnauthorized StatusCode = 4001 // 未授权
	StatusCodeForbidden    StatusCode = 4003 // 禁止访问
	StatusCodeNotFound     StatusCode = 4004 // 没有找到
	StatusCodeServerError  StatusCode = 5000 // 服务器开小差了
)
