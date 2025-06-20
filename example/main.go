package main

import (
	"cachex/cache"
	"context"
	"fmt"
	"time"
)

func main() {
	local, _ := cache.NewLocal(1<<20, time.Minute)
	redis := cache.NewRedis("localhost:6379", "", 0)
	c := cache.NewHybridCache(local, redis)

	val, err := c.Get(context.Background(), "user:123", func() (any, error) {
		fmt.Println("fetch from DB")
		return "jack", nil
	}, time.Minute)
	if err != nil {
		panic(err)
	}

	fmt.Println("got:", val)
}

