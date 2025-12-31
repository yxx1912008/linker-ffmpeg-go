//go:build darwin && arm64

package ffmpeg

import _ "embed"

//go:embed bin/ffmpeg_darwin_arm64
var ffmpegMacarm64Binary []byte

func init() {
	ffmpegBinary = ffmpegMacarm64Binary
}
