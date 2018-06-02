package youku

import (
	"io/ioutil"

	"sync"

	"net/url"

	"strings"

	"bufio"

	"path"
	"strconv"

	"github.com/schwarzeni/go-get-v2/parser/model"
	"github.com/schwarzeni/go-get-v2/util"
)

type YouKuParser struct{}

func (y YouKuParser) BuildParser() model.Parser {
	return YouKuParser{}
}

func (y YouKuParser) GetVideoListAndSavePathForChrome(videoInfo model.SingleVideoInJson) ([]model.Video, string) {
	var videos []model.Video
	apiu, _ := url.Parse(videoInfo.ApiUrl)
	pageu, _ := url.Parse(videoInfo.WebpageUrl)
	strs, e := y.ParsePlayFrameTxt(videoInfo.ApiUrl,
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
	vs := y.GenerateVideoModelsForChrome(strs, videoInfo)
	videos = append(videos, vs...)
	return videos, videoInfo.SavePath
}

func (y YouKuParser) GenerateVideoModelsForChrome(urls []string, videoInfo model.SingleVideoInJson) []model.Video {
	var videos []model.Video
	for i, vu := range urls {
		u, e := url.Parse(vu)
		if e != nil {
			util.LogFatal("Error: in YouKuParser.GenerateVideoModelsForChrome: " + e.Error())
		}
		videos = append(videos, YoukuVideo{
			Url:      u.String(),
			Origin:   "http://v.youku.com",
			Host:     u.Host,
			Refer:    videoInfo.WebpageUrl,
			Cookie:   videoInfo.Cookie,
			SavePath: path.Join(videoInfo.SavePath, strconv.Itoa(i)+".ts")})
	}
	return videos
}

////////////////////// 以下为原方法 ///////////////////////////

func (y YouKuParser) GetVideoListAndSavePath() ([]model.Video, []string) {
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
			strs, e := y.ParsePlayFrameTxt(info.Data[idx].ApiUrl,
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
			vs := y.GenerateVideoModels(strs, &info, idx)
			*videos = append(*videos, vs...)
			*paths = append(*paths, info.Data[idx].SavePath)
			wg.Done()
		}(info, &videos, &pathLists, &wg, i)
	}
	wg.Wait()
	return videos, pathLists
}

// 解析api返回的文件并获得分段列表
func (y YouKuParser) ParsePlayFrameTxt(url string, header map[string]string) ([]string, error) {
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

func (y YouKuParser) GenerateVideoModels(urls []string, info *model.Config, idx int) []model.Video {
	var videos []model.Video
	for i, vu := range urls {
		u, e := url.Parse(vu)
		if e != nil {
			util.LogFatal("Error: in YouKuParser.GenerateVideoModels: " + e.Error())
		}
		videos = append(videos, YoukuVideo{
			Url:      u.String(),
			Origin:   "http://v.youku.com",
			Host:     u.Host,
			Refer:    info.Data[idx].WebpageUrl,
			Cookie:   info.Cookie,
			SavePath: path.Join(info.Data[idx].SavePath, strconv.Itoa(i)+".ts")})
	}
	return videos
}
