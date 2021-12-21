package ffmpeg

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/asumsi/api.inlive/pkg"
)

func GetFileSDPPath(slug string) string {
	return pkg.GetCurrentDir() + "/" + os.Getenv("SDP_DIRECTORY") + "/" + slug + "-" + os.Getenv("SDP_FILENAME")
}

func GetStreamPath(slug string) string {
	return os.Getenv("DASH_SERVER") + "/ldash/" + slug + "/" + os.Getenv("MANIFEST_FILENAME")
}

func Execute(fileSDP string, urlStream string, pid chan int) {
	cmd := exec.Command("/bin/sh", "-c", "ffmpeg -protocol_whitelist file,udp,rtp -i "+fileSDP+" -c:v libx264 -g 120 -keyint_min 120 -b:v 4000k -vf \"fps=30\" -format_options movflags=cmaf -write_prft 1 -sc_threshold 0 -method PUT -ldash 1 -seg_duration 4 -frag_duration 1 -streaming 1 -http_persistent 1 -utc_timing_url \"https://time.akamai.com/?iso\" -use_timeline 0 -use_template 1 -media_seg_name 'chunk-stream-$RepresentationID$-$Number%05d$.m4s' -init_seg_name 'init-stream1-$RepresentationID$.m4s' -window_size 5 -extra_window_size 10 -remove_at_exit 1 -adaptation_sets \"id=0,streams=v id=1,streams=a\" -fflags nobuffer -f dash "+urlStream)

	stdout, err := cmd.StdoutPipe()

	cmd.Stderr = cmd.Stdout
	if err != nil {
		return
	}
	if err = cmd.Start(); err != nil {
		return
	}

	actualPid := cmd.Process.Pid
	pid <- actualPid

	i := 0

	for {
		tmp := make([]byte, 1024)
		_, err := stdout.Read(tmp)
		fmt.Print(string(tmp))
		if err != nil {
			break
		}
		i++
	}
}

func End(pid string) (string, error) {
	child, _ := exec.Command("pgrep", "-P", pid).Output()
	fmt.Println(child)
    if string(child) != "" {
		pid = string(child)
	}

	cmd := exec.Command("/bin/sh", "-c","kill -9 " + pid)

	stdout, err := cmd.Output()

    if err != nil {
        fmt.Println(err)
        return "gagal", err
    }

    // Print the output
    return string(stdout), nil
}
