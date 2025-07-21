package ffmpegcore

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

type TinyFfmpegX11 struct {
	Framerate  int
	OutputPath string
}

var framerate = 25
var outputPath = "outputs/"

func (t *TinyFfmpegX11) RecordScreen() error {
	// ffmpeg -video_size $(xrandr | grep '*' | awk '{print $1}') -framerate 25 -f x11grab -i :1+0,0 -f pulse -ac 2 -i default output.mkv
	resolution, err := t.Resolution()
	if err != nil {
		return err
	}
	fmt.Println("screen resolution ", resolution)
	cmd := exec.Command("ffmpeg",
		"-video_size", resolution,
		"-framerate", strconv.Itoa(framerate),
		"-f", "x11grab",
		"-i", ":1+0,0",
		"-f", "pulse",
		"-ac", "2",
		"-i", "default",
		"-y",
		fmt.Sprintf("%s%s", outputPath, "output.mkv"),
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to record screen: %w", err)
	}

	fmt.Println(string(output))
	return nil
}

func (t *TinyFfmpegX11) Resolution() (string, error) {
	cmd := exec.Command("bash", "-c", "xrandr | grep '*' | awk '{print $1}'")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	
	// Convert output to string and split by lines
	lines := strings.Split(string(output), "\n")
	if len(lines) == 0 || lines[0] == "" {
		return "", fmt.Errorf("no resolution found")
	}

	// Return the first resolution line trimmed
	return strings.TrimSpace(lines[0]), nil
	// return strings.TrimSpace(string(output)), nil
}


func (t *TinyFfmpegX11) StartStream(protocol, address string) error {
	resolution, err := t.Resolution()
	if err != nil {
		return err
	}
	
	cmd := exec.Command("ffmpeg",
		"-video_size", resolution,
		"-framerate", "25",
		"-f", "x11grab",
		"-i", ":1+0,0",
		"-c:v", "libx264",
		"-preset", "ultrafast",
		"-f", "mpegts",
		fmt.Sprintf("%s://%s", protocol, address),
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to record screen: %w", err)
	}
	fmt.Println(string(output))
	return nil
}