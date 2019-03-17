package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/mwarzynski/confidence2019_teaser_lottery/app"
)

type AccountGetResponse struct {
	Account app.Account `json:"account"`
	Flag    string      `json:"flag,omitempty"`
}

func (h *Handlers) AccountGet(w http.ResponseWriter, r *http.Request) {
	account, flag, err := h.service.AccountGet(chi.URLParam(r, "name"))
	if err != nil {
		if err == app.ErrNotFound {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	response := AccountGetResponse{
		Account: account,
	}
	if flag {
		response.Flag = h.flag
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(response)
}
