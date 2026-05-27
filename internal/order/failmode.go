package order

import "time"

type Kind string

const (
	None    Kind = "none"
	Fail    Kind = "fail"
	Delay   Kind = "delay"
	Panic   Kind = "panic"
	NoReply Kind = "noreply"
)

// FailMode　テストでエラーや遅延をシミュレートするための構造体です。
type FailMode struct {
	Kind  Kind
	Delay time.Duration
}
