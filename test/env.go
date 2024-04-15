package test

import (
	"context"
	"fmt"
	"github.com/alicebob/miniredis/v2"
	"github.com/go-redis/redismock/v9"
	"github.com/redis/go-redis/v9"
	"goboot/pkg/log"
	"time"
)

func NewLog() *log.Logger {
	return log.NewLog(&log.ZapConf{
		Level:         "debug",
		Prefix:        "",
		Format:        "text",
		Director:      "logs",
		EncodeLevel:   "cap",
		StacktraceKey: "stacktrace",
		MaxAge:        0,
		ShowLine:      true,
		LogInConsole:  true,
	})
}

func MockRedis() (*redis.Client, redismock.ClientMock) {
	return redismock.NewClientMock()
}

func NewMiniRedis() redis.Cmdable {
	// 测试用miniredis
	mr, err := miniredis.Run()
	if err != nil {
		panic(fmt.Errorf("new test redis error: %v", err))
	}
	// 使用miniredis创建client
	client := redis.NewClient(&redis.Options{
		Addr:         mr.Addr(),
		Password:     "",
		DB:           0,
		MaxRetries:   5,
		MinIdleConns: 2,
		TLSConfig:    nil,
	})
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = client.Ping(ctx).Result()
	if err != nil {
		panic(fmt.Sprintf("Redis error: %s", err.Error()))
	}
	fmt.Printf("redis connected, url: %s\n", client.Conn().String())
	return client
}
