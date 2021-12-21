package webrtc

import (
	"fmt"
	"net"
	"os"
	"time"

	"github.com/pion/interceptor"
	"github.com/pion/rtcp"
	"github.com/pion/rtp"
	"github.com/pion/webrtc/v3"
)

type Session struct {
	PeerConnection *webrtc.PeerConnection
	IceCandidates  []*webrtc.ICECandidate
}

var Ports [16383]int

var Sessions = make(map[string]*Session)

type udpConn struct {
	conn        *net.UDPConn
	port        uint16
	payloadType uint8
}

func StopSession(streamSlug string) {
	if session, ok := Sessions[streamSlug]; ok {
		//do something here
		session.PeerConnection.SCTP().Stop()
	}
}

func CreateSession(streamSlug string, offer webrtc.SessionDescription, rtcPort uint16, videoPort uint16, audioPort uint16) (*webrtc.PeerConnection, error) {
	peerConnection := createPeerConnection(rtcPort)
	Sessions[streamSlug] = &Session{PeerConnection: peerConnection}
	udpConns := openUDPPorts(videoPort, audioPort)

	setOnTrackEvent(udpConns, peerConnection)

	setOnIceConnectionStateEvent(peerConnection)

	setOnConnectionStateChangeEvent(peerConnection)

	setOnIceCandidateEvent(streamSlug, peerConnection)

	// Create channel that is blocked until ICE Gathering is complete
	gatherComplete := webrtc.GatheringCompletePromise(peerConnection)
	err := peerConnection.SetRemoteDescription(offer)
	if err != nil {
		fmt.Println("Failed set remote SDP")
		return &webrtc.PeerConnection{}, err
	}
	// Create answer
	answer, err := peerConnection.CreateAnswer(nil)
	if err != nil {
		fmt.Println("Failed to create answer")
		return &webrtc.PeerConnection{}, err
	}
	// Set the remote SessionDescription
	if err = peerConnection.SetLocalDescription(answer); err != nil {
		fmt.Println("Failed set local SDP")
		return &webrtc.PeerConnection{}, err
	}
	// Block until ICE Gathering is complete, disabling trickle ICE
	// we do this because we only can exchange one signaling message
	// in a production application you should exchange ICE Candidates via OnICECandidate
	<-gatherComplete
	return peerConnection, nil
}

func createPeerConnection(rtcPort uint16) *webrtc.PeerConnection {
	// Everything below is the Pion WebRTC API! Thanks for using it ❤️.

	// Create a MediaEngine object to configure the supported codec
	m := &webrtc.MediaEngine{}
	s := webrtc.SettingEngine{}
	s.SetEphemeralUDPPortRange(rtcPort, rtcPort)
	// Setup the codecs you want to use.
	// We'll use a VP8 and Opus but you can also define your own
	if err := m.RegisterCodec(webrtc.RTPCodecParameters{
		RTPCodecCapability: webrtc.RTPCodecCapability{MimeType: webrtc.MimeTypeH264, ClockRate: 90000, Channels: 0, SDPFmtpLine: "", RTCPFeedback: nil},
	}, webrtc.RTPCodecTypeVideo); err != nil {
		panic(err)
	}
	if err := m.RegisterCodec(webrtc.RTPCodecParameters{
		RTPCodecCapability: webrtc.RTPCodecCapability{MimeType: webrtc.MimeTypeOpus, ClockRate: 48000, Channels: 0, SDPFmtpLine: "", RTCPFeedback: nil},
	}, webrtc.RTPCodecTypeAudio); err != nil {
		panic(err)
	}

	// Create a InterceptorRegistry. This is the user configurable RTP/RTCP Pipeline.
	// This provides NACKs, RTCP Reports and other features. If you use `webrtc.NewPeerConnection`
	// this is enabled by default. If you are manually managing You MUST create a InterceptorRegistry
	// for each PeerConnection.
	i := &interceptor.Registry{}

	// Use the default set of Interceptors
	if err := webrtc.RegisterDefaultInterceptors(m, i); err != nil {
		panic(err)
	}

	// Create the API object with the MediaEngine
	api := webrtc.NewAPI(webrtc.WithMediaEngine(m), webrtc.WithInterceptorRegistry(i))

	// Prepare the configuration
	config := webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{
				URLs: []string{"stun:stun.l.google.com:19302"},
			},
		},
	}

	// Create a new RTCPeerConnection
	peerConnection, err := api.NewPeerConnection(config)

	if err != nil {
		panic(err)
	}
	// defer func() {
	// 	if cErr := peerConnection.Close(); cErr != nil {
	// 		fmt.Printf("cannot close peerConnection: %v\n", cErr)
	// 	}
	// }()

	// Allow us to receive 1 audio track, and 1 video track
	if _, err = peerConnection.AddTransceiverFromKind(webrtc.RTPCodecTypeAudio); err != nil {
		panic(err)
	} else if _, err = peerConnection.AddTransceiverFromKind(webrtc.RTPCodecTypeVideo); err != nil {
		panic(err)
	}
	return peerConnection
}

