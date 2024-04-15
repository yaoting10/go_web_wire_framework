package apis

import (
	"fmt"
	"goboot/test/testx"
	"net/url"
)

func login(email string, pwd string) string {
	uri := "/login/sign_in"
	var params = url.Values{}
	params.Set("email", email)
	params.Set("password", pwd)
	var m = make(map[string]any)
	key := testx.ResponseKey(fixToken)
	iv := key
	err := testx.Http.Post(host+uri, fixToken, testx.Param(params), nil).Content().Decode(key, iv).Unmarshal(&m)
	if err != nil {
		fmt.Println("error:", err.Error())
		return ""
	}
	return m["token"].(string)
}
