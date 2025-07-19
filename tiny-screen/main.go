package main

import (
	tf "tiny-screen/ffmpeg-core"
)

func main() {
	// Create an instance of TinyFfmpegX11
	x11 := tf.TinyFfmpegX11{}
	x11.RecordScreen()
}
