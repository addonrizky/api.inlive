package streams

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/asumsi/api.inlive/internal/models/stream"
	"github.com/asumsi/api.inlive/pkg/api"
)

type response struct {
	stream.Stream
	IsLive bool `json:"is_live"`
}

func (controller *Controller) Get(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	if id == "" {
		api.RespondJSON(w, api.Response{Code: http.StatusBadRequest, Message: "No slug or ID in get request", Data: ""})
	} else {
		result, err := stream.GetBySlugOrId(id)
		if err != nil {
			api.RespondJSON(w, api.Response{Code: http.StatusBadRequest, Message: err.Error(), Data: err})
		}
		_, ok := controller.FFmpegs[id]
		response := response{Stream: result, IsLive: ok}
		api.RespondJSON(w, api.Response{Code: http.StatusOK, Message: http.StatusText((http.StatusOK)), Data: response})
	}

}
