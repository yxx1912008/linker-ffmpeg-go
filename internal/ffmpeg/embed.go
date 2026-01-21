package ffmpeg

// ffmpegBinary 存储当前平台的FFmpeg二进制数据
var ffmpegBinary []byte

// ffprobeBinary 存储当前平台的FFprobe二进制数据
var ffprobeBinary []byte

// GetFFmpegBinary 获取当前平台的FFmpeg二进制数据
func GetFFmpegBinary() []byte {
	return ffmpegBinary
}

// GetFFprobeBinary 获取当前平台的FFprobe二进制数据
func GetFFprobeBinary() []byte {
	return ffprobeBinary
}
