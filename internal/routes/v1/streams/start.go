package streams

import (
	"net/http"
	"encoding/json"
	"fmt"
	"time"

	"github.com/asumsi/api.inlive/internal/models/stream"
	"github.com/asumsi/api.inlive/pkg"
	"github.com/asumsi/api.inlive/pkg/api"
	"github.com/asumsi/api.inlive/pkg/ffmpeg"
	"github.com/gorilla/mux"

	// "gopkg.in/go-playground/validator.v10"
)

type FFMPeg struct {
	pid int
	url string
}

// startStream godoc
// @Summary      Start stream
// @Description  Start stream send chunk video to dash server using FFMPEG
// @Tags         stream
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Stream ID"
// @Param 		 body body stream.StartStreamRequest false "Body Request"
// @Success      200  {object}  stream.ResponseSwagStartStreamSuccess
// @Failure		 200  {object}	stream.ResponseSwagStartStreamFail
// @Router       /v1/streams/{id}/start [post]
func (controller *Controller) Start(w http.ResponseWriter, r *http.Request) {
	var err error

	params := mux.Vars(r)
	id := params["id"]

	// get data stream by slug or id, to next be updated
	streamObj, err := stream.GetBySlugOrId(id)
	if err != nil {
		api.RespondJSON(w, api.Response{Code: http.StatusNotFound, Message: "Stream not found", Data: ""})
		return
	}

	// check existence of stream session
	if controller.Sessions[id] == nil {
		api.RespondJSON(w, api.Response{Code: http.StatusBadRequest, Message: "Stream never initiated", Data: nil})
		return
	}

	startStreamRequest := stream.StartStreamRequest{}
	decoder := json.NewDecoder(r.Body)
	
	//decode body to object startStreamRequest
	if err = decoder.Decode(&startStreamRequest); err != nil {
		api.RespondJSON(w, api.Response{Code: http.StatusBadRequest, Message: "Decode request body fail on start Stream", Data:  err.Error()})
		return
	}
	
	// validate request body of /start
	err = pkg.ValidateRequest(startStreamRequest)
	if err != nil {
		api.RespondJSON(w, api.Response{Code: http.StatusBadRequest, Message: "Validation Error on start Stream", Data:  err.Error()})
		return
	}

	// check if FFMPEG related to slug or id already exist/running
	if _, ok := controller.FFmpegs[startStreamRequest.Slug]; ok {
		api.RespondJSON(w, api.Response{Code: http.StatusOK, Message: "Streaming already started, cant be interfered" , Data: ""})
		return
	}

	// start the streaming
	ffmpegObject, err := startStreaming(startStreamRequest.Slug)

	// if exist any error on run ffmpeg 
	if err != nil {
		api.RespondJSON(w, api.Response{Code: http.StatusInternalServerError, Message: "Stream  encoding failed to run", Data: err})
		return
	}

	// prepare update data on streams
	streamObj.ManifestPath = ffmpegObject["url_stream"].(string)
	streamObj.StartDate = pkg.TimePtr(time.Now())
	
	// do update table streams in database
	_,err = streamObj.Update()
	if err != nil {
		fmt.Println(err)
		api.RespondJSON(w, api.Response{Code: http.StatusInternalServerError, Message: "Failed to save stream manifest", Data: ""})
		return
	}
	
	controller.FFmpegs[startStreamRequest.Slug] = ffmpegObject["pid"].(int)
	fmt.Println(controller.FFmpegs[startStreamRequest.Slug])
	api.RespondJSON(w, api.Response{Code: http.StatusOK, Message: "Stream run successfully", Data: ffmpegObject})

}

// StartStreaming trigger FFMPEG to start running, consuming audiport and videoport produced by rtpforwarder
// also send stream as chunk to dash server
// as output, it will return pid of FFMPEG and manifest url (that can be played on bifrost video player)
func startStreaming(slug string) (map[string]interface{}, error) {

	// channel to receive pid answer from FFMPEG instance
	pid := make(chan int, 1)

	// get path of SDP file to be consumed by FFMPEG instance
	// also construct url stream which can be used as parameter in FFMPEG (-f) to send chunked stream to dash server

	fileSDP := ffmpeg.GetFileSDPPath(slug)
	urlStream := ffmpeg.GetStreamPath(slug)

	// start ffmpeg instance as a goroutine, which give PID as channel result
	go ffmpeg.Execute(fileSDP, urlStream, pid)
	valuePid := <-pid

	result := map[string]interface{}{
		"pid" : valuePid,
		"url_stream" : urlStream,
	}

	//return pid and url_stream of newly created ffmpeg
	return result, nil
}
