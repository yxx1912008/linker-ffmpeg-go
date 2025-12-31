package ffmpeg_tools

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestNewFFmpeg 测试创建FFmpeg实例
func TestNewFFmpeg(t *testing.T) {
	// 定义进度回调
	callback := func(progress *Progress) {
		t.Logf("Progress: %.2f%%, Current: %dms, Total: %dms, Status: %s",
			progress.Percentage, progress.Current, progress.Total, progress.Status)
	}

	// 创建FFmpeg实例
	ffmpeg, err := NewFFmpeg(callback)
	if err != nil {
		t.Fatalf("Failed to create FFmpeg instance: %v", err)
	}

	// 检查回调是否设置
	if ffmpeg.Callback == nil {
		t.Fatal("Callback is nil")
	}

	// 路径可能为空，因为我们没有嵌入实际的FFmpeg二进制文件
	// 这是预期行为，用户可以后续设置路径
	t.Logf("FFmpeg instance created successfully, path: %s, extract path: %s", ffmpeg.FFmpegPath, ffmpeg.ExtractPath)
}

// TestNewFFmpegWithExtractPath 测试使用指定释放路径创建FFmpeg实例
func TestNewFFmpegWithExtractPath(t *testing.T) {
	// 定义进度回调
	callback := func(progress *Progress) {
		t.Logf("Progress: %.2f%%, Current: %dms, Total: %dms, Status: %s",
			progress.Percentage, progress.Current, progress.Total, progress.Status)
	}

	// 创建临时目录作为释放路径
	tempDir, err := os.MkdirTemp("", "ffmpeg_test_extract")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// 使用指定释放路径创建FFmpeg实例
	ffmpeg, err := NewFFmpegWithExtractPath(tempDir, callback)
	if err != nil {
		t.Fatalf("Failed to create FFmpeg instance with extract path: %v", err)
	}

	// 检查回调是否设置
	if ffmpeg.Callback == nil {
		t.Fatal("Callback is nil")
	}

	// 检查释放路径是否正确设置
	if ffmpeg.ExtractPath != tempDir {
		t.Fatalf("Expected extract path to be %s, got %s", tempDir, ffmpeg.ExtractPath)
	}

	// 验证FFmpeg路径是否包含指定的释放路径
	if ffmpeg.FFmpegPath != "" && !strings.Contains(ffmpeg.FFmpegPath, tempDir) {
		t.Fatalf("Expected FFmpeg path to contain %s, got %s", tempDir, ffmpeg.FFmpegPath)
	}

	t.Logf("FFmpeg instance created successfully with extract path, ffmpeg path: %s, extract path: %s", ffmpeg.FFmpegPath, ffmpeg.ExtractPath)
}

// TestNewFFmpegWithPath 测试使用指定路径创建FFmpeg实例
func TestNewFFmpegWithPath(t *testing.T) {
	// 定义进度回调
	callback := func(progress *Progress) {
		t.Logf("Progress: %.2f%%, Current: %dms, Total: %dms, Status: %s",
			progress.Percentage, progress.Current, progress.Total, progress.Status)
	}

	// 使用一个临时文件作为模拟的FFmpeg路径
	tempFile, err := os.CreateTemp("", "ffmpeg_test")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	tempPath := tempFile.Name()
	tempFile.Close()
	defer os.Remove(tempPath)

	// 创建FFmpeg实例
	ffmpeg, err := NewFFmpegWithPath(tempPath, callback)
	if err != nil {
		t.Fatalf("Failed to create FFmpeg instance with path: %v", err)
	}

	// 检查FFmpeg路径是否正确设置
	if ffmpeg.FFmpegPath != tempPath {
		t.Fatalf("Expected FFmpeg path to be %s, got %s", tempPath, ffmpeg.FFmpegPath)
	}

	// 检查回调是否设置
	if ffmpeg.Callback == nil {
		t.Fatal("Callback is nil")
	}

	t.Logf("FFmpeg instance created successfully with path: %s", ffmpeg.FFmpegPath)
}

