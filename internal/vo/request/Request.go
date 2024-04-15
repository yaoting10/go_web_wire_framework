package request

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/hankmor/gotools/errs"
	"strings"
)

type Request struct {
	Uuid        int64 `json:"uuid" form:"uuid"`
	LoginUserId uint  `json:"loginUserId" form:"loginUserId"`
}

func BindBase(c *gin.Context) Request {
	req := Request{}
	if strings.Index(c.GetHeader("Content-Type"), "multipart/form-data") == 0 {
		errs.Throw(c.MustBindWith(&req, binding.Form))
	} else {
		errs.Throw(c.Bind(&req))
	}
	return req
}

func Bind[T any](c *gin.Context, obj T) T {
	errs.Throw(c.Bind(&obj))
	return obj
}
