package streams

import (
	"encoding/json"
	"net/http"
	//"fmt"

	"github.com/asumsi/api.inlive/internal/models/stream"
	"github.com/asumsi/api.inlive/pkg"
	"github.com/asumsi/api.inlive/pkg/api"
	rtc "github.com/asumsi/api.inlive/pkg/webrtc"
	"github.com/gorilla/mux"
	"github.com/pion/webrtc/v3"
	"strconv"
)

// InitSession godoc
// @Summary      Init session stream
// @Description  Init session send offer session to WebRTC
// @Tags         stream
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Stream ID"
// @Param 		 body body stream.InitSessionRequestSwag false "Body Request"
// @Success      200  {object}  stream.ResponseSwagInitSessSuccess
// @Failure		 400  {object}	stream.ResponseSwagInitSessFail
// @Router       /v1/streams/{id}/init [post]
func (controller *Controller) Init(w http.ResponseWriter, r *http.Request) {
	var err error
	var peerConnection *webrtc.PeerConnection
	var sessionDescription webrtc.SessionDescription
	
	//resultPeerConnection := make(chan *webrtc.PeerConnection, 1)

	params := mux.Vars(r)
	id := params["id"]

	_, err = stream.GetBySlugOrId(id)
	if err != nil {
		api.RespondJSON(w, api.Response{Code:http.StatusNotFound, Message: "Stream not found", Data: ""})
		return
	}

	if controller.Sessions[id] != nil {
		api.RespondJSON(w, api.Response{Code: http.StatusOK, Message: "Stream already initiated", Data: nil})
		return
	}

	initSessionRequest := stream.InitSessionRequest{}
	decoder := json.NewDecoder(r.Body)
	
	//decode body to object initSessionRequest
	if err = decoder.Decode(&initSessionRequest); err != nil {
		api.RespondJSON(w, api.Response{Code: http.StatusBadRequest, Message: "Decode request body fail on init Stream", Data:  err.Error()})
		return
	}
	
	// validate request body of /init
	err = pkg.ValidateRequest(initSessionRequest)
	if err != nil {
		api.RespondJSON(w, api.Response{Code: http.StatusBadRequest, Message: "Validation Error on init Stream", Data:  err.Error()})
		return
	}

	// Session parameter must exist
	if sessionDescription = initSessionRequest.SessionDescription; sessionDescription.SDP == "" {
		api.RespondJSON(w, api.Response{Code: http.StatusBadRequest, Message: "Validation Error on init Stream", Data:  "json session must exist and valid"})
		return
	}
		
	// establish p2p connection
	if peerConnection, err = initSession(controller, initSessionRequest.Slug, sessionDescription); err != nil {
		api.RespondJSON(w, api.Response{Code: http.StatusBadRequest, Message: "Fail to establish P2P connection on init Stream", Data:  err.Error()})
		return
	}

	//go pkg.StartP2P(sessionDescription, resultPeerConnection)
	//peerConnection = <-resultPeerConnection

	// save session in persistent map Session
	controller.Sessions[id] = peerConnection
	controller.Sessions[initSessionRequest.Slug] = peerConnection

	// set answer to be return to client
	answerSdp := peerConnection.LocalDescription()

	api.RespondJSON(w, api.Response{Code: http.StatusOK, Message: "Stream initiated", Data: answerSdp})
	return
}

// InitSession generate answer session base on offer from browser,
// also spawn rtpforwarder instance with randomized udp's video port and audio port
func initSession(controller *Controller, streamSlug string, offer webrtc.SessionDescription) (*webrtc.PeerConnection, error) {
	var videoPort uint16
	var audioPort uint16
	//var rtcPort uint16
	var err error

	resultPeerConnection := make(chan *webrtc.PeerConnection, 1)

	for {
		audioPort = pkg.RandUint16(49152, 65535)
		videoPort = pkg.RandUint16(49152, 65535)

		if audio, okAudio := controller.Ports[audioPort]; okAudio {
			if audio {
				continue
			}
			if video, okVideo := controller.Ports[videoPort]; okVideo {
				if video {
					continue
				}
			} else {
				controller.Ports[audioPort] = true
				break
			}
		} else {
			controller.Ports[videoPort] = true
			break
		}

	}

	//peerConnection, err := rtc.CreateSession(streamSlug, offer, rtcPort, videoPort, audioPort)
	//go rtc.CreateSession(streamSlug, offer, rtcPort, videoPort, audioPort, resultPeerConnection)

	go rtc.StartP2P(offer, int(audioPort), int(videoPort), resultPeerConnection)
	peerConnection := <-resultPeerConnection

	if err != nil {
		return &webrtc.PeerConnection{}, err
	}

	//generate SDP file containing IP,audioport, and videoport, for next to be used by FFMPEG
	_ = pkg.GenerateSDP(
		strconv.Itoa(int(audioPort)),
		strconv.Itoa(int(videoPort)), 
		streamSlug,
	)

	return peerConnection, nil
}
