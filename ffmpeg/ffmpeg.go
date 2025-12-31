package ffmpeg_tools

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"

	ffmpeg_go "github.com/u2takey/ffmpeg-go"
)

var ffmpegBinary []byte

// extractFFmpeg 提取FFmpeg二进制文件到指定路径
// 如果extractPath为空，则使用临时目录
func extractFFmpeg(extractPath string) (string, error) {
	// 确定释放目录
	var outputDir string
	if extractPath != "" {
		outputDir = extractPath
		// 确保释放目录存在
		if err := os.MkdirAll(outputDir, 0755); err != nil {
			return "", fmt.Errorf("failed to create extract directory: %w", err)
		}
	} else {
		// 使用临时目录
		outputDir = os.TempDir()
	}

	// 根据平台获取FFmpeg文件名
	ffmpegFileName := "ffmpeg"
	if runtime.GOOS == "windows" {
		ffmpegFileName += ".exe"
	}
	// 构建输出路径
	outputPath := filepath.Join(outputDir, ffmpegFileName)

	// 检查文件是否已存在
	if _, err := os.Stat(outputPath); err == nil {
		// 文件已存在，返回路径
		return outputPath, nil
	}

	// 直接使用已经初始化的ffmpegBinary变量
	if len(ffmpegBinary) == 0 {
		// 如果二进制数据为空，返回空路径，让用户自行设置
		return "", nil
	}

	// 创建输出文件
	outFile, err := os.Create(outputPath)
	if err != nil {
		return "", fmt.Errorf("failed to create output file: %w", err)
	}
	defer outFile.Close()

	// 写入二进制数据
	_, err = outFile.Write(ffmpegBinary)
	if err != nil {
		return "", fmt.Errorf("failed to write file content: %w", err)
	}

	// 设置执行权限
	if runtime.GOOS != "windows" {
		if err := outFile.Chmod(0755); err != nil {
			return "", fmt.Errorf("failed to set executable permission: %w", err)
		}
	}

	return outputPath, nil
}

// NewFFmpeg 创建并初始化FFmpeg工具
// 保持向后兼容，extractPath默认为空（使用临时目录）
func NewFFmpeg(callback ProgressCallback) (*FFmpeg, error) {
	return NewFFmpegWithExtractPath("", callback)
}

// NewFFmpegWithExtractPath 使用指定释放路径创建FFmpeg工具
// extractPath: FFmpeg二进制文件释放路径，为空则使用临时目录
func NewFFmpegWithExtractPath(extractPath string, callback ProgressCallback) (*FFmpeg, error) {
	// 提取FFmpeg到指定路径
	ffmpegPath, err := extractFFmpeg(extractPath)
	if err != nil {
		return nil, err
	}

	return &FFmpeg{
		FFmpegPath:  ffmpegPath,
		ExtractPath: extractPath,
		Callback:    callback,
	}, nil
}

// NewFFmpegWithPath 使用指定路径创建FFmpeg工具
func NewFFmpegWithPath(ffmpegPath string, callback ProgressCallback) (*FFmpeg, error) {
	// 检查指定的FFmpeg路径是否存在
	if _, err := os.Stat(ffmpegPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("ffmpeg executable not found at: %s", ffmpegPath)
	}

	return &FFmpeg{
		FFmpegPath: ffmpegPath,
		Callback:   callback,
	}, nil
}

// 设置FFmpeg路径
func (f *FFmpeg) SetFFmpegPath(path string) {
	f.FFmpegPath = path
}

// 设置进度回调
func (f *FFmpeg) SetProgressCallback(callback ProgressCallback) {
	f.Callback = callback
}

// 解析ffmpeg输出的进度信息
func (f *FFmpeg) parseProgress(output string) {
	if f.Callback == nil {
		return
	}

	// 匹配时间进度信息
	timeRegex := regexp.MustCompile(`time=(\d+):(\d+):(\d+)\.(\d+)`)
	durationRegex := regexp.MustCompile(`Duration: (\d+):(\d+):(\d+)\.(\d+)`)

	timeMatches := timeRegex.FindStringSubmatch(output)
	durationMatches := durationRegex.FindStringSubmatch(output)

	if len(timeMatches) < 5 || len(durationMatches) < 5 {
		return
	}

	// 解析当前时间
	currentHours, _ := strconv.Atoi(timeMatches[1])
	currentMinutes, _ := strconv.Atoi(timeMatches[2])
	currentSeconds, _ := strconv.Atoi(timeMatches[3])
	currentMilliseconds, _ := strconv.Atoi(timeMatches[4])

	current := int64(currentHours*3600*1000 + currentMinutes*60*1000 + currentSeconds*1000 + currentMilliseconds*10)

	// 解析总时长
	totalHours, _ := strconv.Atoi(durationMatches[1])
	totalMinutes, _ := strconv.Atoi(durationMatches[2])
	totalSeconds, _ := strconv.Atoi(durationMatches[3])
	totalMilliseconds, _ := strconv.Atoi(durationMatches[4])

	total := int64(totalHours*3600*1000 + totalMinutes*60*1000 + totalSeconds*1000 + totalMilliseconds*10)

	// 计算百分比
	percentage := 0.0
	if total > 0 {
		percentage = (float64(current) / float64(total)) * 100
		if percentage > 100 {
			percentage = 100
		}
	}

	// 调用回调函数
	f.Callback(&Progress{
		Percentage: percentage,
		Current:    current,
		Total:      total,
		Status:     "processing",
	})
}

