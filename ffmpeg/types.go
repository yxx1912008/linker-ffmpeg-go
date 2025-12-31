package ffmpeg_tools

// ProgressCallback 定义进度回调函数类型
type ProgressCallback func(progress *Progress)

// Progress 定义进度信息
type Progress struct {
	Percentage float64 // 进度百分比 (0-100)
	Current    int64   // 当前处理时间 (毫秒)
	Total      int64   // 总时长 (毫秒)
	Status     string  // 当前状态
}

// FFmpeg 定义FFmpeg工具结构体
type FFmpeg struct {
	FFmpegPath  string           // FFmpeg二进制文件路径
	ExtractPath string           // FFmpeg二进制文件释放路径
	Callback    ProgressCallback // 进度回调函数
}

// ExtractAudioParams 提取音频流参数
type ExtractAudioParams struct {
	InputPath  string // 输入视频文件路径
	OutputPath string // 输出音频文件路径
}

// SplitVideoParams 视频分段参数
type SplitVideoParams struct {
	InputPath    string // 输入视频文件路径
	OutputDir    string // 输出目录
	SegmentTime  int    // 分段时长 (秒)
	OutputPrefix string // 输出文件名前缀
}

// ExtractKeyFramesParams 提取关键帧参数
type ExtractKeyFramesParams struct {
	InputPath     string // 输入视频文件路径
	OutputDir     string // 输出目录
	FrameInterval int    // 关键帧间隔 (秒)
	OutputPrefix  string // 输出文件名前缀
}
