package util

import (
	"github.com/gin-gonic/gin"
	"github.com/hankmor/gotools/errs"
	"goboot/pkg/redisx"
	"strconv"
)

const (
	tokenKey    = "Authorization"
	tokenPrefix = "Bearer"
)

func GetUserId(c *gin.Context, redis redisx.Client) uint {
	userIdStr, err := redis.Get(c, GetToken(c)).Result()
	errs.Throw(err)
	userId, err := strconv.ParseUint(userIdStr, 10, 64)
	errs.Throw(err)
	return uint(userId)
}

func GetRealIp(c *gin.Context) string {
	reqIP := c.Request.Header.Get("X-Forwarded-For")
	if reqIP == "" {
		reqIP = c.ClientIP()
	}
	if reqIP == "::1" {
		reqIP = "127.0.0.1"
	}
	return reqIP
}
