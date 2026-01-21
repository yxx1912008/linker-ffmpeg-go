//go:build windows && amd64

package ffmpeg

import (
	_ "embed"
)

//go:embed bin/ffmpeg_windows_amd64.exe.gz
var ffmpegWindowsx64CompressedBinary []byte

//go:embed bin/ffprobe_windows_amd64.exe.gz
var ffprobeWindowsx64CompressedBinary []byte

func init() {
	ffmpegCompressedBinary = ffmpegWindowsx64CompressedBinary
	ffprobeCompressedBinary = ffprobeWindowsx64CompressedBinary
}
