### Base command for screen recording with ffmpeg

```bash
ffmpeg -video_size 1024x768 -framerate 25 -f x11grab -i :1+100,200 output.mp4
```

### Get screen resolution:

#### Using srandr

```bash
xrandr | grep '*' | awk '{print $1}'
```

#### Using xdpyinfo:

```bash
xdpyinfo -display :1 | grep dimensions | awk '{print $2}'
```

### Examples of screen recording commands`

incorporate screen resolution with base:

```bash
ffmpeg -video_size $(xrandr | grep '*' | awk '{print $1}') -framerate 25 -f x11grab -i :1+100,200 output.mp4
ffmpeg -video_size $(xdpyinfo -display :1 | grep dimensions | awk '{print $2}') -framerate 25 -f x11grab -i :1+100,200 output.mp4
```

In case the above command throws output or offset errors :

```bash
ffmpeg -video_size $(xrandr | grep '*' | awk '{print $1}') -framerate 25 -f x11grab -i :1+0,0 -f pulse -ac 2 -i default output.mkv
ffmpeg -video_size $(xrandr | grep '*' | awk '{print $1}') -framerate 25 -f x11grab -i :1+0,0 output.mp4
ffmpeg -video_size $(xdpyinfo -display :1 | grep dimensions | awk '{print $2}') -framerate 25 -f x11grab -i :1+0,0 output.mp4
```

** Play recorded video **
```bash 
ffplay outputs/output.mp4 -autoexit
```

### Stream the recording 
To stream the captured video from your Ubuntu system using ffmpeg (with the x11grab command you provided) over TCP or UDP, you can modify the ffmpeg command to output to a streaming protocol like MPEG-TS over TCP or UDP. Since you’re using a Go wrapper, I’ll provide the ffmpeg command for streaming and explain how to integrate it into a Go program, assuming the wrapper executes ffmpeg commands. Below, I’ll cover streaming over TCP and UDP, and how to handle this in a Go context.

1. Streaming with ffmpeg
Current code to capture recording
```bash
ffmpeg -video_size $(xrandr | grep '*' | awk '{print $1}') -framerate 25 -f x11grab -i :1+0,0 output.mp4
```

** Streaming over TCP **
MPEG-TS (MPEG Transport Stream) is commonly used for streaming over TCP. Modify the command to stream to a TCP server:
```bash 
ffmpeg -video_size $(xrandr | grep '*' | awk '{print $1}') -framerate 25 -f x11grab -i :1+0,0 -c:v libx264 -preset ultrafast -f mpegts tcp://127.0.0.1:1234
```
- -c:v libx264: Encodes the video with H.264 for efficient streaming.
- -preset ultrafast: Optimizes for low latency and real-time streaming.
- -f mpegts: Uses MPEG-TS format, suitable for TCP streaming.
- tcp://127.0.0.1:1234: Streams to localhost on port 1234. Replace 127.0.0.1 with the target IP if streaming to another machine.

** Receiver Side (to view the stream): **
Run a client to receive and play the stream, e.g., using ffplay:
```bash 
ffplay tcp://127.0.0.1:1234
```

** Streaming over UDP **
UDP is simpler but less reliable (packets may be lost). Modify the command:
```bash 
ffmpeg -video_size $(xrandr | grep '*' | awk '{print $1}') -framerate 25 -f x11grab -i :1+0,0 -c:v libx264 -preset ultrafast -f mpegts udp://127.0.0.1:1234
```
- -f mpegts: Uses MPEG-TS for UDP streaming.
- udp://127.0.0.1:1234: Streams to localhost on port 1234. Replace with the target IP as needed.

** Receiver Side: **
```bash 
ffplay udp://127.0.0.1:1234
```

** Adding Audio (Optional) **

If you want to include audio (e.g., using PulseAudio), add the audio input:
```bash
ffmpeg -video_size $(xrandr | grep '*' | awk '{print $1}') -framerate 25 -f x11grab -i :1+0,0 -f pulse -ac 2 -i default -c:v libx264 -c:a aac -preset ultrafast -f mpegts tcp://127.0.0.1:1234
```

- -f pulse -ac 2 -i default: Captures stereo audio from PulseAudio.
- -c:a aac: Encodes audio with AAC.

** Integrating with a Go Wrapper **
```bash
package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func getScreenResolution() (string, error) {
	cmd := exec.Command("bash", "-c", "xrandr | grep '*' | awk '{print $1}'")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get screen resolution: %v", err)
	}
	return strings.TrimSpace(string(output)), nil
}

func startFFmpegStream(resolution, protocol, address string) error {
	ffmpegCmd := []string{
		"ffmpeg",
		"-video_size", resolution,
		"-framerate", "25",
		"-f", "x11grab",
		"-i", ":1+0,0",
		"-c:v", "libx264",
		"-preset", "ultrafast",
		"-f", "mpegts",
		fmt.Sprintf("%s://%s", protocol, address),
	}

	cmd := exec.Command(ffmpegCmd[0], ffmpegCmd[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Printf("Running: %s\n", strings.Join(ffmpegCmd, " "))
	return cmd.Run()
}

func main() {
	// Get screen resolution
	resolution, err := getScreenResolution()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Stream to TCP or UDP
	protocol := "tcp" // or "udp"
	address := "127.0.0.1:1234" // Change to target IP:port
	err = startFFmpegStream(resolution, protocol, address)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Streaming failed: %v\n", err)
		os.Exit(1)
	}
}
```

#### On the receiver side
For TCP
```bash
ffplay tcp://127.0.0.1:1234
```
Or for UDP
```bash
ffplay -fflags nobuffer -flags low_delay -probesize 32 -analyzeduration 0 udp://127.0.0.1:1234 -autoexit
ffplay udp://127.0.0.1:1234
```
** check window session **
```bash
echo $XDG_SESSION_TYPE
```

** Capture image with in built cam with delay of 3 sec **
```bash
sleep 3 && ffmpeg -y -f v4l2 -i /dev/video0 -frames:v 1 -update 1 outputs/cam.jpg

# to draw text use flag 
# -vf "drawtext=text='%{localtime}':fontcolor=white:fontsize=24:x=10:y=10" 
# before outputs
```

** Capture Screenshot **
```bash
ffmpeg -y -video_size $(xrandr | grep '*' | awk '{print $1}') -f x11grab -i :1+0,0 -frames:v 1 outputs/screenshot.png
```

** Merge cam image with screenshot **
```bash
ffmpeg -y -i outputs/screenshot.png -i outputs/cam.jpg -filter_complex "overlay=W-w-10:H-h-10" outputs/merge.jpg
```