//go:build windows && amd64

package ffmpeg_tools

import (
	_ "embed"
	"fmt"
)

//go:embed bin/ffmpeg_windows_amd64.exe
var ffmpegWindowsx64Binary []byte

func init() {
	ffmpegBinary = ffmpegWindowsx64Binary
	fmt.Println("ffmpegBinary init windows amd64")
}
