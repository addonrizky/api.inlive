package api

import (
	"encoding/json"
	"net/http"
)

func RespondJSON(w http.ResponseWriter, r Response) {
	w.Header().Set("Content-Type", "application/json")
	var resp []byte
	var err error
	w.WriteHeader(r.Code)
	resp, err = json.Marshal(r)

	if err != nil {
		w.Write(resp)
	} else {
		w.Write(resp)
	}
}
