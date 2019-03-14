package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/mwarzynski/confidence_web/app"
)

func (h *Handlers) AccountAdd(w http.ResponseWriter, r *http.Request) {
	account, err := h.service.AccountAdd()
	if err != nil {
		if err == app.ErrAlreadyExists {
			http.Error(w, "account (with randomly generated name) already exists, try again", http.StatusConflict)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(account)
}
