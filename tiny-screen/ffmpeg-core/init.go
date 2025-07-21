package ffmpegcore


type TinyFfmpegCore interface{
	RecordScreen() error
	StartStream(protocol, address string) error
}

type FFmpegCore interface {
	RecordScreen(outputFile string, display string, offsetX int, offsetY int, width int, height int) error
	RecordAudio(outputFile string, device string) error
	RecordScreenWithAudio(outputFile string, display string, offsetX int, offsetY int, width int, height int, audioDevice string) error
	RecordScreenWithAudioAndFramerate(outputFile string, display string, offsetX int, offsetY int, width int, height int, audioDevice string, framerate int) error
	RecordScreenWithAudioAndFramerateAndFormat(outputFile string, display string, offsetX int, offsetY int, width int, height int, audioDevice string, framerate int, format string) error
	RecordScreenWithAudioAndFramerateAndFormatAndCodec(outputFile string, display string, offsetX int, offsetY int, width int, height int, audioDevice string, framerate int, format string, codec string) error
	RecordScreenWithAudioAndFramerateAndFormatAndCodecAndBitrate(outputFile string, display string, offsetX int, offsetY int, width int, height int, audioDevice string, framerate int, format string, codec string, bitrate string) error
	RecordScreenWithAudioAndFramerateAndFormatAndCodecAndBitrateAndQuality(outputFile string, display string, offsetX int, offsetY int, width int, height int, audioDevice string, framerate int, format string, codec string, bitrate string, quality int) error
	RecordScreenWithAudioAndFramerateAndFormatAndCodecAndBitrateAndQualityAndPreset(outputFile string, display string, offsetX int, offsetY int, width int, height int, audioDevice string, framerate int, format string, codec string, bitrate string, quality int, preset string) error
	RecordScreenWithAudioAndFramerateAndFormatAndCodecAndBitrateAndQualityAndPresetAndThreads(outputFile string, display string, offsetX int, offsetY int, width int, height int, audioDevice string, framerate int, format string, codec string, bitrate string, quality int, preset string, threads int) error
	RecordScreenWithAudioAndFramerateAndFormatAndCodecAndBitrateAndQualityAndPresetAndThreadsAndLogLevel(outputFile string, display string, offsetX int, offsetY int, width int, height int, audioDevice string, framerate int, format string, codec string, bitrate string, quality int, preset string, threads int, logLevel string) error
	RecordScreenWithAudioAndFramerateAndFormatAndCodecAndBitrateAndQualityAndPresetAndThreadsAndLogLevelAndExtraArgs(outputFile string, display string, offsetX int, offsetY int, width int, height int, audioDevice string, framerate int, format string, codec string, bitrate string, quality int, preset string, threads int, logLevel string, extraArgs []string) error
	RecordScreenWithAudioAndFramerateAndFormatAndCodecAndBitrateAndQualityAndPresetAndThreadsAndLogLevelAndExtraArgsAndEnv(outputFile string, display string, offsetX int, offsetY int, width int, height int, audioDevice string, framerate int, format string, codec string, bitrate string, quality int, preset string, threads int, logLevel string, extraArgs []string, env map[string]string) error
	RecordScreenWithAudioAndFramerateAndFormatAndCodecAndBitrateAndQualityAndPresetAndThreadsAndLogLevelAndExtraArgsAndEnvAndTimeout(outputFile string, display string, offsetX int, offsetY int, width int, height int, audioDevice string, framerate int, format string, codec string, bitrate string, quality int, preset string, threads int, logLevel string, extraArgs []string, env map[string]string, timeout int) error
	RecordScreenWithAudioAndFramerateAndFormatAndCodecAndBitrateAndQualityAndPresetAndThreadsAndLogLevelAndExtraArgsAndEnvAndTimeoutAndProgressCallback(outputFile string, display string, offsetX int, offsetY int, width int, height int, audioDevice string, framerate int, format string, codec string, bitrate string, quality int, preset string, threads int, logLevel string, extraArgs []string, env map[string]string, timeout int, progressCallback func(progress float64)) error
}