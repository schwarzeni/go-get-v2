package bilibili

// 配置文件的json格式
type Config struct {
	Data   []VedioInfoConfig `json:"data"`
	Cookie string            `json:"cookie"`
}

// 单个配置视频的列表
type VedioInfoConfig struct {
	ApiUrl     string `json:"apiUrl"`
	WebpageUrl string `json:"webpageUrl"`
	SavePath   string `json:"savePath"`
}

// 请求视频播放列表返回的部分json内容
type VideoListJson struct {
	Durl  []SingleVedioApi
	Refer string
}

// 单个视频分段的部分信息
type SingleVedioApi struct {
	Length int    `json:"length"`
	Order  int    `json:"order"`
	Size   int    `json:"size"`
	Url    string `json:"url"`
}
