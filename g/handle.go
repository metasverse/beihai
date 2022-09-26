package g

import (
	"github.com/eatmoreapple/regia"
	"lihood/internal/enum"
	"lihood/pkg/jwt"
	"strings"
)

type handleFunc func(*regia.Context) error

func Wrapper(f handleFunc) regia.HandleFunc {
	return func(c *regia.Context) {
		if err := f(c); err != nil {
			if v, ok := err.(ResponseWriter); ok {
				_ = v.Write(c)
			} else {
				c.Logger().Error(err)
				_ = ServerError(c)
			}
		}
	}
}

func Recover(context *regia.Context) {
	defer func() {
		if v := recover(); v != nil {
			context.Logger().Error(v)
			_ = ServerError(context)
		}
	}()
	context.Next()
}

func JWTRequired() regia.HandleFunc {
	return func(context *regia.Context) {
		authorization := context.Request.Header.Get("Authorization")
		token := strings.TrimPrefix(authorization, "Bearer ")
		if token == "" {
			context.AbortWithJSON(&Resp[any]{
				Code: enum.StatusCodeUnauthorized,
				Msg:  enum.StatusCodeUnauthorized.String(),
			})
			return
		}
		id, err := jwt.ParseToken(token)
		if err != nil {
			context.AbortWithJSON(&Resp[any]{
				Code: enum.StatusCodeUnauthorized,
				Msg:  enum.StatusCodeUnauthorized.String(),
			})
			return
		}
		context.SetValue(userID, id)
		context.Next()
	}
}