// TestSetFFmpegPath 测试设置FFmpeg路径
func TestSetFFmpegPath(t *testing.T) {
	// 创建FFmpeg实例
	ffmpeg, err := NewFFmpeg(nil)
	if err != nil {
		t.Fatalf("Failed to create FFmpeg instance: %v", err)
	}

	// 设置新的FFmpeg路径
	newPath := "/custom/ffmpeg/path"
	ffmpeg.SetFFmpegPath(newPath)

	// 检查路径是否更新
	if ffmpeg.FFmpegPath != newPath {
		t.Fatalf("Expected FFmpeg path to be %s, got %s", newPath, ffmpeg.FFmpegPath)
	}

	t.Logf("FFmpeg path set successfully: %s", ffmpeg.FFmpegPath)
}

// TestSetProgressCallback 测试设置进度回调
func TestSetProgressCallback(t *testing.T) {
	// 创建FFmpeg实例
	ffmpeg, err := NewFFmpeg(nil)
	if err != nil {
		t.Fatalf("Failed to create FFmpeg instance: %v", err)
	}

	// 检查初始回调是否为nil
	if ffmpeg.Callback != nil {
		t.Fatal("Expected initial callback to be nil")
	}

	// 设置新的回调
	called := false
	newCallback := func(progress *Progress) {
		called = true
	}

	ffmpeg.SetProgressCallback(newCallback)

	// 检查回调是否更新
	if ffmpeg.Callback == nil {
		t.Fatal("Callback is still nil after setting")
	}

	// 触发回调
	ffmpeg.Callback(&Progress{
		Percentage: 50,
		Current:    1000,
		Total:      2000,
		Status:     "testing",
	})

	// 检查回调是否被调用
	if !called {
		t.Fatal("Callback was not called")
	}

	t.Log("Progress callback set and called successfully")
}

// TestGetVideoDuration 测试获取视频时长（模拟）
func TestGetVideoDuration(t *testing.T) {
	// 创建FFmpeg实例
	ffmpeg, err := NewFFmpeg(nil)
	if err != nil {
		t.Fatalf("Failed to create FFmpeg instance: %v", err)
	}

	// 注意：这个测试需要实际的视频文件才能通过
	// 这里只是测试函数调用，不会实际执行
	// 在实际使用中，应该提供一个测试视频文件路径
	t.Skip("Skipping TestGetVideoDuration - requires actual video file")

	// 替换为实际的视频文件路径
	videoPath := "test_video.mp4"
	duration, err := ffmpeg.GetVideoDuration(videoPath)
	if err != nil {
		t.Fatalf("Failed to get video duration: %v", err)
	}

	if duration <= 0 {
		t.Fatalf("Expected positive duration, got %d", duration)
	}

	t.Logf("Video duration: %dms", duration)
}

// TestExtractAudio 测试提取音频流（模拟）
func TestExtractAudio(t *testing.T) {
	// 创建FFmpeg实例
	ffmpeg, err := NewFFmpeg(func(progress *Progress) {
		t.Logf("ExtractAudio progress: %.2f%%", progress.Percentage)
	})
	if err != nil {
		t.Fatalf("Failed to create FFmpeg instance: %v", err)
	}

	// 注意：这个测试需要实际的视频文件才能通过
	// 这里只是测试函数结构，不会实际执行
	t.Skip("Skipping TestExtractAudio - requires actual video file")

	// 替换为实际的视频文件路径和输出路径
	params := &ExtractAudioParams{
		InputPath:  "test_video.mp4",
		OutputPath: "test_audio.mp3",
	}

	err = ffmpeg.ExtractAudio(params)
	if err != nil {
		t.Fatalf("Failed to extract audio: %v", err)
	}

	t.Log("Audio extracted successfully")
}

