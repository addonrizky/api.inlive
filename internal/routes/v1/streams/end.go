package streams

import (
	"encoding/json"
	"net/http"
	"time"
	"fmt"

	"github.com/asumsi/api.inlive/internal/models/stream"
	"github.com/asumsi/api.inlive/pkg"
	"github.com/asumsi/api.inlive/pkg/api"
	//"github.com/asumsi/api.inlive/pkg/ffmpeg"
	"github.com/gorilla/mux"
	// "gopkg.in/go-playground/validator.v10"
)

// endStream godoc
// @Summary      End stream
// @Description  End stream stop process of send chunk video to dash server using FFMPEG
// @Tags         stream
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Stream ID"
// @Param 		 body body stream.StartStreamRequest false "Body Request"
// @Success      200  {object}  stream.ResponseSwagEndStreamSuccess
// @Failure		 400  {object}	stream.ResponseSwagEndStreamFail
// @Router       /v1/streams/{id}/end [post]
func (controller *Controller) End(w http.ResponseWriter, r *http.Request) {
	var err error
	params := mux.Vars(r)
	id := params["id"]
	var endResult string
	var ok bool

	// check the existence of slug or id first
	result, err := stream.GetBySlugOrId(id)
	if err != nil {
		api.RespondJSON(w, api.Response{Code: http.StatusNotFound, Message: "Can't get the stream data", Data: err})
		return
	}

	// check existence of stream session
	if controller.Sessions[id] == nil {
		api.RespondJSON(w, api.Response{Code: http.StatusBadRequest, Message: "Stream never initiated", Data: nil})
		return
	}

	endStreamRequest := stream.StartStreamRequest{}
	decoder := json.NewDecoder(r.Body)

	//decode body to object startStreamRequest
	if err = decoder.Decode(&endStreamRequest); err != nil {
		api.RespondJSON(w, api.Response{Code: http.StatusBadRequest, Message: "Decode request body fail on end Stream", Data:  err.Error()})
		return
	}

	// validate request body of /end
	err = pkg.ValidateRequest(endStreamRequest)
	if err != nil {
		api.RespondJSON(w, api.Response{Code: http.StatusBadRequest, Message: "Validation Error on end Stream", Data:  err.Error()})
		return
	}

	_, ok = controller.FFmpegs[endStreamRequest.Slug]

	// stream must on running state, so available to be ended
	if !ok {
		api.RespondJSON(w, api.Response{Code: http.StatusOK, Message: "No Streaming to be ended" , Data: ""})
		return
	}

	/*
	if endResult, err = ffmpeg.End(endStreamRequest.Slug); err != nil {
		api.RespondJSON(w, api.Response{Code: http.StatusBadRequest, Message: "Fail to end stream" , Data: err.Error()})
		return
	}
	*/

	session := controller.Sessions[endStreamRequest.Slug]
	//err = session.SCTP().Stop()
	err = session.Close()
	session = nil
	if err != nil {
		fmt.Println("gagal maning")
	}

	delete(controller.FFmpegs, id)
	delete(controller.Sessions, id)

	result.EndDate = pkg.TimePtr(time.Now())

	_, err = result.Update()
	if err != nil {
		api.RespondJSON(w, api.Response{Code: http.StatusInternalServerError, Message: "Failed to save stream end_date", Data: ""})
		return
	}

	api.RespondJSON(w, api.Response{Code: http.StatusOK, Message: "Streaming Stop", Data: endResult})
}