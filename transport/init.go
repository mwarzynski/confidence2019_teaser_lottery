package transport

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/mwarzynski/confidence2019_teaser_lottery/app"
	"github.com/mwarzynski/confidence2019_teaser_lottery/transport/handlers"
)

func InitRouter(service *app.Service, flag string) http.Handler {
	handler := handlers.New(service, flag)
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Post("/account", handler.AccountAdd)
	r.Post("/account/{name}/amount", handler.AccountAddAmount)
	r.Get("/account/{name}", handler.AccountGet)
	r.Post("/lottery/add", handler.LotteryAdd)
	r.Get("/lottery/results", handler.LotteryResults)

	r.Get("/readiness", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})

	return r
}
