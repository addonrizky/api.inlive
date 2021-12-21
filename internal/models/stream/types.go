package stream

import (
	"time"

	"github.com/asumsi/api.inlive/pkg/api"
	"github.com/pion/webrtc/v3"
)

type Stream struct {
	ID           int64      `json:"id"`
	Name         string     `json:"name" validate:"required"`
	Slug         string     `json:"slug" validate:"required"`
	StartDate    *time.Time  `json:"start_date`
	EndDate      *time.Time  `json:"end_date`
	ManifestPath string	`json:"manifest_path"`
	Description  string     `json:"description"`
	CreatedBy    int64      `json:"created_by"`
	CreatedDate  time.Time  `json:"created_date"`
	UpdatedBy    *int64     `json:"updated_by"`
	UpdatedDate  *time.Time `json:"updated_date"`
}

type StreamResponse struct {
	api.Response
	Data Stream `json:"data"`
}

type StreamsResponse struct {
	api.Response
	Data []Stream `json:"data"`
}

type StreamParams struct{
	Live bool
}

type InitSessionRequest struct {
	Slug string `json:"slug" validate:"required"`
	SessionDescription webrtc.SessionDescription `json:"session_description" validate:"required"`
}

type StartStreamRequest struct {
	Slug string `json:"slug" validate:"required" example:"aff-suzuki-2020"`
}

type EndSessionRequest struct {
	PID string `json:"pid" validate:"required"`
}

type ErrorObject struct {
	ErrorActual  error
	ErrorCode    string
	ErrorMessage string
}

/* #################################################### */

type InitSessionRequestSwag struct {
	Slug string `json:"slug" validate:"required" example:"debat-kusir-pilpres"`
	SessionDescription SessionDescriptionSwag `json:"session_description" validate:"required"`
}

