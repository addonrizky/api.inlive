package streams

import (
	"encoding/json"
	"net/http"

	"github.com/asumsi/api.inlive/internal/models/stream"
	"github.com/asumsi/api.inlive/pkg/api"
)

func Create(w http.ResponseWriter, r *http.Request) {
	var newStream stream.Stream
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&newStream); err == nil {
		result, err := stream.Create(newStream)
		if err == nil {
			api.RespondJSON(w, api.Response{Code: http.StatusOK, Message: http.StatusText((http.StatusOK)), Data: result})
		} else {
			api.RespondJSON(w, api.Response{Code: http.StatusBadRequest, Message: err.Error(), Data: err})
		}
	} else {
		api.RespondJSON(w, api.Response{Code: http.StatusBadRequest, Message: err.Error(), Data: err})
	}

}
