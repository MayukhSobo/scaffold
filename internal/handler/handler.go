package handler

import (
	"golang-di/pkg/log"
)

type Handler struct {
	logger log.Logger
}

func NewHandler(logger log.Logger) *Handler {
	return &Handler{
		logger: logger,
	}
}

// GetLogger returns the logger for use in derived handlers
func (h *Handler) GetLogger() log.Logger {
	return h.logger
}
