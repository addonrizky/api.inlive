package pkg

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func GenerateSDP(portAudio string, portVideo string, slug string) string {
	var err error
	currentDir := GetCurrentDir()
	fileSDP := currentDir + "/" + os.Getenv("SDP_DIRECTORY") + "/" + slug + "-" + os.Getenv("SDP_FILENAME")
	f, err := os.Create(fileSDP)
	check(err)

	defer f.Close()

	line := []byte("v=0\n")
	_, err = f.Write(line)
	check(err)

	line = []byte("o=- 0 0 IN IP4 127.0.0.1\n")
	_, err = f.Write(line)
	check(err)

	line = []byte("s=Pion WebRTC\n")
	_, err = f.Write(line)
	check(err)

	line = []byte("c=IN IP4 127.0.0.1\n")
	_, err = f.Write(line)
	check(err)

	line = []byte("t=0 0\n")
	_, err = f.Write(line)
	check(err)

	line = []byte("m=audio " + portAudio + " RTP/AVP 111\n")
	_, err = f.Write(line)
	check(err)

	line = []byte("a=rtpmap:111 OPUS/48000/2\n")
	_, err = f.Write(line)
	check(err)

	line = []byte("m=video " + portVideo + " RTP/AVP 96\n")
	_, err = f.Write(line)
	check(err)

	line = []byte("a=rtpmap:96 VP8/90000")
	_, err = f.Write(line)
	check(err)

	f.Sync()

	return fileSDP
}

func UpdateFFMPEGFile(sdpFile string, urlStream string) string {
	input, err := ioutil.ReadFile("ffmpeg.sh")
	if err != nil {
		fmt.Println(err)
	}

	lines := strings.Split(string(input), "\n")

	for i, line := range lines {
		fmt.Println(line)
		if i == 4 {
			lines[i] = sdpFile
		}

		if i == 33 {
			lines[i] = urlStream
		}
	}
	output := strings.Join(lines, "\n")
	err = ioutil.WriteFile("ffmpeg.sh", []byte(output), 0644)
	if err != nil {
		fmt.Println(err)
	}

	return ""
}

func GetCurrentDir() string {
	pwd, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return pwd
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
