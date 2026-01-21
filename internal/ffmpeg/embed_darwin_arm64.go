//go:build darwin && arm64

package ffmpeg

import _ "embed"

//go:embed bin/ffmpeg_darwin_arm64.gz
var ffmpegMacarm64CompressedBinary []byte

//go:embed bin/ffprobe_darwin_arm64.gz
var ffprobeMacarm64CompressedBinary []byte

func init() {
	ffmpegCompressedBinary = ffmpegMacarm64CompressedBinary
	ffprobeCompressedBinary = ffprobeMacarm64CompressedBinary
}
