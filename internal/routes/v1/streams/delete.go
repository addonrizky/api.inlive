package streams

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/asumsi/api.inlive/internal/models/stream"
	"github.com/asumsi/api.inlive/pkg/api"
)

func Delete(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	if id == "" {
		api.RespondJSON(w, api.Response{Code: http.StatusBadRequest, Message: "No slug or ID in delete request", Data: ""})
	} else {
		result, err := stream.GetBySlugOrId(id)
		delete, err := result.Delete()
		if err == nil {
			api.RespondJSON(w, api.Response{Code: http.StatusOK, Message: http.StatusText((http.StatusOK)), Data: delete})
		} else {
			api.RespondJSON(w, api.Response{Code: http.StatusBadRequest, Message: err.Error(), Data: err})
		}
	}

}
