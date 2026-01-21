//go:build windows && amd64

package ffmpeg

import (
	_ "embed"
)

//go:embed bin/ffmpeg_windows_amd64.exe
var ffmpegWindowsx64Binary []byte

//go:embed bin/ffprobe_windows_amd64.exe
var ffprobeWindowsx64Binary []byte

func init() {
	ffmpegBinary = ffmpegWindowsx64Binary
	ffprobeBinary = ffprobeWindowsx64Binary
}
