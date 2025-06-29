package service

import "github.com/MayukhSobo/scaffold/pkg/log"

type Service struct {
	logger log.Logger
}

func NewService(logger log.Logger) *Service {
	return &Service{
		logger: logger,
	}
}
