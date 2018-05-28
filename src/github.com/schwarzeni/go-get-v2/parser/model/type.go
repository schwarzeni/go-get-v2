package model

import (
	"net/http"
)

type BasicInfo struct {
	Url    string
	Cookie string
	Refer  string
	Host   string
	Origin string
}

type Video interface {
	// 下载功能，设置header等操作，获取response
	Download() (*http.Response, error)
	GetSavePath() string
}

type Parser interface {
	// savepath用于最后链接文件
	GetVideoListAndSavePath() ([]Video, []string)
	// 生成一个parser实例，用于对parser信息的初始化
	BuildParser() Parser
}

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
