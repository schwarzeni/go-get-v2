package bilibili

import (
	"net/http"

	"net/url"

	"strconv"

	"github.com/schwarzeni/go-get-v2/parser"
	"github.com/schwarzeni/go-get-v2/util"
)

type BilibiliVideo struct {
	url      url.URL
	savePath string
}

func (b BilibiliVideo) GetUrlString() string {
	return b.url.String()
}

func (b BilibiliVideo) Download() (http.Response, error) {
	// TODO: just a test
	util.LogP("send download reqest for: " + b.url.String())
	return http.Response{}, nil
}

func (b BilibiliVideo) GetSavePath() string {
	return b.savePath
}

type BilibiliParser struct {
}

// 获取视频列表
func (b BilibiliParser) GetVideoList(apiUrl url.URL) []parser.Video {
	// TODO: just a test
	util.SleepAtRandomTime()
	var videoList []parser.Video
	for i := 0; i < 4; i++ {
		tmpUrl, _ := url.Parse(apiUrl.String() + strconv.Itoa(i))
		videoList = append(videoList, BilibiliVideo{url: *tmpUrl, savePath: tmpUrl.String()})
	}
	return videoList
}
