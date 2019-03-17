package handlers

import (
	"github.com/mwarzynski/confidence2019_teaser_lottery/app"
)

type Handlers struct {
	service *app.Service
	flag    string
}

func New(service *app.Service, flag string) *Handlers {
	return &Handlers{
		service: service,
		flag:    flag,
	}
}
