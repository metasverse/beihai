package order

import (
	"fmt"
	"github.com/eatmoreapple/regia"
	"lihood/g"
	"lihood/internal/enum"
	"lihood/internal/requests"
	"lihood/internal/services"
	"log"
	"time"
)

type controller struct{}

// 下单
func (controller) commit() regia.HandleFunc {
	return g.Wrapper(func(context *regia.Context) error {
		var req requests.OrderRequest
		if err := context.Data(&req); err != nil {
			return g.BadRequest(context, err.Error())
		}
		uid := g.CurrentUserID(context)

		// 查一下当前的作品有没有锁住
		count, err := g.Redis.Exists(context.Request.Context(), fmt.Sprintf("lock:order:%d", req.PID)).Result()
		if err != nil {
			return err
		}
		if count > 0 {
			// 说明被锁住了，直接返回
			return g.Error("当前作品不可购买")
		}
		// 获取锁单时间
		seconds, err := g.Redis.Get(context.Request.Context(), "order_lock_times").Int64()
		if err != nil {
			return err
		}
		tx, err := g.DB.Begin()
		if err != nil {
			return err
		}
		order, err := services.NewOrderService(tx).NewOrder(uid, req.PID, enum.PayType(req.PayType))
		if err != nil {
			tx.Rollback()
			return err
		}
		// 锁住当前作品
		if err := g.Redis.Set(context.Request.Context(), fmt.Sprintf("lock:order:%d", req.PID), 1, time.Second*time.Duration(seconds)).Err(); err != nil {
			tx.Rollback()
			return err
		}
		tx.Commit()
		fmt.Println(order, err)
		return g.OK(context, order)
	})
}

func (c controller) callback() regia.HandleFunc {
	return g.Wrapper(func(context *regia.Context) error {
		defer context.String("success")
		orderId := context.Params.Get("orderId").Text()
		if err := services.NewOrderService(g.DB).OrderCallback(orderId); err != nil {
			log.Println(err)
		}
		return nil
	})
}
