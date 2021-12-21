package streams

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/asumsi/api.inlive/internal/models/stream"
	"github.com/asumsi/api.inlive/pkg/api"
)


func List(w http.ResponseWriter, r *http.Request){
	queryParams := r.URL.Query()
	live_string := queryParams.Get("live")
	live := true

	if live_string != ""{
		var err error
		live,err = strconv.ParseBool(live_string)
		if err != nil {
			api.RespondJSON(w, api.Response{Code: http.StatusUnprocessableEntity, Message: "invalid value", Data: ""})
			return
		}
	}

	res, err := stream.GetAll(stream.StreamParams{Live: live})

	if err != nil {
		api.RespondJSON(w, api.Response{Code: http.StatusBadRequest, Message: "Error happened", Data: ""})
		fmt.Println(err)
		return
	}
	
	api.RespondJSON(w, api.Response{Code: http.StatusOK, Message: "List of streams", Data: res})
	return

}