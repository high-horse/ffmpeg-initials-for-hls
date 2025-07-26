package main

import (
	"fmt"
	"log"
	"time"
	tf "tiny-screen/ffmpeg-core"
)

func main() {
	encodeHlsRetry(3)
	serve()
}

func encodeHlsRetry(maxRetry int) error{
	x11 := tf.TinyFfmpegX11{}
	var err error
	for attempt := 1; attempt <= maxRetry; attempt++ {
		fmt.Printf("Attempt %d of %d", attempt, maxRetry)
		err = x11.HlsEncode("input-samples/sample-video.mp4", "", true)
		if err == nil {
			return nil
		}
		fmt.Printf("FFmpeg failed: %v\n", err)
		if attempt < maxRetry {
			time.Sleep(time.Second * 2)
		}
	}
	return fmt.Errorf("all %d ffmpeg attempts failed: last error: %w", maxRetry, err)
}

func recordScreen() {
	x11 := tf.TinyFfmpegX11{}
	if err := x11.RecordScreen(); err != nil {
		log.Fatal(err)
	}
}

func getResolution() {
	x11 := tf.TinyFfmpegX11{}
	log.Println(x11.Resolution())
}

func captureImg() {
	x11 := tf.TinyFfmpegX11{}
	if err := x11.CaptureImage(1); err != nil {
		log.Fatal(err)
	}
}

func streamScreen() {
	x11 := tf.TinyFfmpegX11{}
	protocol := "udp"
	address := "127.0.0.1:1234"
	if err := x11.StartStream(protocol, address); err != nil {
		log.Fatal(err)
	}
}
