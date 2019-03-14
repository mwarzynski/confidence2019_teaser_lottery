package handlers

import (
	"crypto/md5"
	"encoding/json"
	"io"
	"net/http"
)

type ResponseLotteryResults struct {
	Winners []string `json:"winners"`
}

func (h *Handlers) LotteryResults(w http.ResponseWriter, r *http.Request) {
	winners := h.service.LotteryResults()
	resp := ResponseLotteryResults{
		Winners: make([]string, 0, 0),
	}
	for _, winner := range winners {
		h := md5.New()
		if _, err := io.WriteString(h, winner); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		h.Sum(nil)
		resp.Winners = append(resp.Winners, string(h.Sum(nil)))
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}