func openUDPPorts(videoPort uint16, audioPort uint16) map[string]*udpConn {
	// Create a local addr
	var laddr *net.UDPAddr
	var err error
	if laddr, err = net.ResolveUDPAddr("udp", "127.0.0.1:"); err != nil {
		panic(err)
	}
	// Prepare udp conns
	// Also update incoming packets with expected PayloadType, the browser may use
	// a different value. We have to modify so our stream matches what rtp-forwarder.sdp expects
	udpConns := map[string]*udpConn{
		"audio": {port: audioPort, payloadType: 111},
		"video": {port: videoPort, payloadType: 96},
	}
	for _, c := range udpConns {
		// Create remote addr
		var raddr *net.UDPAddr
		if raddr, err = net.ResolveUDPAddr("udp", fmt.Sprintf("127.0.0.1:%d", c.port)); err != nil {
			panic(err)
		}

		// Dial udp
		if c.conn, err = net.DialUDP("udp", laddr, raddr); err != nil {
			panic(err)
		}
		defer func(conn net.PacketConn) {
			if closeErr := conn.Close(); closeErr != nil {
				panic(closeErr)
			}
		}(c.conn)
	}
	return udpConns
}

func setOnTrackEvent(udpConns map[string]*udpConn, peerConnection *webrtc.PeerConnection) {
	// Set a handler for when a new remote track starts, this handler will forward data to
	// our UDP listeners.
	// In your application this is where you would handle/process audio/video
	peerConnection.OnTrack(func(track *webrtc.TrackRemote, receiver *webrtc.RTPReceiver) {
		var err error
		// Retrieve udp connection
		c, ok := udpConns[track.Kind().String()]
		if !ok {
			return
		}

		// Send a PLI on an interval so that the publisher is pushing a keyframe every rtcpPLIInterval
		go func() {
			ticker := time.NewTicker(time.Second * 2)
			for range ticker.C {
				if rtcpErr := peerConnection.WriteRTCP([]rtcp.Packet{&rtcp.PictureLossIndication{MediaSSRC: uint32(track.SSRC())}}); rtcpErr != nil {
					fmt.Println(rtcpErr)
				}
			}
		}()

		b := make([]byte, 1500)
		rtpPacket := &rtp.Packet{}
		for {
			// Read
			n, _, readErr := track.Read(b)
			if readErr != nil {
				if readErr.Error() == "EOF" {
					break
				}
				panic(readErr)
			}

			// Unmarshal the packet and update the PayloadType
			if err = rtpPacket.Unmarshal(b[:n]); err != nil {
				panic(err)
			}
			rtpPacket.PayloadType = c.payloadType

			// Marshal into original buffer with updated PayloadType
			if n, err = rtpPacket.MarshalTo(b); err != nil {
				panic(err)
			}

			// Write
			if _, err = c.conn.Write(b[:n]); err != nil {
				// For this particular example, third party applications usually timeout after a short
				// amount of time during which the user doesn't have enough time to provide the answer
				// to the browser.
				// That's why, for this particular example, the user first needs to provide the answer
				// to the browser then open the third party application. Therefore we must not kill
				// the forward on "connection refused" errors
				if opError, ok := err.(*net.OpError); ok && (opError.Err.Error() == "write: connection refused" || opError.Err.Error() == "use of closed network connection") {
					continue
				}
				panic(err)
			}
		}
	})
}

func setOnIceCandidateEvent(streamSlug string, peerConnection *webrtc.PeerConnection) {
	peerConnection.OnICECandidate(func(iceCandidate *webrtc.ICECandidate) {
		session := Sessions[streamSlug]
		session.IceCandidates = append(session.IceCandidates, iceCandidate)
	})
}

func setOnIceConnectionStateEvent(peerConnection *webrtc.PeerConnection) {
	// Set the handler for ICE connection state
	// This will notify you when the peer has connected/disconnected
	peerConnection.OnICEConnectionStateChange(func(connectionState webrtc.ICEConnectionState) {
		fmt.Printf("Connection State has changed %s \n", connectionState.String())

		if connectionState == webrtc.ICEConnectionStateConnected {
			fmt.Println("Ctrl+C the remote client to stop the demo")
		}
	})
}

func setOnConnectionStateChangeEvent(peerConnection *webrtc.PeerConnection) {
	// Set the handler for Peer connection state
	// This will notify you when the peer has connected/disconnected
	peerConnection.OnConnectionStateChange(func(s webrtc.PeerConnectionState) {
		fmt.Printf("Peer Connection State has changed: %s\n", s.String())

		if s == webrtc.PeerConnectionStateFailed {
			// Wait until PeerConnection has had no network activity for 30 seconds or another failure. It may be reconnected using an ICE Restart.
			// Use webrtc.PeerConnectionStateDisconnected if you are interested in detecting faster timeout.
			// Note that the PeerConnection may come back from PeerConnectionStateDisconnected.
			fmt.Println("Done forwarding")
			os.Exit(0)
		}
	})
}
