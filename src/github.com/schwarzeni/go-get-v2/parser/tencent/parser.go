package tencent

import (
	"bufio"
	"io/ioutil"
	"net/url"
	"path"
	"strconv"
	"strings"
	"sync"

	"github.com/schwarzeni/go-get-v2/parser/model"
	"github.com/schwarzeni/go-get-v2/util"
)

type TencentParser struct{}

func (TencentParser) BuildParser() model.Parser {
	return TencentParser{}
}

func (t TencentParser) GetVideoListAndSavePathForChrome(videoInfo model.SingleVideoInJson) ([]model.Video, string) {
	var videos []model.Video

	apiu, _ := url.Parse(videoInfo.ApiUrl)
	pageu, _ := url.Parse(videoInfo.WebpageUrl)
	strs, e := t.ParsePlayFrameTxt(videoInfo.ApiUrl,
		map[string]string{
			"Cookie":     videoInfo.Cookie,
			"Connection": "keep-alive",
			"Host":       apiu.Host,
			"Origin":     pageu.Scheme + "://" + pageu.Host,
			"Referer":    pageu.String(),
			"User-Agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/66.0.3359.181 Safari/537.36"})
	if e != nil {
		util.LogFatal("some error in parse api list: " + e.Error())
	}
	vs := t.GenerateVideoModelsForChrome(strs, videoInfo)
	videos = append(videos, vs...)

	return videos, videoInfo.SavePath
}

func (t TencentParser) GenerateVideoModelsForChrome(urls []string, videoInfo model.SingleVideoInJson) []model.Video {
	var videos []model.Video

	apiStrUrl, e := url.Parse(videoInfo.ApiUrl)
	if e != nil {
		util.LogFatal("Error: in TencentParser.GenerateVideoModelsForChrome parse apiStrUrl : " + e.Error())
	}
	p := apiStrUrl.Path
	result := strings.Split(p, "/")
	realPath := ""
	for i := 0; i < len(result)-1; i++ {
		realPath = realPath + result[i] + "/"
	}
	apiStrUrl.RawQuery = ""
	apiStrUrl.Path = realPath

	for i, vu := range urls {
		u := apiStrUrl.String() + vu
		videos = append(videos, TencentVideo{
			Url:      u,
			Origin:   "https://v.qq.com",
			Host:     apiStrUrl.Host,
			Refer:    videoInfo.WebpageUrl,
			Cookie:   videoInfo.Cookie,
			SavePath: path.Join(videoInfo.SavePath, strconv.Itoa(i)+".ts")})
	}
	return videos
}

////////////////////// 以下为原方法 ///////////////////////////

func (t TencentParser) GetVideoListAndSavePath() ([]model.Video, []string) {
	info := util.ParseConfigFile()
	num := len(info.Data)
	var videos []model.Video
	var pathLists []string
	var wg sync.WaitGroup
	wg.Add(num)
	// TODO: ************ do here *************
	for i := 0; i < num; i++ {
		go func(info model.Config, videos *[]model.Video, paths *[]string, wg *sync.WaitGroup, idx int) {
			apiu, _ := url.Parse(info.Data[idx].ApiUrl)
			pageu, _ := url.Parse(info.Data[idx].WebpageUrl)
			strs, e := t.ParsePlayFrameTxt(info.Data[idx].ApiUrl,
				map[string]string{
					"Cookie":     info.Cookie,
					"Connection": "keep-alive",
					"Host":       apiu.Host,
					"Origin":     pageu.Scheme + "://" + pageu.Host,
					"Referer":    pageu.String(),
					"User-Agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/66.0.3359.181 Safari/537.36"})
			if e != nil {
				util.LogFatal("some error in parse api list: " + e.Error())
			}
			vs := t.GenerateVideoModels(strs, &info, idx)
			*videos = append(*videos, vs...)
			*paths = append(*paths, info.Data[idx].SavePath)
			wg.Done()
		}(info, &videos, &pathLists, &wg, i)
	}
	wg.Wait()

	return videos, pathLists
}

// 解析api返回的文件并获得分段列表
func (t TencentParser) ParsePlayFrameTxt(url string, header map[string]string) ([]string, error) {
	var strs []string

	resp, err := util.MethodGet(url, header)
	if err != nil {
		return []string{}, err
	}
	defer resp.Body.Close()
	bodyBytes, err2 := ioutil.ReadAll(resp.Body)
	if err2 != nil {
		return []string{}, err2
	}
	bodyString := string(bodyBytes)
	scanner := bufio.NewScanner(strings.NewReader(bodyString))
	for scanner.Scan() {
		str := scanner.Text()
		if string([]rune(str)[0]) != "#" {
			strs = append(strs, str)
		}
	}
	return strs, nil
}

func (t TencentParser) GenerateVideoModels(urls []string, info *model.Config, idx int) []model.Video {
	var videos []model.Video

	apiStrUrl, e := url.Parse(info.Data[idx].ApiUrl)
	if e != nil {
		util.LogFatal("Error: in TencentParser.GenerateVideoModels parse apiStrUrl : " + e.Error())
	}
	p := apiStrUrl.Path
	result := strings.Split(p, "/")
	realPath := ""
	for i := 0; i < len(result)-1; i++ {
		realPath = realPath + result[i] + "/"
	}
	apiStrUrl.RawQuery = ""
	apiStrUrl.Path = realPath

	for i, vu := range urls {
		u := apiStrUrl.String() + vu
		videos = append(videos, TencentVideo{
			Url:      u,
			Origin:   "https://v.qq.com",
			Host:     apiStrUrl.Host,
			Refer:    info.Data[idx].WebpageUrl,
			Cookie:   info.Cookie,
			SavePath: path.Join(info.Data[idx].SavePath, strconv.Itoa(i)+".ts")})
	}
	return videos
}
