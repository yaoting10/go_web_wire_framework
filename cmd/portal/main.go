package main

import (
	"fmt"
	"goboot/cmd/portal/conf"
	"goboot/internal/config"
	"goboot/pkg/http"
	"goboot/pkg/log"
)

func main() {
	app, cleanup, err := newApp()
	if err != nil {
		panic(err)
	}
	c := app.Conf
	fmt.Printf("server start: %s:%d\n", "http://"+c.Http().Ip, c.Http().Port)

	http.Run(app.Engine, fmt.Sprintf(":%d", c.Http().Port))
	defer cleanup()
}

// 配置provider
func ConfProvider() *conf.PortalConf {
	return conf.NewConfig(func(conf config.Conf) {
	}, "cmd/portal/conf/local.yml")
}

// LoggerProvider provide Logger
func LoggerProvider(c *conf.PortalConf) *log.Logger {
	zap := c.Zap()
	return log.NewLog(&zap)
}
