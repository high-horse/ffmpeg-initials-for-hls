# ffmpeg

# todo HLS server over http 

### Basic HLS Export (VOD – Video on Demand) (Basic ffmpeg command)
```bash
ffmpeg -i input.mp4 \
       -codec: copy \
       -start_number 0 \
       -hls_time 10 \
       -hls_list_size 0 \
       -f hls output.m3u8

```

### Adaptive Bitrate HLS (Multiple Qualities)
```bash
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
/project/
├── main.go
└── hls/
    ├── master.m3u8
    ├──
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