// ExtractAudio 提取音频流
func (f *FFmpeg) ExtractAudio(params *ExtractAudioParams) error {
	// 设置环境变量指定ffmpeg路径
	origFfmpegPath := os.Getenv("FFMPEG_PATH")
	os.Setenv("FFMPEG_PATH", f.FFmpegPath)
	defer os.Setenv("FFMPEG_PATH", origFfmpegPath)

	// 创建命令
	cmd := exec.Command(f.FFmpegPath, "-i", params.InputPath, "-vn", "-acodec", "copy", params.OutputPath)

	// 处理输出
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}

	// 启动命令
	if err := cmd.Start(); err != nil {
		return err
	}

	// 读取输出并解析进度
	outputChan := make(chan string)

	// 读取stderr
	go func() {
		buf := make([]byte, 1024)
		for {
			n, err := stderr.Read(buf)
			if err != nil {
				close(outputChan)
				return
			}
			if n > 0 {
				outputChan <- string(buf[:n])
			}
		}
	}()

	// 解析输出
	go func() {
		for output := range outputChan {
			f.parseProgress(output)
		}
	}()

	// 等待命令完成
	if err := cmd.Wait(); err != nil {
		return err
	}

	// 发送完成进度
	if f.Callback != nil {
		f.Callback(&Progress{
			Percentage: 100,
			Status:     "completed",
		})
	}

	return nil
}

// SplitVideo 视频分段
func (f *FFmpeg) SplitVideo(params *SplitVideoParams) ([]string, error) {
	// 设置环境变量指定ffmpeg路径
	origFfmpegPath := os.Getenv("FFMPEG_PATH")
	os.Setenv("FFMPEG_PATH", f.FFmpegPath)
	defer os.Setenv("FFMPEG_PATH", origFfmpegPath)

	// 确保输出目录存在
	if err := os.MkdirAll(params.OutputDir, 0755); err != nil {
		return nil, err
	}

	// 构建输出文件名模式
	// 注意：segment muxer只支持单个占位符，我们先使用序号
	outputPattern := fmt.Sprintf("%s/%s%%03d.mp4", params.OutputDir, params.OutputPrefix)

	// 创建命令
	cmd := exec.Command(f.FFmpegPath, "-i", params.InputPath, "-c", "copy", "-f", "segment", "-segment_time", strconv.Itoa(params.SegmentTime), "-reset_timestamps", "1", outputPattern)

	// 处理输出
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, err
	}

	// 启动命令
	if err := cmd.Start(); err != nil {
		return nil, err
	}

	// 读取输出并解析进度
	outputChan := make(chan string)
	var stderrOutput strings.Builder

	go func() {
		buf := make([]byte, 1024)
		for {
			n, err := stderr.Read(buf)
			if err != nil {
				close(outputChan)
				return
			}
			if n > 0 {
				output := string(buf[:n])
				stderrOutput.WriteString(output)
				outputChan <- output
			}
		}
	}()

	// 解析输出
	go func() {
		for output := range outputChan {
			f.parseProgress(output)
		}
	}()

	// 等待命令完成
	if err := cmd.Wait(); err != nil {
		return nil, fmt.Errorf("ffmpeg command failed: %w\n%s", err, stderrOutput.String())
	}

	// 发送完成进度
	if f.Callback != nil {
		f.Callback(&Progress{
			Percentage: 100,
			Status:     "completed",
		})
	}

	// 获取分段文件列表
	files, err := os.ReadDir(params.OutputDir)
	if err != nil {
		return nil, err
	}

	var segmentFiles []string
	for _, file := range files {
		if !file.IsDir() && strings.HasPrefix(file.Name(), params.OutputPrefix) && strings.HasSuffix(file.Name(), ".mp4") {
			segmentFiles = append(segmentFiles, fmt.Sprintf("%s/%s", params.OutputDir, file.Name()))
		}
	}

	// 为分段文件添加时间戳信息
	for i, segmentPath := range segmentFiles {
		// 计算时间戳（秒级）
		timestamp := i * params.SegmentTime

		// 获取文件目录
		dir := filepath.Dir(segmentPath)

		// 生成新文件名，包含序号和时间戳
		newName := fmt.Sprintf("%s%03d_%d.mp4", params.OutputPrefix, i, timestamp)
		newPath := filepath.Join(dir, newName)

		// 重命名文件
		if err := os.Rename(segmentPath, newPath); err != nil {
			return nil, fmt.Errorf("failed to rename segment file: %w", err)
		}

		// 更新文件列表中的路径
		segmentFiles[i] = newPath
	}

	return segmentFiles, nil
}

