package order

import "github.com/asynkron/protoactor-go/actor"

// compensationFn は補償メッセージを送る関数
type compensationFn func(ctx actor.Context)

// compensations は成功済みステップを逆順に補償するためのスタック
type compensations struct {
	fns []compensationFn
}

func (c *compensations) push(fn compensationFn) {
	c.fns = append(c.fns, fn)
}

// len は未補償ステップ数を返す
func (c *compensations) len() int {
	return len(c.fns)
}

// dispatchNext は逆順で次の補償を1件だけ送信する
func (c *compensations) dispatchNext(ctx actor.Context) {
	i := len(c.fns) - 1
	c.fns[i](ctx)
	c.fns = c.fns[:i]
}
