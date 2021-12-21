//go:build !js
// +build !js

package webrtc

import (
	"fmt"
	"net"
	//"os"
	"time"

	"github.com/pion/interceptor"
	"github.com/pion/rtcp"
	"github.com/pion/rtp"
	"github.com/pion/webrtc/v3"

	"bufio"
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"os/exec"
	//"github.com/pion/webrtc/v3/examples/internal/signal"
)

type udpConnRtpf struct {
	conn        *net.UDPConn
	port        int
	payloadType uint8
}

const compress = false

func StartP2P(offer webrtc.SessionDescription, audioPort int, videoPort int, result chan *webrtc.PeerConnection) {
	// Everything below is the Pion WebRTC API! Thanks for using it ❤️.

	// Create a MediaEngine object to configure the supported codec
	m := &webrtc.MediaEngine{}

	// Setup the codecs you want to use.
	// We'll use a VP8 and Opus but you can also define your own
	if err := m.RegisterCodec(webrtc.RTPCodecParameters{
		RTPCodecCapability: webrtc.RTPCodecCapability{MimeType: webrtc.MimeTypeVP8, ClockRate: 90000, Channels: 0, SDPFmtpLine: "", RTCPFeedback: nil},
	}, webrtc.RTPCodecTypeVideo); err != nil {
		fmt.Println(err)
	}
	if err := m.RegisterCodec(webrtc.RTPCodecParameters{
		RTPCodecCapability: webrtc.RTPCodecCapability{MimeType: webrtc.MimeTypeOpus, ClockRate: 48000, Channels: 0, SDPFmtpLine: "", RTCPFeedback: nil},
	}, webrtc.RTPCodecTypeAudio); err != nil {
		fmt.Println(err)
	}

	// Create a InterceptorRegistry. This is the user configurable RTP/RTCP Pipeline.
	// This provides NACKs, RTCP Reports and other features. If you use `webrtc.NewPeerConnection`
	// this is enabled by default. If you are manually managing You MUST create a InterceptorRegistry
	// for each PeerConnection.
	i := &interceptor.Registry{}

	// Use the default set of Interceptors
	if err := webrtc.RegisterDefaultInterceptors(m, i); err != nil {
		fmt.Println(err)
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
		fmt.Println(err)
	}
	defer func() {
		if cErr := peerConnection.Close(); cErr != nil {
			fmt.Printf("cannot close peerConnection: %v\n", cErr)
		}
	}()

	// Allow us to receive 1 audio track, and 1 video track
	if _, err = peerConnection.AddTransceiverFromKind(webrtc.RTPCodecTypeAudio); err != nil {
		fmt.Println(err)
	} else if _, err = peerConnection.AddTransceiverFromKind(webrtc.RTPCodecTypeVideo); err != nil {
		fmt.Println(err)
	}

	// Create a local addr
	var laddr *net.UDPAddr
	if laddr, err = net.ResolveUDPAddr("udp", "127.0.0.1:"); err != nil {
		fmt.Println(err)
	}

	// Prepare udp conns
	// Also update incoming packets with expected PayloadType, the browser may use
	// a different value. We have to modify so our stream matches what rtp-forwarder.sdp expects
	udpConns := map[string]*udpConnRtpf{
		"audio": {port: audioPort, payloadType: 111},
		"video": {port: videoPort, payloadType: 96},
	}
	for _, c := range udpConns {
		// Create remote addr
		var raddr *net.UDPAddr
		if raddr, err = net.ResolveUDPAddr("udp", fmt.Sprintf("127.0.0.1:%d", c.port)); err != nil {
			fmt.Println(err)
		}

		// Dial udp
		if c.conn, err = net.DialUDP("udp", laddr, raddr); err != nil {
			fmt.Println(err)
		}
		defer func(conn net.PacketConn) {
			if closeErr := conn.Close(); closeErr != nil {
				fmt.Println(closeErr)
			}
		}(c.conn)
	}

	// Set a handler for when a new remote track starts, this handler will forward data to
	// our UDP listeners.
	// In your application this is where you would handle/process audio/video
	peerConnection.OnTrack(func(track *webrtc.TrackRemote, receiver *webrtc.RTPReceiver) {
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
					fmt.Println("cek yg ticker")
					fmt.Println(rtcpErr)
				}
			}
		}()

		b := make([]byte, 1500)
		rtpPacket := &rtp.Packet{}
		counterDie := 0
		for {
			//modify bodoh2an
			if counterDie > 200 {
				break
			}

			// Read
			n, _, readErr := track.Read(b)
			if readErr != nil {
				fmt.Println("cek yg readbyte")
				fmt.Println(readErr)
			}

			// Unmarshal the packet and update the PayloadType
			if err = rtpPacket.Unmarshal(b[:n]); err != nil {
				counterDie++
				fmt.Println("cek yg unmarshal")
				fmt.Println(err)
			}
			rtpPacket.PayloadType = c.payloadType

			// Marshal into original buffer with updated PayloadType
			if n, err = rtpPacket.MarshalTo(b); err != nil {
				fmt.Println("cek yg marshal")
				fmt.Println(err)
			}

			// Write
			if _, err = c.conn.Write(b[:n]); err != nil {
				// For this particular example, third party applications usually timeout after a short
				// amount of time during which the user doesn't have enough time to provide the answer
				// to the browser.
				// That's why, for this particular example, the user first needs to provide the answer
				// to the browser then open the third party application. Therefore we must not kill
				// the forward on "connection refused" errors
				fmt.Println("error conn write nih")
				fmt.Println(err)

				if opError, ok := err.(*net.OpError); ok && opError.Err.Error() == "write: connection refused" {
					fmt.Println("masuk sini kah connection refused ???")
					continue
				}
				fmt.Println("cek yg ontrack")
				fmt.Println(err)
			}
		}
	})

	// Set the handler for ICE connection state
	// This will notify you when the peer has connected/disconnected
	peerConnection.OnICEConnectionStateChange(func(connectionState webrtc.ICEConnectionState) {
		fmt.Printf("Connection State has changed %s \n", connectionState.String())

		if connectionState == webrtc.ICEConnectionStateConnected {
			fmt.Println("Ctrl+C the remote client to stop the demo")
		}
	})

	// Set the handler for Peer connection state
	// This will notify you when the peer has connected/disconnected
	peerConnection.OnConnectionStateChange(func(s webrtc.PeerConnectionState) {
		fmt.Printf("Peer Connection State has changed: %s\n", s.String())

		if s == webrtc.PeerConnectionStateFailed {
			// Wait until PeerConnection has had no network activity for 30 seconds or another failure. It may be reconnected using an ICE Restart.
			// Use webrtc.PeerConnectionStateDisconnected if you are interested in detecting faster timeout.
			// Note that the PeerConnection may come back from PeerConnectionStateDisconnected.
			fmt.Println("Done forwarding")
			//os.Exit(0)
		}
	})

	// Wait for the offer to be pasted
	//offer := webrtc.SessionDescription{}
	//Decode(stdin, &offer)

	// Set the remote SessionDescription
	if err = peerConnection.SetRemoteDescription(offer); err != nil {
		fmt.Println(err)
	}

	// Create answer
	answer, err := peerConnection.CreateAnswer(nil)
	if err != nil {
		fmt.Println(err)
	}

	// Create channel that is blocked until ICE Gathering is complete
	gatherComplete := webrtc.GatheringCompletePromise(peerConnection)

	// Sets the LocalDescription, and starts our UDP listeners
	if err = peerConnection.SetLocalDescription(answer); err != nil {
		fmt.Println(err)
	}

	// Block until ICE Gathering is complete, disabling trickle ICE
	// we do this because we only can exchange one signaling message
	// in a production application you should exchange ICE Candidates via OnICECandidate
	<-gatherComplete

	// Output the answer in base64 so we can paste it in browser
	//fmt.Println(Encode(*peerConnection.LocalDescription()))

	result <- peerConnection

	// Block forever
	select {}

	return
}

