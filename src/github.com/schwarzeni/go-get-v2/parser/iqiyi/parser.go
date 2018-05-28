package iqiyi

import (
	"encoding/json"
	"io/ioutil"

	"sync"

	"net/url"

	"regexp"

	"strings"

	"path"

	"strconv"

	"github.com/schwarzeni/go-get-v2/parser/model"
	"github.com/schwarzeni/go-get-v2/util"
)

var dataCenter = "http://data.video.iqiyi.com/videos"

type IqiyiParser struct {
	IsVip bool
}

func (IqiyiParser) BuildParser() model.Parser {
	return IqiyiParser{IsVip: false}
}

// 获取视频列表
func (i IqiyiParser) GetVideoListAndSavePath() ([]model.Video, []string) {
	config := util.ParseConfigFile()
	var savePaths []string
	var videos []model.Video
	for _, data := range config.Data {
		savePaths = append(savePaths, data.SavePath)
	}
	iqiyiVideoUrlQuests := i.GenerateDownloadQuestUrl(config)

	// 获取视频真实地址
	num := len(iqiyiVideoUrlQuests)
	var wg sync.WaitGroup
	wg.Add(num)
	for idx, quest := range iqiyiVideoUrlQuests {
		go func(idx int, videos *[]model.Video, quest IqiyiVideoUrlQuest, wg *sync.WaitGroup) {
			quest.SelfConstruct()
			resp, err := util.MethodGet(quest.Url.String(), quest.GenerateHttpRequestHeader())
			if err != nil {
				util.LogFatal("in IqiyiParser.GetVideoListAndSavePath get read video url " + quest.Url.String() + " " + err.Error())
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
			*videos = append(*videos, IqiyiVideo{
				Refer:    quest.Referer,
				Origin:   quest.Origin,
				Host:     quest.Host,
				Cookie:   config.Cookie,
				Url:      realurl,
				SavePath: path.Join(quest.SavePath, strconv.Itoa(idx)+".f4v")})
			wg.Done()
		}(idx, &videos, quest, &wg)
	}
	wg.Wait()
	return videos, savePaths
}

// 获取一个包含视频请求下载列表的js文件并解析出链接
func (i IqiyiParser) GenerateDownloadQuestUrl(config model.Config) []IqiyiVideoUrlQuest {
	num := len(config.Data)
	var wg sync.WaitGroup
	wg.Add(num)
	var iqiyiVideoUrlQuests []IqiyiVideoUrlQuest
	for idx := 0; idx < num; idx++ {
		go func(idx int, config model.Config, iqiyiVideoUrlQuests *[]IqiyiVideoUrlQuest, wg *sync.WaitGroup) {
			ul, e := url.Parse(config.Data[idx].ApiUrl)
			if e != nil {
				util.LogFatal("in IqiyiParser.GenerateDownloadQuestUrl parse url " + ul.String() + " " + e.Error())
			}
			resp, err := util.MethodGet(ul.String(), map[string]string{
				"Cookie":        config.Cookie,
				"Host":          ul.Host,
				"Refer":         config.Data[idx].WebpageUrl,
				"User-Agen":     "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/66.0.3359.181 Safari/537.36",
				"Cache-Control": "no-cache",
				"Connection":    "keep-alive",
				"Pragma":        "no-cache"})
			if err != nil {
				util.LogFatal("in IqiyiParser.GenerateDownloadQuestUrl fetch url " + ul.String() + " " + err.Error())
			}
			defer resp.Body.Close()
			bodyBytes, err2 := ioutil.ReadAll(resp.Body)
			if err2 != nil {
				util.LogFatal("in IqiyiParser.GenerateDownloadQuestUrl read response " + ul.String() + " " + err.Error())
			}
			bodyString := string(bodyBytes)
			strs := i.parseJsonFromJsFile(bodyString, ul)
			for _, s := range strs {
				tu, e := url.Parse(s)
				if e != nil {
					util.LogFatal("in IqiyiParser.GenerateDownloadQuestUrl parse url " + s + " " + e.Error())
				}
				*iqiyiVideoUrlQuests = append(*iqiyiVideoUrlQuests, IqiyiVideoUrlQuest{
					Referer:  dataCenter,
					FromUrl:  ul,
					Url:      tu,
					SavePath: config.Data[idx].SavePath})
			}
			wg.Done()
		}(idx, config, &iqiyiVideoUrlQuests, &wg)
	}
	wg.Wait()
	return iqiyiVideoUrlQuests
}

// 那个js文件中json有用的数据格式
type videoJsonFromJsStruct struct {
	Data struct {
		Program struct {
			Video []struct {
				Fs []struct {
					L string `json:"l"`
				} `json:"fs"`
			} `json:"video"`
		} `json:"program"`
	} `json:"data"`
}

// 处理收到的js文件，提取出链接并做处理
func (i IqiyiParser) parseJsonFromJsFile(str string, fromUrl *url.URL) []string {
	str = strings.Trim(str, "")
	var re = regexp.MustCompile(`^[\w]+{[\w\d]+\(`)
	str = re.ReplaceAllString(str, "")
	re = regexp.MustCompile(`\);}catch\(e\){};`)
	str = re.ReplaceAllString(str, "")
	var vs videoJsonFromJsStruct
	json.Unmarshal([]byte(str), &vs)

	result := []string{}
	for _, data := range vs.Data.Program.Video {
		if len(data.Fs) != 0 {
			for _, f := range data.Fs {
				result = append(result, dataCenter+f.L)
			}
		}
	}
	return result
}
