# ffmpeg

# todo HLS server over http 

### Basic HLS Export (VOD – Video on Demand) (Basic ffmpeg command)
```bash
# correct
ffmpeg -i samples/sample-video.mp4 -c copy -hls_time 10 -hls_list_size 0 -start_number 0 -hls_segment_filename outputs/hls/segment_%03d.ts -f hls outputs/hls/sample-hls-op-m3u8

# incorrect
ffmpeg -i input.mp4 \
       -codec: copy \
       -start_number 0 \
       -hls_time 10 \
       -hls_list_size 0 \
       -f hls output.m3u8

```

### Adaptive Bitrate HLS (Multiple Qualities)
```bash
# correct
ffmpeg -i samples/sample-video.mp4 \
  -filter_complex "[0:v]split=3[v1][v2][v3]; \
                   [v1]scale=w=1920:h=1080[vout1]; \
                   [v2]scale=w=1280:h=720[vout2]; \
                   [v3]scale=w=854:h=480[vout3]" \
  -map "[vout1]" -map 0:a -c:v:0 libx264 -b:v:0 3000k -c:a:0 aac -f hls \
    -hls_time 6 -hls_playlist_type vod -hls_segment_filename outputs/hls-adaptive/1080p_%03d.ts outputs/hls-adaptive/1080p.m3u8 \
  -map "[vout2]" -map 0:a -c:v:1 libx264 -b:v:1 1500k -c:a:1 aac -f hls \
    -hls_time 6 -hls_playlist_type vod -hls_segment_filename outputs/hls-adaptive/720p_%03d.ts outputs/hls-adaptive/720p.m3u8 \
  -map "[vout3]" -map 0:a -c:v:2 libx264 -b:v:2 800k  -c:a:2 aac -f hls \
    -hls_time 6 -hls_playlist_type vod -hls_segment_filename outputs/hls-adaptive/480p_%03d.ts outputs/hls-adaptive/480p.m3u8

# incorrect
ffmpeg -i input.mp4 \
  -map 0:v -map 0:a -c:v libx264 -c:a aac -b:v:0 3000k -s:v:0 1920x1080 -hls_time 6 -hls_segment_filename "1080p_%03d.ts" -f hls 1080p.m3u8 \
  -map 0:v -map 0:a -c:v libx264 -c:a aac -b:v:1 1500k -s:v:1 1280x720  -hls_time 6 -hls_segment_filename "720p_%03d.ts"  -f hls 720p.m3u8 \
  -map 0:v -map 0:a -c:v libx264 -c:a aac -b:v:2 800k  -s:v:2 854x480   -hls_time 6 -hls_segment_filename "480p_%03d.ts"  -f hls 480p.m3u8
```

** Output dir **
```bash
/project/
├── main.go
└── hls/
    ├── master.m3u8
    ├──

```

** Server **
```bash
package main

import (
	"log"
	"net/http"
)

func main() {
	// Serve files from the "hls" directory
	fs := http.FileServer(http.Dir("./hls"))

	// Route all /hls/ URLs to the hls directory
	http.Handle("/hls/", http.StripPrefix("/hls/", fs))

	// Optional: simple root handler
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html") // or just say hello
	})

	// Start the server
	port := ":8080"
	log.Println("Serving HLS on http://localhost" + port + "/hls/master.m3u8")
	log.Fatal(http.ListenAndServe(port, nil))
}
```

** client **
```bash
<!DOCTYPE html>
<html>
<head>
  <meta charset="UTF-8">
  <title>HLS Streaming</title>
</head>
<body>
  <h1>HLS Demo</h1>
  <video id="video" controls width="640"></video>
  <script src="https://cdn.jsdelivr.net/npm/hls.js@latest"></script>
  <script>
    const video = document.getElementById('video');
    const videoSrc = '/hls/master.m3u8';

    if (Hls.isSupported()) {
      const hls = new Hls();
      hls.loadSource(videoSrc);
      hls.attachMedia(video);
    } else if (video.canPlayType('application/vnd.apple.mpegurl')) {
      video.src = videoSrc;
    }
  </script>
</body>
</html>
```

**  Optional: Web Player (HTML5 + hls.js) **
```bash
<!DOCTYPE html>
<html>
<head>
  <meta charset="UTF-8">
  <title>HLS Streaming</title>
</head>
<body>
  <h1>HLS Demo</h1>
  <video id="video" controls width="640"></video>
  <script src="https://cdn.jsdelivr.net/npm/hls.js@latest"></script>
  <script>
    const video = document.getElementById('video');
    const videoSrc = '/hls/master.m3u8';

    if (Hls.isSupported()) {
      const hls = new Hls();
      hls.loadSource(videoSrc);
      hls.attachMedia(video);
    } else if (video.canPlayType('application/vnd.apple.mpegurl')) {
      video.src = videoSrc;
    }
  </script>
</body>
</html>
```