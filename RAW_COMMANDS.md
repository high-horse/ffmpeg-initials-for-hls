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
