# linker-ffmpeg-go

一个基于FFmpeg的Go语言封装库，提供简单易用的视频处理API，支持音频提取、视频分段和关键帧提取等功能。

## 功能特性

- **音频提取**：从视频文件中提取音频流
- **视频分段**：将视频文件分割为指定时长的小段
- **关键帧提取**：按时间间隔提取视频关键帧
- **视频时长获取**：获取视频文件的总时长
- **跨平台支持**：内置多种平台的FFmpeg二进制文件
- **进度回调**：实时获取处理进度

## 安装方法

```bash
go get -u github.com/yxx1912008/linker-ffmpeg-go/pkg/ffmpeg
```

## 快速开始

```go
package main

import (
	"fmt"
	"os"
	"path/filepath"

	ffmpeg "github.com/yxx1912008/linker-ffmpeg-go/pkg/ffmpeg"
)

func main() {
	// 创建FFmpeg实例，设置进度回调
	ffmpegInstance, err := ffmpeg.NewFFmpeg(func(progress *ffmpeg.Progress) {
		fmt.Printf("Progress: %.2f%%\n", progress.Percentage)
	})
	if err != nil {
		panic(err)
	}

	// 使用示例：视频分段
	segmentParams := &ffmpeg.SplitVideoParams{
		InputPath:    "input.mp4",
		OutputDir:    "/tmp/segments",
		SegmentTime:  10, // 每10秒一个分段
		OutputPrefix: "segment_",
	}

	segments, err := ffmpegInstance.SplitVideo(segmentParams)
	if err != nil {
		fmt.Printf("Failed to split video: %v\n", err)
		return
	}

	fmt.Printf("Split video into %d segments\n", len(segments))
}
```

## API参考

### 主要类型

- `FFmpeg`：FFmpeg工具实例
- `Progress`：进度信息结构体
- `ExtractAudioParams`：音频提取参数
- `SplitVideoParams`：视频分段参数
- `ExtractKeyFramesParams`：关键帧提取参数

### 主要方法

#### 创建实例
- `NewFFmpeg(callback ProgressCallback) (*FFmpeg, error)`：创建FFmpeg实例，使用默认临时目录
- `NewFFmpegWithExtractPath(extractPath string, callback ProgressCallback) (*FFmpeg, error)`：使用指定释放路径创建实例
- `NewFFmpegWithPath(ffmpegPath string, callback ProgressCallback) (*FFmpeg, error)`：使用指定FFmpeg路径创建实例

#### 视频处理
- `ExtractAudio(params *ExtractAudioParams) error`：提取音频流
- `SplitVideo(params *SplitVideoParams) ([]string, error)`：视频分段
- `ExtractKeyFrames(params *ExtractKeyFramesParams) ([]string, error)`：提取关键帧
- `GetVideoDuration(inputPath string) (int64, error)`：获取视频时长

## 示例代码

### 1. 音频提取

```go
params := &ffmpeg.ExtractAudioParams{
	InputPath:  "input.mp4",
	OutputPath: "output.mp3",
}

err := ffmpegInstance.ExtractAudio(params)
if err != nil {
	fmt.Printf("Failed to extract audio: %v\n", err)
}
```

### 2. 视频分段

```go
params := &ffmpeg.SplitVideoParams{
	InputPath:    "input.mp4",
	OutputDir:    "/tmp/segments",
	SegmentTime:  10, // 每10秒一个分段
	OutputPrefix: "segment_",
}

segments, err := ffmpegInstance.SplitVideo(params)
if err != nil {
	fmt.Printf("Failed to split video: %v\n", err)
}
```

### 3. 关键帧提取

```go
params := &ffmpeg.ExtractKeyFramesParams{
	InputPath:     "input.mp4",
	OutputDir:     "/tmp/keyframes",
	FrameInterval: 5, // 每5秒提取一个关键帧
	OutputPrefix:  "keyframe_",
}

keyframes, err := ffmpegInstance.ExtractKeyFrames(params)
if err != nil {
	fmt.Printf("Failed to extract keyframes: %v\n", err)
}
```

## 注意事项

1. 首次使用时，库会自动提取FFmpeg二进制文件到指定目录或临时目录
2. 处理大文件时，建议设置适当的进度回调，以便监控处理进度
3. 确保输入文件路径正确，并且有足够的权限访问和写入输出目录
4. 不同平台的FFmpeg二进制文件已内置，无需额外安装

## 许可证

MIT License
