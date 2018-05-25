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
	GetVideoListAndSavePath() ([]Video, []string)
}
