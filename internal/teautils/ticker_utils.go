package teautils

import "time"

// 定时运行某个函数
func Every(duration time.Duration, f func(ticker *Ticker)) *Ticker {
	ticker := NewTicker(duration)
	go func() {
		for ticker.Next() {
			f(ticker)
		}
	}()

	return ticker
}
