package ffmpeg

// ffmpegBinary 存储当前平台的FFmpeg二进制数据
var ffmpegBinary []byte

// GetFFmpegBinary 获取当前平台的FFmpeg二进制数据
func GetFFmpegBinary() []byte {
	return ffmpegBinary
}