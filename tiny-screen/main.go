package main

import (
	"log"
	tf "tiny-screen/ffmpeg-core"
)

func main() {
	// encodeHls()
	serve()
}

func encodeHls() {
	x11 := tf.TinyFfmpegX11{}
	if err := x11.HlsEncode("input-samples/sample-video.mp4", "", true); err != nil {
		log.Fatal(err)
	}
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
