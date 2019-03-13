package handlers

import (
	"github.com/mwarzynski/confidence_web/app"
)

type Handlers struct {
	service *app.Service
}

func New(service *app.Service) *Handlers {
	return &Handlers{
		service: service,
	}
}
