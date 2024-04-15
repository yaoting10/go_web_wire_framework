package apis

import "fmt"

const (
	// local
	host = "http://localhost:8000"
	// dev
	//host = "http://nps.lianduoduo.pro:10304"

	fixToken = "cfcd208495d565ef66e7dff9f98764da-6a127489a99660bf8362766769c0b7a5-437abbaa78f29486b6cfa6fc4bf25f3f"
)

var token string

func getToken() {
	if token == "" {
		//email := "t10@gmail.com"
		email := "hxq1@gmail.com"
		pwd := "123456"
		token = login(email, pwd)
		if token == "" {
			panic(fmt.Errorf("get token failed: %v", token))
		}
	}
}
