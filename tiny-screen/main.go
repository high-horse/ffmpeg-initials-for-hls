package main

import (
	"log"
	tf "tiny-screen/ffmpeg-core"
)

func main() {
	// Create an instance of TinyFfmpegX11
	x11 := tf.TinyFfmpegX11{}
	// log.Println(x11.Resolution()) 
	// return
	// err :=x11.RecordScreen()
	protocol := "udp"
	address := "127.0.0.1:1234" 
	err := x11.StartStream(protocol, address )
	if err != nil {
		log.Fatal(err)
	}
	
}
