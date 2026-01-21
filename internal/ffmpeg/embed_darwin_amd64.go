//go:build darwin && amd64

package ffmpeg

import _ "embed"

//go:embed bin/ffmpeg_darwin_amd64.gz
var ffmpegMacx64CompressedBinary []byte

//go:embed bin/ffprobe_darwin_amd64.gz
var ffprobeMacx64CompressedBinary []byte

func init() {
	ffmpegCompressedBinary = ffmpegMacx64CompressedBinary
	ffprobeCompressedBinary = ffprobeMacx64CompressedBinary
}
