package parser

import (
	"net/http"
	"net/url"
)

type VideoWrapper struct {
	Video Video
}

type Video interface {
	// 下载功能，设置header等操作，获取response
	Download() (http.Response, error)
	// 获取保存的路径
	GetSavePath() string
	// 获取url
	GetUrlString() string
}

type VideoList interface {
	GetNextVideo() Video
}

type Parser interface {
	GetVideoList(url url.URL) []Video
}
