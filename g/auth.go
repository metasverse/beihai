package g

import "github.com/eatmoreapple/regia"

const userID = "user_id"

func CurrentUserID(ctx *regia.Context) int64 {
	value, exist := ctx.GetValue(userID)
	if !exist {
		panic("user_id not found")
	}
	return value.(int64)
}
