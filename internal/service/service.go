package service

import (
	"goboot/internal/config"
	"goboot/internal/repository/repo"
	"goboot/pkg/helper/sid"
	"goboot/pkg/log"
)

type Service struct {
	R   *repo.Repository
	L   *log.Logger
	sid *sid.Sid
	C   config.Conf
}

func NewService(repo *repo.Repository, logger *log.Logger, sid *sid.Sid, conf config.Conf) *Service {
	return &Service{
		R:   repo,
		L:   logger,
		sid: sid,
		C:   conf,
	}
}
