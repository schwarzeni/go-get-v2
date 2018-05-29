package iqiyi

import (
	"net/http"
	"net/url"
	"strings"

	"strconv"
	"time"

	"encoding/json"

	"github.com/schwarzeni/go-get-v2/util"
)

type IqiyiVideo struct {
	Refer    string
	Origin   string
	Host     string
	Cookie   string
	Url      string
	SavePath string
}

func (i IqiyiVideo) Download() (*http.Response, error) {
	i.SetDownloadUnixTime()
	util.LogP("fetch read download url for: " + i.SavePath)
	i.FetchRealVIdeoUrl()
	u, _ := url.Parse(i.Url)
	util.LogP("send download reqest for: " + u.Host + u.Path)
	u, e := url.Parse(strings.Trim(i.Url, " "))
	if e != nil {
		util.LogFatal("Error in IqiyiVideo.Download: " + e.Error())
	}
	headers := map[string]string{
		"Host":            u.Host,
		"Origin":          i.Origin,
		"Referer":         i.Refer,
		"Cookie":          i.Cookie,
		"Range":           "bytes=0-",
		"Connection":      "keep-alive",
		"Accept-Encoding": "gzip, deflate, br",
		"Accept-Language": "en,zh-CN;q=0.9,zh;q=0.8,zh-TW;q=0.7",
		"User-Agent":      "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/66.0.3359.181 Safari/537.36"}

	return util.MethodGet(i.Url, headers)
}

func (i *IqiyiVideo) SetDownloadUnixTime() {
	u, _ := url.Parse(i.Url)
	q, _ := url.ParseQuery(u.RawQuery)
	q.Add("rn", strconv.Itoa(int(time.Now().Unix())))
	u.RawQuery = q.Encode()
	i.Url = u.String()
}

func (i *IqiyiVideo) FetchRealVIdeoUrl() {
	resp, err := util.MethodGet(i.Url, map[string]string{
		"Connection": "keep-alive",
		"Host":       i.Host,
		"Origin":     i.Origin,
		"Pragma":     "no-cache",
		"Referer":    i.Refer,
		"User-Agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/66.0.3359.181 Safari/537.36"})
	if err != nil {
		util.LogFatal("in IqiyiParser.GetVideoListAndSavePath get read video url " + i.Url + " " + err.Error())
	}
	defer resp.Body.Close()
	content, err := util.ResponseBodyToString(resp.Body)
	if err != nil {
		util.LogFatal("n IqiyiParser.GetVideoListAndSavePath convert response.body to string " + err.Error())
	}
	var jsonBodyRawResult map[string]*json.RawMessage
	json.Unmarshal([]byte(content), &jsonBodyRawResult)
	realurl := string(*jsonBodyRawResult["l"])
	realurl = strings.Trim(realurl, `""`)
	i.Url = realurl
}

func (i IqiyiVideo) GetSavePath() string {
	return i.SavePath
}
