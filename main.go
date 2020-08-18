package main

import (
	"fmt"
	"time"

	"github.com/liu-jianhao/token_bucket/token_bucket"
)

func main() {
	// 初始令牌桶有5个，之后每2秒加一个令牌
	tb := token_bucket.New(2*time.Second, 5)
	fmt.Printf("token_bucket avail %d\n", tb.Available())

	count := 0
	for {
		count++
		time.Sleep(1 * time.Second)

		// 拿一个令牌
		if tb.TryTake(1) {
			fmt.Println(count, ": pass")
		} else {
			fmt.Println(count, ": block")
		}
	}
}
