package g

import (
	"github.com/eatmoreapple/regia"
	"lihood/internal/enum"
)

type ResponseWriter interface {
	Write(ctx *regia.Context) error
}

func NewRespWriter[T any](data T) ResponseWriter {
	return &Resp[T]{
		Code:    enum.StatusCodeOK,
		Data:    data,
		Success: true,
	}
}

// Resp is a JSON response envelope.
type Resp[T any] struct {
	Code    enum.StatusCode `json:"code"`
	Msg     string          `json:"message"`
	Data    T               `json:"data"`
	Success bool            `json:"success"`
}

// Write writes the response to the given context.
func (r Resp[T]) Write(ctx *regia.Context) error {
	return ctx.JSON(r)
}

// Error implements the error interface.
func (r Resp[T]) Error() string {
	return r.Msg
}

// NewResp returns a JSON response with the given status code and message.
func NewResp(ctx *regia.Context, success bool, code enum.StatusCode, msg string, data interface{}) error {
	return ctx.JSON(Resp[any]{
		Code:    code,
		Msg:     msg,
		Data:    data,
		Success: success,
	})
}

func Error(msg string) error {
	return &Resp[any]{
		Code:    enum.StatusCodeBadRequest,
		Msg:     msg,
		Data:    nil,
		Success: false,
	}
}

// ForbiddenError is a convenience method for a JSON response with an enum.StatusCodeForbidden status.
func ForbiddenError(msg string) error {
	return &Resp[any]{
		Code:    enum.StatusCodeForbidden,
		Msg:     msg,
		Data:    nil,
		Success: false,
	}
}

// OK is a convenience method for a JSON response with an enum.StatusCodeOK status.
func OK(ctx *regia.Context, data interface{}) error {
	return NewResp(ctx, true, enum.StatusCodeOK, enum.StatusCodeOK.String(), data)
}

// BadRequest is a convenience method for a JSON response with an enum.StatusCodeBadRequest status.
func BadRequest(ctx *regia.Context, msg string) error {
	return NewResp(ctx, false, enum.StatusCodeBadRequest, msg, nil)
}

// ServerError is a convenience method for a JSON response with an enum.StatusCodeServerError status.
func ServerError(ctx *regia.Context) error {
	return NewResp(ctx, false, enum.StatusCodeServerError, enum.StatusCodeServerError.String(), nil)
}

// Fail is a convenience method for a JSON response with an enum.StatusCodeOK status.
func Fail(ctx *regia.Context, msg string) error {
	return NewResp(ctx, false, enum.StatusCodeOK, msg, nil)
}

// Many is a convenience method for a JSON response with an enum.StatusCodeOK status.
func Many(ctx *regia.Context, data interface{}, count int64) error {
	return OK(ctx, regia.Map{"result": data, "count": count})
}

func NotFound(ctx *regia.Context) error {
	return NewResp(ctx, false, enum.StatusCodeNotFound, enum.StatusCodeNotFound.String(), nil)
}

type listData struct {
	Result interface{} `json:"result"`
	Count  int64       `json:"count"`
}

type ListResp struct {
	Code    enum.StatusCode `json:"code"`
	Msg     string          `json:"message"`
	Data    listData        `json:"data"`
	Success bool            `json:"success"`
}

func (l ListResp) Error() string {
	return l.Msg
}

// Write writes the response to the given context.
func (l ListResp) Write(ctx *regia.Context) error {
	return ctx.JSON(l)
}

var EmptyListRespWriter ResponseWriter = &ListResp{
	Code:    enum.StatusCodeOK,
	Success: true,
	Data: listData{
		Result: make([]interface{}, 0),
		Count:  0,
	}}

func NewListResponseWriter(data interface{}, count int64) ResponseWriter {
	return &ListResp{
		Code:    enum.StatusCodeOK,
		Success: true,
		Data: listData{
			Result: data,
			Count:  count,
		}}
}