func Decode(in string, obj interface{}) {
	b, err := base64.StdEncoding.DecodeString(in)
	if err != nil {
		fmt.Println(err)
	}

	if compress {
		b = unzip(b)
	}

	err = json.Unmarshal(b, obj)
	if err != nil {
		fmt.Println(err)
	}
}

func Encode(obj interface{}) string {
	b, err := json.Marshal(obj)
	if err != nil {
		fmt.Println(err)
	}

	if compress {
		b = zip(b)
	}

	return base64.StdEncoding.EncodeToString(b)
}

func unzip(in []byte) []byte {
	var b bytes.Buffer
	_, err := b.Write(in)
	if err != nil {
		fmt.Println(err)
	}
	r, err := gzip.NewReader(&b)
	if err != nil {
		fmt.Println(err)
	}
	res, err := ioutil.ReadAll(r)
	if err != nil {
		fmt.Println(err)
	}
	return res
}

func zip(in []byte) []byte {
	var b bytes.Buffer
	gz := gzip.NewWriter(&b)
	_, err := gz.Write(in)
	if err != nil {
		fmt.Println(err)
	}
	err = gz.Flush()
	if err != nil {
		fmt.Println(err)
	}
	err = gz.Close()
	if err != nil {
		fmt.Println(err)
	}
	return b.Bytes()
}

func ExecuteRTPForwarder(portAudio string, portVideo string, offerSession string, answerSession chan string, pid chan int) {
	commandExecRTPForwarder := "echo " + portAudio + " " + portVideo + " " + offerSession + " | ./rtp-forwarder"
	fmt.Println(commandExecRTPForwarder)
	cmd := exec.Command("/bin/sh", "-c", commandExecRTPForwarder)

	stdout, err := cmd.StdoutPipe()

	cmd.Stderr = cmd.Stdout
	if err != nil {
		fmt.Println("err1")
		fmt.Println(err)
	}
	if err = cmd.Start(); err != nil {
		fmt.Println("err2")
		fmt.Println(err)
	}

	actualPid := cmd.Process.Pid
	fmt.Println(actualPid)
	pid <- actualPid

	i := 0
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		line := scanner.Text()
		if i == 1 {
			fmt.Println(line)
			answerSession <- line
		}
		i++
	}
}
