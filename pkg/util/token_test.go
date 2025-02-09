package util

import (
	"fmt"
	"testing"
	"time"

	"github.com/hankmor/gotools/ciphers"
	"github.com/hankmor/gotools/conv"
)

func TestProductToken(t *testing.T) {
	userId := 0
	prefix := ciphers.MD5(conv.Int64ToStr(int64(userId)))
	middle := ciphers.MD5(UUID32())
	suffix := ciphers.MD5(conv.Int64ToStr(time.Now().UnixMilli())) // 时间戳

	token := prefix + "-" + middle + "-" + suffix
	fmt.Println(token)
}
