package main

import (
	"fmt"
	"os"
	"path/filepath"

	tools "github.com/yxx1912008/linker-ffmpeg-go/ffmpeg"
)

func main() {
	ffmpeg, err := tools.NewFFmpegWithExtractPath("D://", func(progress *tools.Progress) {
		fmt.Printf("Progress: %.2f%%\n", progress.Percentage)
	})
	if err != nil {
		panic(err)
	}

	// 主输出目录
	mainOutputDir := filepath.Join("D://output", "ffmpeg_test")
	os.RemoveAll(mainOutputDir) // 清理之前的测试数据

	// 1. 视频分段
	segmentOutputDir := filepath.Join(mainOutputDir, "segments")
	segmentParams := &tools.SplitVideoParams{
		InputPath:    "input.mp4",
		OutputDir:    segmentOutputDir,
		SegmentTime:  10, // 每10秒一个分段
		OutputPrefix: "segment_",
	}

	segments, err := ffmpeg.SplitVideo(segmentParams)
	if err != nil {
		fmt.Printf("Failed to split video: %v\n", err)
		return
	}

	fmt.Printf("Split video into %d segments\n", len(segments))
	fmt.Printf("Segments: %v\n", segments)

	// 2. 对每个分段提取关键帧
	for i, segmentPath := range segments {
		// 为每个分段创建独立的关键帧输出目录
		keyframesOutputDir := filepath.Join(mainOutputDir, fmt.Sprintf("keyframes_segment_%03d", i))

		// 提取关键帧
		keyframeParams := &tools.ExtractKeyFramesParams{
			InputPath:     segmentPath,
			OutputDir:     keyframesOutputDir,
			FrameInterval: 1, // 每1秒提取一个关键帧
			OutputPrefix:  fmt.Sprintf("keyframe_segment_%03d_", i),
		}

		keyframes, err := ffmpeg.ExtractKeyFrames(keyframeParams)
		if err != nil {
			fmt.Printf("Failed to extract keyframes for segment %d: %v\n", i, err)
			continue
		}

		if len(keyframes) == 0 {
			fmt.Printf("No keyframes created for segment %d\n", i)
		} else {
			fmt.Printf("Extracted %d keyframes for segment %d\n", len(keyframes), i)
			fmt.Printf("Keyframes for segment %d: %v\n", i, keyframes)
		}
	}

	fmt.Println("All tasks completed successfully")
}
