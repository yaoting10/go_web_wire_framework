package testx

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"github.com/hankmor/gotools/ciphers"
	"github.com/hankmor/gotools/conv"
	"github.com/hankmor/gotools/errs"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	signKey     = "sign"
	tokenKey    = "Authorization"
	tokenPrefix = "Bearer"
	uuidKey     = "uuid"
)

type Param url.Values

func (p Param) Encode(key, iv string) []byte {
	if _, ok := p[uuidKey]; !ok {
		p[uuidKey] = []string{conv.Int64ToStr(time.Now().Unix())}
	}
	// 展平map，url.Values 中的value为[]string，转json后是数组，后台无法接收
	var flatMap = map[string]string{}
	for k, v := range p {
		flatMap[k] = strings.Join(v, ",")
	}
	bs, err := json.Marshal(flatMap)
	errs.Throw(err)
	r, err := ciphers.AES.Encrypt(bs, []byte(key), ciphers.CBC, []byte(iv))
	errs.Throw(err)
	return r
}

var Http = &requester{}

type requester struct {
}

func (r *requester) Get(url string, token string, v Param, headers map[string]string) *Resp {
	key := RequestKey(token)
	iv := key
	bs := v.Encode(key, iv)
	requrl := url + "?" + signKey + "=" + hex.EncodeToString(bs)

	req, err := http.NewRequest(http.MethodGet, requrl, nil)
	req.Header.Set(tokenKey, tokenStr(token))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	errs.Throw(err)
	return NewResp(http.DefaultClient.Do(req))
}

func (r *requester) Post(url string, token string, v Param, headers map[string]string) *Resp {
	key := RequestKey(token)
	bs := v.Encode(key, key)

	encodedParams := signKey + "=" + hex.EncodeToString(bs)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader([]byte(encodedParams)))
	errs.Throw(err)
	req.Header.Set(tokenKey, tokenStr(token))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	return NewResp(http.DefaultClient.Do(req))
}
