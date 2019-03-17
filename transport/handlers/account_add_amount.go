package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/mwarzynski/confidence2019_teaser_lottery/app"
	"github.com/pkg/errors"
)

type RequestAddAmount struct {
	Amount int `json:"amount"`
}

func (h *Handlers) AccountAddAmount(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var req RequestAddAmount
	if err := json.Unmarshal(body, &req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := h.service.AccountAddAmount(chi.URLParam(r, "name"), req.Amount); err != nil {
		cErr := errors.Cause(err)
		if cErr == app.ErrNotFound {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		if cErr == app.ErrInvalidData {
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
