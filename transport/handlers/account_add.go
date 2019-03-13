package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/mwarzynski/confidence_web/app"
)

func (h *Handlers) AccountAdd(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var account app.Account
	if err := json.Unmarshal(body, &account); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := account.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}
	account.Amounts = make([]int, 0, 0)
	if err := h.service.AccountAdd(account); err != nil {
		if err == app.ErrAlreadyExists {
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}
