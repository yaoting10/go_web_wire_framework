package testx

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/hankmor/gotools/ciphers"
	"github.com/hankmor/gotools/errs"
	"io"
	"net/http"
)

type body struct {
	raw     []byte
	decoded []byte
	success bool
}

func newBody(bs []byte) *body {
	return &body{raw: bs}
}

func (b *body) Decode(key, iv string) *body {
	if b.success {
		bs, err := hex.DecodeString(string(b.raw))
		errs.Throw(err)
		bs, err = ciphers.AES.Decrypt(bs, []byte(key), ciphers.CBC, []byte(iv))
		errs.Throw(err)
		b.decoded = bs
	} else {
		b.decoded = b.raw
	}
	return b
}

func (b *body) String() string {
	return string(b.decoded)
}

func (b *body) Unmarshal(v any) error {
	if b.success {
		return json.Unmarshal(b.decoded, &v)
	} else {
		return fmt.Errorf(b.String())
	}
}

type Resp struct {
	*http.Response
	body *body
	err  error
}

func NewResp(r *http.Response, err error) *Resp {
	return &Resp{Response: r, err: err}
}

func (r *Resp) Content() *body {
	bs, err := io.ReadAll(r.Body)
	errs.Throw(err)
	if r.StatusCode == http.StatusOK && r.err == nil {
		r.body = newBody(bs)
		r.body.success = true
	} else {
		r.body = newBody([]byte(fmt.Sprintf("request failed, status: %s, content: %s\n", r.Status, string(bs))))
	}
	return r.body
}

func (r *Resp) Error() error {
	return r.err
}