type SessionDescriptionSwag struct {
	Type string `json:"type" validate:"required" example:"offer"`
	Sdp string `json:"sdp" validate:"required" example:"v=0\r\no=- 1476873138974450762 2 IN IP4 127.0.0.1\r\ns=-\r\nt=0 0\r\na=group:BUNDLE 0 1\r\na=msid-semantic: WMS 5dceda4c-fa94-4f54-91f2-632d4e28fb04\r\nm=video 51742 UDP/TLS/RTP/SAVPF 96 97 98 99 100 101 127 125 104\r\nc=IN IP4 111.95.233.143\r\na=rtcp:9 IN IP4 0.0.0.0\r\na=candidate:1914645175 1 udp 2113937151 192.168.0.190 51742 typ host generation 0 network-cost 999\r\na=candidate:842163049 1 udp 1677729535 111.95.233.143 51742 typ srflx raddr 192.168.0.190 rport 51742 generation 0 network-cost 999\r\na=ice-ufrag:9AVD\r\na=ice-pwd:Qaqw3KJClEglwziAV5jY7gHA\r\na=ice-options:trickle\r\na=fingerprint:sha-256 D1:39:22:9F:9C:79:E6:65:89:29:FD:4E:CC:8B:A8:D9:90:11:57:76:A3:33:D0:43:4B:FD:A7:AD:11:8E:8D:DC\r\na=setup:actpass\r\na=mid:0\r\na=extmap:14 urn:ietf:params:rtp-hdrext:toffset\r\na=extmap:2 http://www.webrtc.org/experiments/rtp-hdrext/abs-send-time\r\na=extmap:13 urn:3gpp:video-orientation\r\na=extmap:3 http://www.ietf.org/id/draft-holmer-rmcat-transport-wide-cc-extensions-01\r\na=extmap:12 http://www.webrtc.org/experiments/rtp-hdrext/playout-delay\r\na=extmap:11 http://www.webrtc.org/experiments/rtp-hdrext/video-content-type\r\na=extmap:7 http://www.webrtc.org/experiments/rtp-hdrext/video-timing\r\na=extmap:8 http://tools.ietf.org/html/draft-ietf-avtext-framemarking-07\r\na=extmap:9 http://www.webrtc.org/experiments/rtp-hdrext/color-space\r\na=extmap:4 urn:ietf:params:rtp-hdrext:sdes:mid\r\na=extmap:5 urn:ietf:params:rtp-hdrext:sdes:rtp-stream-id\r\na=extmap:6 urn:ietf:params:rtp-hdrext:sdes:repaired-rtp-stream-id\r\na=sendrecv\r\na=msid:5dceda4c-fa94-4f54-91f2-632d4e28fb04 4c1cce82-0dcb-4aeb-8a45-09ed464b8e7b\r\na=rtcp-mux\r\na=rtcp-rsize\r\na=rtpmap:96 H264/90000\r\na=rtcp-fb:96 goog-remb\r\na=rtcp-fb:96 transport-cc\r\na=rtcp-fb:96 ccm fir\r\na=rtcp-fb:96 nack\r\na=rtcp-fb:96 nack pli\r\na=fmtp:96 level-asymmetry-allowed=1;packetization-mode=1;profile-level-id=640c1f\r\na=rtpmap:97 rtx/90000\r\na=fmtp:97 apt=96\r\na=rtpmap:98 H264/90000\r\na=rtcp-fb:98 goog-remb\r\na=rtcp-fb:98 transport-cc\r\na=rtcp-fb:98 ccm fir\r\na=rtcp-fb:98 nack\r\na=rtcp-fb:98 nack pli\r\na=fmtp:98 level-asymmetry-allowed=1;packetization-mode=1;profile-level-id=42e01f\r\na=rtpmap:99 rtx/90000\r\na=fmtp:99 apt=98\r\na=rtpmap:100 VP8/90000\r\na=rtcp-fb:100 goog-remb\r\na=rtcp-fb:100 transport-cc\r\na=rtcp-fb:100 ccm fir\r\na=rtcp-fb:100 nack\r\na=rtcp-fb:100 nack pli\r\na=rtpmap:101 rtx/90000\r\na=fmtp:101 apt=100\r\na=rtpmap:127 red/90000\r\na=rtpmap:125 rtx/90000\r\na=fmtp:125 apt=127\r\na=rtpmap:104 ulpfec/90000\r\na=ssrc-group:FID 3021001483 2138149310\r\na=ssrc:3021001483 cname:aryKLH6+2OkdDLnG\r\na=ssrc:3021001483 msid:5dceda4c-fa94-4f54-91f2-632d4e28fb04 4c1cce82-0dcb-4aeb-8a45-09ed464b8e7b\r\na=ssrc:3021001483 mslabel:5dceda4c-fa94-4f54-91f2-632d4e28fb04\r\na=ssrc:3021001483 label:4c1cce82-0dcb-4aeb-8a45-09ed464b8e7b\r\na=ssrc:2138149310 cname:aryKLH6+2OkdDLnG\r\na=ssrc:2138149310 msid:5dceda4c-fa94-4f54-91f2-632d4e28fb04 4c1cce82-0dcb-4aeb-8a45-09ed464b8e7b\r\na=ssrc:2138149310 mslabel:5dceda4c-fa94-4f54-91f2-632d4e28fb04\r\na=ssrc:2138149310 label:4c1cce82-0dcb-4aeb-8a45-09ed464b8e7b\r\nm=audio 60415 UDP/TLS/RTP/SAVPF 111 103 9 102 0 8 105 13 110 113 126\r\nc=IN IP4 111.95.233.143\r\na=rtcp:9 IN IP4 0.0.0.0\r\na=candidate:1914645175 1 udp 2113937151 192.168.0.190 60415 typ host generation 0 network-cost 999\r\na=candidate:842163049 1 udp 1677729535 111.95.233.143 60415 typ srflx raddr 192.168.0.190 rport 60415 generation 0 network-cost 999\r\na=ice-ufrag:9AVD\r\na=ice-pwd:Qaqw3KJClEglwziAV5jY7gHA\r\na=ice-options:trickle\r\na=fingerprint:sha-256 D1:39:22:9F:9C:79:E6:65:89:29:FD:4E:CC:8B:A8:D9:90:11:57:76:A3:33:D0:43:4B:FD:A7:AD:11:8E:8D:DC\r\na=setup:actpass\r\na=mid:1\r\na=extmap:1 urn:ietf:params:rtp-hdrext:ssrc-audio-level\r\na=extmap:2 http://www.webrtc.org/experiments/rtp-hdrext/abs-send-time\r\na=extmap:3 http://www.ietf.org/id/draft-holmer-rmcat-transport-wide-cc-extensions-01\r\na=extmap:4 urn:ietf:params:rtp-hdrext:sdes:mid\r\na=extmap:5 urn:ietf:params:rtp-hdrext:sdes:rtp-stream-id\r\na=extmap:6 urn:ietf:params:rtp-hdrext:sdes:repaired-rtp-stream-id\r\na=sendrecv\r\na=msid:5dceda4c-fa94-4f54-91f2-632d4e28fb04 8d2ef1f4-8be7-415a-a10c-cce040166924\r\na=rtcp-mux\r\na=rtpmap:111 opus/48000/2\r\na=rtcp-fb:111 transport-cc\r\na=fmtp:111 minptime=10;useinbandfec=1\r\na=rtpmap:103 ISAC/16000\r\na=rtpmap:9 G722/8000\r\na=rtpmap:102 ILBC/8000\r\na=rtpmap:0 PCMU/8000\r\na=rtpmap:8 PCMA/8000\r\na=rtpmap:105 CN/16000\r\na=rtpmap:13 CN/8000\r\na=rtpmap:110 telephone-event/48000\r\na=rtpmap:113 telephone-event/16000\r\na=rtpmap:126 telephone-event/8000\r\na=ssrc:3405135432 cname:aryKLH6+2OkdDLnG\r\na=ssrc:3405135432 msid:5dceda4c-fa94-4f54-91f2-632d4e28fb04 8d2ef1f4-8be7-415a-a10c-cce040166924\r\na=ssrc:3405135432 mslabel:5dceda4c-fa94-4f54-91f2-632d4e28fb04\r\na=ssrc:3405135432 label:8d2ef1f4-8be7-415a-a10c-cce040166924\r\n"`
}

