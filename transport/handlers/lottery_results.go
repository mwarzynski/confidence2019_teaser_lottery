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
	var resp ResponseLotteryResults
	for _, w := range winners {
		h := md5.New()
		io.WriteString(h, w)
		h.Sum(nil)
		resp.Winners = append(resp.Winners, string(h.Sum(nil)))
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
