package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"goboot/pkg/log"
	"net/http"
	"runtime/debug"
)

func Recover(logger *log.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				url := c.Request.URL.Path
				method := c.Request.Method
				logger.Errorf("%-30s: %-5s %s", url, method, err)

				// 控制台打印堆栈信息
				logger.Errorf("%s", string(debug.Stack()))

				c.String(http.StatusInternalServerError, fmt.Sprintf("%v", err))
				c.Abort()
			}
		}()
		c.Next()
	}
}
