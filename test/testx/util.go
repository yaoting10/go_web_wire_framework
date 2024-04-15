package testx

import "strings"

func tokenStr(token string) string {
	return tokenPrefix + " " + token
}

func RequestKey(token string) string {
	return strings.Split(token, "-")[0][0:16]
}

func ResponseKey(token string) string {
	s := strings.Split(token, "-")
	return s[len(s)-1][0:16]
}
