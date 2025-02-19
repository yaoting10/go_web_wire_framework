//go:build wireinject
// +build wireinject

package main

import (
	"goboot/cmd/api/conf"
	"goboot/cmd/api/server"
	"goboot/internal/config"
	"goboot/internal/handler"
	"goboot/internal/repository"
	"goboot/internal/repository/repo"
	"goboot/internal/service"
	"goboot/pkg/aws/s3"
	"goboot/pkg/helper/sid"
	"goboot/pkg/mail"

	"github.com/google/wire"
)

// 全局唯一id
var sidSet = wire.NewSet(sid.NewSid)

var confSet = wire.NewSet(ConfProvider, LoggerProvider, wire.Bind(new(config.Conf), new(*conf.ApiConf)))

// http服务
var serverSet = wire.NewSet(server.NewServerHTTP)

// sys相关集合
var sysHandlerSet = wire.NewSet(handler.NewSysSettingHandler)
var sysSrvSet = wire.NewSet(service.NewAppVersionService, service.NewSysDeviceService, service.NewSysSettingService)
var sysRepoSet = wire.NewSet(repository.NewSysDeviceRepository, repository.NewSysSettingRepository, repository.NewSysVersionRepository)

// user相关集合
var userHandlerSet = wire.NewSet(handler.NewLoginHandler, handler.NewUserHandler)
var userSrvSet = wire.NewSet(service.NewUserService)
var userRepoSet = wire.NewSet(repository.NewUserRepository)

// 集合组
var handlerGroup = wire.NewSet(
	handler.NewHandler, // 顶层handler
	handler.NewDemoApi,
	sysHandlerSet,
	userHandlerSet,
)

var srvGroup = wire.NewSet(
	conf.NewS3Conf,
	s3.NewS3,
	mail.NewSender,
	service.NewService, // 顶层service
	service.NewImageService,
	sysSrvSet,
	userSrvSet,
)

var repoGroup = wire.NewSet(
	repo.NewWDB,
	repo.NewRDB,
	repo.NewRedis,
	repo.NewMongoDB,
	repo.NewRepository,
	sysRepoSet,
	userRepoSet,
)

func newApp() (*server.App, func(), error) {
	panic(wire.Build(
		confSet,
		sidSet,
		repoGroup,
		srvGroup,
		handlerGroup,
		serverSet,
	))
}
