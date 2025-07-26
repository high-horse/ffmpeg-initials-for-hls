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

func (t *TinyFfmpegX11) ClickImage() error {
	return nil
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
		"-tune", "zerolatency", // Reduce encoding latency
		"-g", "25", // keyframe every 25 frames
		"-x264-params", "repeat-headers=1", // repeat SPS/PPS before keyframes
		"-b:v", "3000k", // Set bitrate to 3 Mbps
		"-maxrate", "3000k",
		"-bufsize", "1000k", // Smaller buffer size
		"-flags", "low_delay", // Low-latency encoding
		"-max_delay", "0", // Minimize output buffer
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

func (t *TinyFfmpegX11) CaptureImage(delay int) error {
	cmdStr := fmt.Sprintf(
		"sleep %d && ffmpeg -y -f v4l2 -i /dev/video0 -frames:v 1 -update 1 outputs/appCam.jpg",
		delay,
	)

	cmd := exec.Command("sh", "-c", cmdStr)

	op, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to capture screenshot: %w\nOutput: %s", err, string(op))
	}

	fmt.Println(string(op))
	return nil
}

func (t *TinyFfmpegX11) CaptureImageOld(delay int) error {
	// sleep 3 && ffmpeg -y -f v4l2 -i /dev/video0 -frames:v 1 -update 1 outputs/cam.jpg

	cmd := exec.Command(fmt.Sprintf("sleep %d && ffmpeg", delay),
		"-y",
		"if", "v4l2",
		"-i", "/dev/video0",
		"-frames:v", "1",
		"-update", "1",
		"outputs/appCam.jpg",
	)
	op, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to capture screenshot: %w", err)
	}
	fmt.Sprintln(string(op))
	return nil
}

func (t *TinyFfmpegX11) HlsEncode(mediaPath, outputPath string, adaptive bool) error {
	return hlsEncodeNormal(mediaPath, outputPath)
}

// func hlsEncodeAdaptive(mediaPath, outputPath string) error {
// 	// ffmpeg -i samples/sample-video.mp4 \
// 	//  -filter_complex "[0:v]split=3[v1][v2][v3]; \
// 	//                   [v1]scale=w=1920:h=1080[vout1]; \
// 	//                   [v2]scale=w=1280:h=720[vout2]; \
// 	//                   [v3]scale=w=854:h=480[vout3]" \
// 	//  -map "[vout1]" -map 0:a -c:v:0 libx264 -b:v:0 3000k -c:a:0 aac -f hls \
// 	//    -hls_time 6 -hls_playlist_type vod -hls_segment_filename outputs/hls-adaptive/1080p_%03d.ts outputs/hls-adaptive/1080p.m3u8 \
// 	//  -map "[vout2]" -map 0:a -c:v:1 libx264 -b:v:1 1500k -c:a:1 aac -f hls \
// 	//    -hls_time 6 -hls_playlist_type vod -hls_segment_filename outputs/hls-adaptive/720p_%03d.ts outputs/hls-adaptive/720p.m3u8 \
// 	//  -map "[vout3]" -map 0:a -c:v:2 libx264 -b:v:2 800k  -c:a:2 aac -f hls \
// 	//    -hls_time 6 -hls_playlist_type vod -hls_segment_filename outputs/hls-adaptive/480p_%03d.ts outputs/hls-adaptive/480p.m3u8
// 	//
// 	cmd := exec.Command("ffmpeg",
// 		"-i", mediaPath,
// 		"-filter_complex"
// 	)
// 	return nil
// }

func hlsEncodeNormal(mediaPath, outputPath string) error {
	fmt.Println(outputPath)
	cmd := exec.Command("ffmpeg",
		"-i", mediaPath,
		"-c", "copy",
		"-hls_time", "10",
		"-hls_list_size", "0",
		"-start_number", "0",
		"-hls_segment_filename", "outputs/hls/segment_%03d.ts",
		"-f", "hls", "outputs/hls/sample-hls-op.m3u8",
	)
	output, err := cmd.CombinedOutput()
	fmt.Println(string(output))
	if err != nil {
		return fmt.Errorf("failed to encode hls: %w", err)
	}
	return nil
}