// ExtractKeyFrames 提取关键帧
func (f *FFmpeg) ExtractKeyFrames(params *ExtractKeyFramesParams) ([]string, error) {
	// 设置环境变量指定ffmpeg路径
	origFfmpegPath := os.Getenv("FFMPEG_PATH")
	os.Setenv("FFMPEG_PATH", f.FFmpegPath)
	defer os.Setenv("FFMPEG_PATH", origFfmpegPath)

	// 确保输出目录存在
	if err := os.MkdirAll(params.OutputDir, 0755); err != nil {
		return nil, err
	}

	// 构建输出文件名模式
	outputPattern := fmt.Sprintf("%s/%s%%06d.jpg", params.OutputDir, params.OutputPrefix)

	// 使用ffmpeg的select过滤器结合fps过滤器，实现按时间间隔提取关键帧
	// select='eq(pict_type,I)' 表示只选择关键帧
	// fps=1/%d 表示每%d秒输出1帧
	cmd := exec.Command(f.FFmpegPath, "-i", params.InputPath,
		"-vf", fmt.Sprintf("select='eq(pict_type,I)',fps=1/%d", params.FrameInterval),
		"-vsync", "vfr",
		outputPattern)

	// 处理输出
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, err
	}

	// 启动命令
	if err := cmd.Start(); err != nil {
		return nil, err
	}

	// 读取输出并解析进度
	outputChan := make(chan string)
	var stderrOutput strings.Builder

	go func() {
		buf := make([]byte, 1024)
		for {
			n, err := stderr.Read(buf)
			if err != nil {
				close(outputChan)
				return
			}
			if n > 0 {
				output := string(buf[:n])
				stderrOutput.WriteString(output)
				outputChan <- output
			}
		}
	}()

	// 解析输出
	go func() {
		for output := range outputChan {
			f.parseProgress(output)
		}
	}()

	// 等待命令完成
	if err := cmd.Wait(); err != nil {
		return nil, fmt.Errorf("ffmpeg command failed: %w\n%s", err, stderrOutput.String())
	}

	// 发送完成进度
	if f.Callback != nil {
		f.Callback(&Progress{
			Percentage: 100,
			Status:     "completed",
		})
	}

	// 获取生成的关键帧文件列表
	var keyFrameFiles []string
	files, err := os.ReadDir(params.OutputDir)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if !file.IsDir() && strings.HasPrefix(file.Name(), params.OutputPrefix) && strings.HasSuffix(file.Name(), ".jpg") {
			keyFrameFiles = append(keyFrameFiles, fmt.Sprintf("%s/%s", params.OutputDir, file.Name()))
		}
	}

	return keyFrameFiles, nil
}

// GetVideoDuration 获取视频时长
func (f *FFmpeg) GetVideoDuration(inputPath string) (int64, error) {
	// 设置环境变量指定ffmpeg路径
	origFfmpegPath := os.Getenv("FFMPEG_PATH")
	os.Setenv("FFMPEG_PATH", f.FFmpegPath)
	defer os.Setenv("FFMPEG_PATH", origFfmpegPath)

	// 使用ffmpeg-go获取视频信息
	probeOutput, err := ffmpeg_go.Probe(inputPath, ffmpeg_go.KwArgs{"show_entries": "format=duration"})
	if err != nil {
		return 0, err
	}

	// 解析输出获取时长
	lines := strings.Split(probeOutput, "\n")
	var duration string
	for _, line := range lines {
		if strings.HasPrefix(line, "format_duration=") {
			duration = strings.TrimPrefix(line, "format_duration=")
			break
		}
	}

	if duration == "" {
		return 0, fmt.Errorf("failed to get duration")
	}

	// 转换时长字符串为秒数（浮点数）
	durationSeconds, err := strconv.ParseFloat(duration, 64)
	if err != nil {
		return 0, err
	}

	// 转换为毫秒并返回
	return int64(durationSeconds * 1000), nil
}
