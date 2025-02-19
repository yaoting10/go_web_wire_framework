//go:build wireinject
// +build wireinject

package main

import (
	"goboot/cmd/job/conf"
	"goboot/cmd/job/server"
	"goboot/cmd/job/task"
	"goboot/internal/config"
	"goboot/internal/repository"
	"goboot/internal/repository/repo"
	"goboot/internal/service"
	"goboot/pkg/aws/s3"
	"goboot/pkg/helper/sid"
	"goboot/pkg/mail"

	"github.com/google/wire"
)

var confSet = wire.NewSet(ConfProvider, LoggerProvider, wire.Bind(new(config.Conf), new(*conf.JobConf)))

// http服务
var serverSet = wire.NewSet(server.NewServerHTTP)

// Job相关集合
var jobSrvSet = wire.NewSet(task.NewDemoTask)

// 全局唯一id
var sidSet = wire.NewSet(sid.NewSid)

// sys相关集合
var sysSrvSet = wire.NewSet(service.NewSysSettingService, service.NewSysDeviceService)
var sysRepoSet = wire.NewSet(repository.NewSysSettingRepository, repository.NewSysDeviceRepository)

// user相关集合
var usrSrvSet = wire.NewSet(service.NewUserService)
var usrRepoSet = wire.NewSet(repository.NewUserRepository)

// 集合组
var srvGroup = wire.NewSet(
	conf.NewS3Conf,
	s3.NewS3,
	mail.NewSender,
	service.NewService, // 顶层service
	service.NewImageService,
	usrSrvSet,
	jobSrvSet,
	sysSrvSet,
)

var repoGroup = wire.NewSet(
	repo.NewWDB,
	repo.NewRDB,
	repo.NewRedis,
	repo.NewMongoDB,
	repo.NewRepository,
	usrRepoSet,
	sysRepoSet,
)

func newApp() (*server.App, func(), error) {
	panic(wire.Build(
		confSet,
		sidSet,
		srvGroup,
		repoGroup,
		serverSet,
	))
}