type ResponseSwagInitSessSuccess struct {
	Code string               `json:"code" example:"200"`
	Desc string               `json:"message"  example:"Stream initiated"`
	Data ResponseSwagInitSessSdp `json:"data"`
}

type ResponseSwagInitSessFail struct {
	Code string               `json:"code" example:"400"`
	Desc string               `json:"message"  example:"Validation Error on Init Stream"`
	Data string               `json:"data"  example:"Key: 'InitSessionRequest.Slug' Error:Field validation for 'Slug' failed on the 'required' tag"`
}

type ResponseSwagInitSessSdp struct {
	Type string `json:"type" example:"answer"`
	Sdp  string `json:"sdp" example:"v=0\r\no=- 1709258760602786429 1639914210 IN IP4 0.0.0.0\r\ns=-\r\nt=0 0\r\na=fingerprint:sha-256 4C:07:84:8F:54:2D:9D:35:E1:D1:7D:53:59:FC:1B:9C:82:34:41:AC:77:B6:9F:C7:17:0D:69:4A:9C:F7:41:4D\r\na=group:BUNDLE 0 1\r\nm=video 9 UDP/TLS/RTP/SAVPF 100\r\nc=IN IP4 0.0.0.0\r\na=setup:active\r\na=mid:0\r\na=ice-ufrag:oJMXxoToKKbabFFN\r\na=ice-pwd:GDmxypmUAbeeUIHxiIREGlorXXSZsQan\r\na=rtcp-mux\r\na=rtcp-rsize\r\na=rtpmap:100 VP8/90000\r\na=rtcp-fb:100 goog-remb \r\na=rtcp-fb:100 transport-cc \r\na=rtcp-fb:100 ccm fir\r\na=rtcp-fb:100 nack \r\na=rtcp-fb:100 nack pli\r\na=extmap:3 http://www.ietf.org/id/draft-holmer-rmcat-transport-wide-cc-extensions-01\r\na=ssrc:3899392441 cname:abMWosTLSEkkNJWP\r\na=ssrc:3899392441 msid:abMWosTLSEkkNJWP EzeXWAOrEgIkJhBE\r\na=ssrc:3899392441 mslabel:abMWosTLSEkkNJWP\r\na=ssrc:3899392441 label:EzeXWAOrEgIkJhBE\r\na=msid:abMWosTLSEkkNJWP EzeXWAOrEgIkJhBE\r\na=sendrecv\r\na=candidate:201991694 1 udp 2130706431 192.168.0.190 64507 typ host\r\na=candidate:201991694 2 udp 2130706431 192.168.0.190 64507 typ host\r\na=candidate:94524127 1 udp 2130706431 192.168.2.1 58417 typ host\r\na=candidate:94524127 2 udp 2130706431 192.168.2.1 58417 typ host\r\na=candidate:1268538823 1 udp 1694498815 111.95.233.143 56624 typ srflx raddr 0.0.0.0 rport 56624\r\na=candidate:1268538823 2 udp 1694498815 111.95.233.143 56624 typ srflx raddr 0.0.0.0 rport 56624\r\na=end-of-candidates\r\nm=audio 9 UDP/TLS/RTP/SAVPF 111\r\nc=IN IP4 0.0.0.0\r\na=setup:active\r\na=mid:1\r\na=ice-ufrag:oJMXxoToKKbabFFN\r\na=ice-pwd:GDmxypmUAbeeUIHxiIREGlorXXSZsQan\r\na=rtcp-mux\r\na=rtcp-rsize\r\na=rtpmap:111 opus/48000/2\r\na=fmtp:111 minptime=10;useinbandfec=1\r\na=rtcp-fb:111 transport-cc \r\na=extmap:3 http://www.ietf.org/id/draft-holmer-rmcat-transport-wide-cc-extensions-01\r\na=ssrc:26062976 cname:dpKeKWqkpLAZOHfm\r\na=ssrc:26062976 msid:dpKeKWqkpLAZOHfm cWjrooxUMyGFKbgu\r\na=ssrc:26062976 mslabel:dpKeKWqkpLAZOHfm\r\na=ssrc:26062976 label:cWjrooxUMyGFKbgu\r\na=msid:dpKeKWqkpLAZOHfm cWjrooxUMyGFKbgu\r\na=sendrecv\r\n"`
}

