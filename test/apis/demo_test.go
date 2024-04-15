package apis

import (
	"fmt"
	"github.com/hankmor/gotools/conv"
	"goboot/internal/consts"
	"goboot/test/testx"
	"net/url"
	"testing"
)

func TestQueryGroupList(t *testing.T) {
	getToken()

	params := url.Values{}
	resp := testx.Http.Get(host+"/demo/i18n", token, testx.Param(params), map[string]string{
		consts.AppChannelKey: conv.IntToStr(consts.AppChannelApk),
		"Version":            "0.0.1",
		"Accept-Language":    "en",
	})
	key := testx.ResponseKey(token)
	iv := key
	r := resp.Content().Decode(key, iv).String()
	fmt.Println(r)
}
