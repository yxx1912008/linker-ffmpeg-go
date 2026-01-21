package ffmpeg

import (
	"bytes"
	"compress/gzip"
	"io"
	"sync"
)

// ffmpegCompressedBinary 存储当前平台的FFmpeg压缩二进制数据
var ffmpegCompressedBinary []byte

// ffprobeCompressedBinary 存储当前平台的FFprobe压缩二进制数据
var ffprobeCompressedBinary []byte

// 缓存解压后的二进制数据，避免重复解压
var (
	ffmpegDecompressedBinary  []byte
	ffprobeDecompressedBinary []byte
	ffmpegOnce                sync.Once
	ffprobeOnce               sync.Once
)

// decompressGzip 解压gzip压缩的数据
func decompressGzip(compressed []byte) ([]byte, error) {
	if len(compressed) == 0 {
		return nil, nil
	}

	reader, err := gzip.NewReader(bytes.NewReader(compressed))
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	return io.ReadAll(reader)
}

// GetFFmpegBinary 获取当前平台的FFmpeg二进制数据（解压后）
func GetFFmpegBinary() []byte {
	ffmpegOnce.Do(func() {
		var err error
		ffmpegDecompressedBinary, err = decompressGzip(ffmpegCompressedBinary)
		if err != nil {
			ffmpegDecompressedBinary = nil
		}
	})
	return ffmpegDecompressedBinary
}

// GetFFprobeBinary 获取当前平台的FFprobe二进制数据（解压后）
func GetFFprobeBinary() []byte {
	ffprobeOnce.Do(func() {
		var err error
		ffprobeDecompressedBinary, err = decompressGzip(ffprobeCompressedBinary)
		if err != nil {
			ffprobeDecompressedBinary = nil
		}
	})
	return ffprobeDecompressedBinary
}
