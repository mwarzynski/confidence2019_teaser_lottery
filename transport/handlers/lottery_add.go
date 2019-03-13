package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/mwarzynski/confidence_web/app"
)

type RequestLotteryAdd struct {
	Name string `json:"accountName"`
}

func (h *Handlers) LotteryAdd(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var req RequestLotteryAdd
	if err := json.Unmarshal(body, &req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if req.Name == "" {
		http.Error(w, "'accountName' must be filled", http.StatusUnprocessableEntity)
		return
	}
	if err := h.service.LotteryAdd(req.Name); err != nil {
		if err == app.ErrNotFound {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}
