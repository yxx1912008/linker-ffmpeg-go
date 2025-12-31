//go:build darwin && amd64

package ffmpeg

import _ "embed"

//go:embed bin/ffmpeg_darwin_amd64
var ffmpegMacx64Binary []byte

func init() {
	ffmpegBinary = ffmpegMacx64Binary
}
