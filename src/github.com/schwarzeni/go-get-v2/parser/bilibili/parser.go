package bilibili

import (
	"encoding/json"
	"io/ioutil"

	"sync"

	"net/url"

	"path"
	"strconv"

	"github.com/gpmgo/gopm/modules/log"
	"github.com/schwarzeni/go-get-v2/parser/model"
	"github.com/schwarzeni/go-get-v2/util"
)

type BilibiliParser struct {
}

// 获取视频列表
func (b BilibiliParser) GetVideoListAndSavePath() ([]model.Video, []string) {
	info := b.ParseConfigFile()
	num := len(info.Data)
	var videos []model.Video
	var pathLists []string
	var wg sync.WaitGroup
	wg.Add(num)
	for i := 0; i < num; i++ {
		go func(info Config, videos *[]model.Video, paths *[]string, wg *sync.WaitGroup, idx int) {
			v, e := ParseJsonApi(info.Data[idx].ApiUrl,
				map[string]string{
					"Host":    "bangumi.bilibili.com",
					"Origin":  "https://www.bilibili.com",
					"Referer": info.Data[idx].WebpageUrl,
					"Cookie":  info.Cookie})
			if e != nil {
				util.LogFatal("in GetVideoList " + e.Error())
			}
			vs := ConvertJsonToVideoModels(v, &info, idx)
			*videos = append(*videos, vs...)
			*paths = append(*paths, info.Data[idx].SavePath)
			wg.Done()
		}(info, &videos, &pathLists, &wg, i)
	}
	wg.Wait()

	return videos, pathLists
}

// 解析本地配置文件
func (b BilibiliParser) ParseConfigFile() Config {
	file, e := ioutil.ReadFile("./bilibili-data.json")
	if e != nil {
		util.LogFatal("in ParseConfigFile" + e.Error())
	}

	var jsontype Config
	json.Unmarshal(file, &jsontype)
	return jsontype
}

// 获取一个视频的分段列表并解析，返回
func ParseJsonApi(url string, header map[string]string) (VideoListJson, error) {
	resp, err := MethodGet(url, header)
	if err != nil {
		return VideoListJson{}, err
	}
	defer resp.Body.Close()
	jsonBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return VideoListJson{}, err
	}
	var jsonBodyRawResult map[string]*json.RawMessage
	var jsonBodyResult VideoListJson
	json.Unmarshal(jsonBody, &jsonBodyRawResult)
	json.Unmarshal(*jsonBodyRawResult["durl"], &jsonBodyResult.Durl)
	jsonBodyResult.Refer = header["Refer"]
	return jsonBodyResult, nil
}

// 将json格式的视频信息转换为model
func ConvertJsonToVideoModels(vj VideoListJson, info *Config, idx int) []model.Video {
	var b []model.Video
	for i := 0; i < len(vj.Durl); i++ {
		u, e := url.Parse(vj.Durl[i].Url)
		if e != nil {
			log.Fatal("Error: in ConvertJsonToVideoModels: ", e)
		}
		u.Scheme = "https"
		b = append(b, BilibiliVideo{
			Url:      u.String(),
			Origin:   "https://www.bilibili.com",
			Host:     u.Host,
			Refer:    info.Data[idx].WebpageUrl,
			Cookie:   info.Cookie,
			SavePath: path.Join(info.Data[idx].SavePath, strconv.Itoa(i)+".flv")})
	}
	return b
}
