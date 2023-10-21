package server

import (
	"cardValidator/src/errorValidator"
	"encoding/json"
	"log"
	"net/http"
)

const (
	baseUrl = "https://api.bincodes.com/bin/?format=json&"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	var e errorValidator.ErrorValidator
	w.Header().Set("Content-Type", "application/json")
	cardNumber := e.Decode(r)
	apiKey := e.GetApiKey()
	res := e.SendRequest(baseUrl, apiKey, cardNumber)
	bin := e.UnmarshalBin(res)
	resp := e.MakeJson(bin)
	if e.Err != nil {
		log.Println(e.Message)
		resp.Error = e.Message
	}

	json.NewEncoder(w).Encode(resp)

}