type ResponseSwagStartStreamSuccess struct {
	Code string               `json:"code" example:"200"`
	Desc string               `json:"message"  example:"Stream run successfully"`
	Data ResponseSwagStartStreamFFMPEG `json:"data"`
}

type ResponseSwagStartStreamFFMPEG struct {
	UrlStream string `json:"url_stream" example:"https://bifrost.inlive.app/ldash/aff-suzuki-7/manifest.mpd"`
	Pid int `json:"pid" example:"14981"`
}

type ResponseSwagStartStreamFail struct {
	Code string               `json:"code" example:"200"`
	Desc string               `json:"message"  example:"Stream never initiated"`
	Data string `json:"data" example:"null"`
}

type ResponseSwagEndStreamSuccess struct {
	Code string               `json:"code" example:"200"`
	Desc string               `json:"message"  example:"Streaming stop"`
	Data string `json:"data" example:""`
}

type ResponseSwagEndStreamFail struct {
	Code string               `json:"code" example:"400"`
	Desc string               `json:"message"  example:"Stream never initiated"`
	Data string `json:"data" example:""`
}



func NewErrorObject() *ErrorObject {
	return &ErrorObject{
		ErrorActual:  nil,
		ErrorCode:    "GE",
		ErrorMessage: "General Error",
	}
}
