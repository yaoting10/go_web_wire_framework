package twitter

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"goboot/pkg/stringx"
	"io"
	"net/http"
	"strings"

	"github.com/pkg/errors"
)

const (
	auth2AuthorizeUrlFormat = "https://twitter.com/i/oauth2/authorize?response_type=code&client_id=%s&redirect_uri=%s&scope=%s&state=%s&code_challenge=%s&code_challenge_method=plain"
)

var EmptyAccessToken = AccessToken{}

type StateMap interface {
	Put(key string, state string)

	Get(key string) string

	Del(key string) string
}

type LocalStateMap struct {
	m map[string]*string
}

func NewLocalStateMap() LocalStateMap {
	return LocalStateMap{m: make(map[string]*string)}
}

func (l LocalStateMap) Put(key string, state string) {
	l.m[key] = &state
}

func (l LocalStateMap) Get(key string) string {
	r := l.m[key]
	if r == nil {
		return ""
	}
	return *r
}

func (l LocalStateMap) Del(key string) string {
	r := l.m[key]
	if r == nil {
		return ""
	}
	l.m[key] = nil
	return *r
}

type OAuth2AuthApi struct {
	sm StateMap
}

func NewOAuth2AuthApi(sm StateMap) OAuth2AuthApi {
	if sm == nil {
		sm = NewLocalStateMap()
	}
	return OAuth2AuthApi{sm: sm}
}

func (o OAuth2AuthApi) AuthorizeUrl(clientID, redirectUri string, scope ...Scope) string {
	scopes := formatScopes(scope...)
	code_challenge := stringx.Randn(16)
	state := stringx.Randn(16)
	o.sm.Put(state, code_challenge)
	return fmt.Sprintf(auth2AuthorizeUrlFormat, clientID, redirectUri, scopes, state, code_challenge)
}

func (o OAuth2AuthApi) tokenUrl() string {
	return fmt.Sprintf(oauth2ApiUrlFormat, "/oauth2/token")
}

func (o OAuth2AuthApi) revokeTokenUrl() string {
	return fmt.Sprintf(oauth2ApiUrlFormat, "/oauth2/revoke")
}

func (o OAuth2AuthApi) encodeClient(clientId, clientSecret string) string {
	dest := clientId + ":" + clientSecret
	return base64.StdEncoding.EncodeToString([]byte(dest))
}

func (o OAuth2AuthApi) RequestAccessToken(clientId, clientSecret, code, state, redirectUri string) (AccessToken, error) {
	var challengeCode = o.sm.Get(state)
	if challengeCode == "" { // TODO not graceful
		return EmptyAccessToken, errors.Wrapf(ApiError, "invalid state")
	}
	var builder = stringx.NewBuilder()
	builder.WriteString("code=").WriteString(code).WriteString("&")
	builder.WriteString("grant_type=authorization_code").WriteString("&")
	builder.WriteString("redirect_uri=").WriteString(redirectUri).WriteString("&")
	builder.WriteString("code_verifier=").WriteString(challengeCode)
	var body = strings.NewReader(builder.String())
	req, err := http.NewRequest(http.MethodPost, o.tokenUrl(), body)
	if err != nil {
		return EmptyAccessToken, errors.Wrapf(ApiError, "request error: %v", err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Authorization", "Basic "+o.encodeClient(clientId, clientSecret))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return EmptyAccessToken, errors.Wrapf(ApiError, "request error: %v", err)
	}
	bs, err := io.ReadAll(resp.Body)
	if err != nil {
		return EmptyAccessToken, err
	}
	var accessToken AccessToken
	if err := json.Unmarshal(bs, &accessToken); err != nil {
		return EmptyAccessToken, errors.Wrapf(ApiError, "invalid response: %v", string(bs))
	}
	return accessToken, nil
}

func (o OAuth2AuthApi) RefreshAccessToken(clientId, clientSecret, refreshToken string) (AccessToken, error) {
	var builder = stringx.NewBuilder()
	builder.WriteString("refresh_token=").WriteString(refreshToken).WriteString("&")
	builder.WriteString("grant_type=refresh_token").WriteString("&")
	builder.WriteString("client_id=").WriteString(clientId)
	var body = strings.NewReader(builder.String())
	req, err := http.NewRequest(http.MethodPost, o.tokenUrl(), body)
	if err != nil {
		return EmptyAccessToken, errors.Wrapf(ApiError, "request error: %v", err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Authorization", "Basic "+o.encodeClient(clientId, clientSecret))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return EmptyAccessToken, errors.Wrapf(ApiError, "request error: %v", err)
	}
	bs, err := io.ReadAll(resp.Body)
	if err != nil {
		return EmptyAccessToken, err
	}
	var accessToken AccessToken
	if err := json.Unmarshal(bs, &accessToken); err != nil {
		return EmptyAccessToken, errors.Wrapf(ApiError, "invalid response: %v", string(bs))
	}
	return accessToken, nil
}

func (o OAuth2AuthApi) RevokeAccessToken(clientId, clientSecret string, token string) error {
	var builder = stringx.NewBuilder()
	builder.WriteString("token=").WriteString(token).WriteString("&")
	builder.WriteString("token_type_hint=access_token")
	var body = strings.NewReader(builder.String())
	req, err := http.NewRequest(http.MethodPost, o.revokeTokenUrl(), body)
	if err != nil {
		return errors.Wrapf(ApiError, "request error: %v", err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Authorization", "Basic "+o.encodeClient(clientId, clientSecret))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.Wrapf(ApiError, "request error: %v", err)
	}
	bs, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	var revokeResut map[string]bool
	if err := json.Unmarshal(bs, &revokeResut); err != nil {
		return errors.Wrapf(ApiError, "invalid response: %v", string(bs))
	}
	if !revokeResut["revoked"] {
		return fmt.Errorf("revoke failed")
	}
	return nil
}

type AccessToken struct {
	TokenType    string `json:"token_type"`
	ExpiresIn    int64  `json:"expires_in"`
	AccessToken  string `json:"access_token"`
	Scope        string `json:"scope"`
	RefreshToken string `json:"refresh_token"`
}
