package youku

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/schwarzeni/go-get-v2/util"
)

type YoukuVideo struct {
	Refer    string
	Origin   string
	Host     string
	Cookie   string
	Url      string
	SavePath string
}

func (y YoukuVideo) Download() (*http.Response, error) {
	// TODO: just a test
	u, _ := url.Parse(y.Url)
	util.LogP("send download reqest for: " + u.Host + u.Path)

	u, e := url.Parse(strings.Trim(y.Url, " "))
	if e != nil {
		util.LogFatal("Error in YoukuVideo.Download: " + e.Error())
	}
	headers := map[string]string{
		"Host":            u.Host,
		"Origin":          y.Origin,
		"Referer":         y.Refer,
		"Cookie":          y.Cookie,
		"Range":           "bytes=0-",
		"Connection":      "keep-alive",
		"Accept-Encoding": "gzip, deflate, br",
		"Accept-Language": "en,zh-CN;q=0.9,zh;q=0.8,zh-TW;q=0.7",
		"User-Agent":      "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/66.0.3359.181 Safari/537.36"}

	return util.MethodGet(y.Url, headers)
}

func (y YoukuVideo) GetSavePath() string {
	return y.SavePath
}
