package ffmpeg

// ProgressCallback 定义进度回调函数类型
// 用于在FFmpeg处理过程中实时获取进度信息
// 参数:
//   progress: 进度信息结构体指针，包含当前进度百分比、处理时间等
// 返回值:
//   无
// 示例:
//   callback := func(progress *Progress) {
//       fmt.Printf("Progress: %.2f%%\n", progress.Percentage)
//   }
type ProgressCallback func(progress *Progress)

// Progress 定义进度信息结构体
// 包含FFmpeg处理过程中的进度数据
// 字段:
//   Percentage: 进度百分比，范围0-100
//   Current: 当前处理时间，单位为毫秒
//   Total: 总时长，单位为毫秒
//   Status: 当前状态，如"processing"、"completed"等
type Progress struct {
	Percentage float64 // 进度百分比 (0-100)
	Current    int64   // 当前处理时间 (毫秒)
	Total      int64   // 总时长 (毫秒)
	Status     string  // 当前状态
}

// FFmpeg 定义FFmpeg工具结构体
// 用于封装FFmpeg命令行工具的调用和管理
// 字段:
//   FFmpegPath: FFmpeg二进制文件路径
//   ExtractPath: FFmpeg二进制文件释放路径
//   Callback: 进度回调函数
// 注意:
//   该结构体不应直接实例化，而应通过NewFFmpeg系列函数创建
type FFmpeg struct {
	FFmpegPath  string           // FFmpeg二进制文件路径
	ExtractPath string           // FFmpeg二进制文件释放路径
	Callback    ProgressCallback // 进度回调函数
}

// ExtractAudioParams 提取音频流参数结构体
// 用于配置从视频文件中提取音频流的参数
// 字段:
//   InputPath: 输入视频文件路径
//   OutputPath: 输出音频文件路径
type ExtractAudioParams struct {
	InputPath  string // 输入视频文件路径
	OutputPath string // 输出音频文件路径
}

// SplitVideoParams 视频分段参数结构体
// 用于配置将视频文件分割为多个小段的参数
// 字段:
//   InputPath: 输入视频文件路径
//   OutputDir: 输出目录，用于存放分段后的视频文件
//   SegmentTime: 分段时长，单位为秒
//   OutputPrefix: 输出文件名前缀
type SplitVideoParams struct {
	InputPath    string // 输入视频文件路径
	OutputDir    string // 输出目录
	SegmentTime  int    // 分段时长 (秒)
	OutputPrefix string // 输出文件名前缀
}

// ExtractKeyFramesParams 提取关键帧参数结构体
// 用于配置从视频文件中提取关键帧的参数
// 字段:
//   InputPath: 输入视频文件路径
//   OutputDir: 输出目录，用于存放提取的关键帧图片
//   FrameInterval: 关键帧间隔，单位为秒
//   OutputPrefix: 输出文件名前缀
type ExtractKeyFramesParams struct {
	InputPath     string // 输入视频文件路径
	OutputDir     string // 输出目录
	FrameInterval int    // 关键帧间隔 (秒)
	OutputPrefix  string // 输出文件名前缀
}
