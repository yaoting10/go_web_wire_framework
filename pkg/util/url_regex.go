package util

import (
	"regexp"
	"strings"
)

const (
	blurString       = "*"
	doubleBlurString = "**"
)

var UrlMatcher = &urlMatcher{}

type urlMatcher struct {
}

func (um *urlMatcher) Init(urls []string) []*regexp.Regexp {
	var regexps []*regexp.Regexp
	for _, su := range urls {
		su = um.formatRegexString(su)
		r, _ := regexp.Compile(su)
		regexps = append(regexps, r)
	}
	return regexps
}

func (um *urlMatcher) Match(url string, regexps []*regexp.Regexp) bool {
	for _, regex := range regexps {
		if regex.Match([]byte(url)) {
			return true
		}
	}
	return false
}

// formatRegexString 将配置的 skipUrl 转换为正则形式
func (um *urlMatcher) formatRegexString(s string) string {
	s = "^" + s
	if i := strings.Index(s, doubleBlurString); i > 0 {
		// /**/
		s = strings.ReplaceAll(s, "/"+doubleBlurString+"/", "&&&")
		// /**
		s = strings.ReplaceAll(s, "/"+doubleBlurString, "&&")
	}
	if i := strings.Index(s, blurString); i > 0 {
		s = strings.ReplaceAll(s, blurString, "&")
	}
	s = strings.ReplaceAll(s, "&&&", "/(\\w|/?)*")
	s = strings.ReplaceAll(s, "&&", "/(\\w|/)*")
	s = strings.ReplaceAll(s, "&", "(\\w)+")
	return s + "$"
}