// TestSplitVideo 测试视频分段（模拟）
func TestSplitVideo(t *testing.T) {
	// 创建FFmpeg实例
	ffmpeg, err := NewFFmpeg(func(progress *Progress) {
		t.Logf("SplitVideo progress: %.2f%%", progress.Percentage)
	})
	if err != nil {
		t.Fatalf("Failed to create FFmpeg instance: %v", err)
	}

	// 注意：这个测试需要实际的视频文件才能通过
	// 这里只是测试函数结构，不会实际执行
	t.Skip("Skipping TestSplitVideo - requires actual video file")

	// 创建临时输出目录
	outputDir := filepath.Join(os.TempDir(), "ffmpeg_test_split")
	os.RemoveAll(outputDir) // 清理之前的测试数据

	// 替换为实际的视频文件路径
	params := &SplitVideoParams{
		InputPath:    "test_video.mp4",
		OutputDir:    outputDir,
		SegmentTime:  10, // 10秒分段
		OutputPrefix: "segment_",
	}

	segments, err := ffmpeg.SplitVideo(params)
	if err != nil {
		t.Fatalf("Failed to split video: %v", err)
	}

	if len(segments) == 0 {
		t.Fatal("No segments created")
	}

	t.Logf("Video split into %d segments", len(segments))
}

// TestExtractKeyFrames 测试提取关键帧（模拟）
func TestExtractKeyFrames(t *testing.T) {
	// 创建FFmpeg实例
	ffmpeg, err := NewFFmpeg(func(progress *Progress) {
		t.Logf("ExtractKeyFrames progress: %.2f%%", progress.Percentage)
	})
	if err != nil {
		t.Fatalf("Failed to create FFmpeg instance: %v", err)
	}

	// 注意：这个测试需要实际的视频文件才能通过
	// 这里只是测试函数结构，不会实际执行
	t.Skip("Skipping TestExtractKeyFrames - requires actual video file")

	// 创建临时输出目录
	outputDir := filepath.Join(os.TempDir(), "ffmpeg_test_keyframes")
	os.RemoveAll(outputDir) // 清理之前的测试数据

	// 替换为实际的视频文件路径
	params := &ExtractKeyFramesParams{
		InputPath:     "test_video.mp4",
		OutputDir:     outputDir,
		FrameInterval: 5, // 每5秒提取一个关键帧
		OutputPrefix:  "keyframe_",
	}

	keyframes, err := ffmpeg.ExtractKeyFrames(params)
	if err != nil {
		t.Fatalf("Failed to extract keyframes: %v", err)
	}

	if len(keyframes) == 0 {
		t.Fatal("No keyframes created")
	}

	t.Logf("Extracted %d keyframes", len(keyframes))
}

// TestParseProgress 测试解析进度信息
func TestParseProgress(t *testing.T) {
	// 创建FFmpeg实例
	ffmpeg, err := NewFFmpeg(nil)
	if err != nil {
		t.Fatalf("Failed to create FFmpeg instance: %v", err)
	}

	// 测试进度回调
	progressCalled := false
	ffmpeg.SetProgressCallback(func(progress *Progress) {
		progressCalled = true
		if progress.Percentage < 0 || progress.Percentage > 100 {
			t.Fatalf("Invalid percentage: %.2f", progress.Percentage)
		}
		if progress.Current < 0 {
			t.Fatalf("Invalid current time: %d", progress.Current)
		}
		if progress.Total < 0 {
			t.Fatalf("Invalid total time: %d", progress.Total)
		}
		t.Logf("Parsed progress: %.2f%%, Current: %dms, Total: %dms, Status: %s",
			progress.Percentage, progress.Current, progress.Total, progress.Status)
	})

	// 模拟ffmpeg输出
	testOutput := "Duration: 00:05:00.00, start: 0.000000, bitrate: 1024 kb/s\n"
	testOutput += "frame=  100 fps= 25 q=24.0 size=N/A time=00:00:04.00 bitrate=N/A speed=1.0x\n"

	// 解析进度
	ffmpeg.parseProgress(testOutput)

	if !progressCalled {
		t.Fatal("Progress callback was not called")
	}

	t.Log("Progress parsing test passed")
}
