package bilibili

import (
	"encoding/json"
	"io/ioutil"

	"sync"

	"net/url"

	"path"
	"strconv"

	"bufio"
	"os"

	"strings"

	"github.com/schwarzeni/go-get-v2/parser/model"
	"github.com/schwarzeni/go-get-v2/util"
)

type BilibiliParser struct {
	IsUsePlugin bool
}

func (BilibiliParser) BuildParser() model.Parser {
	parser := BilibiliParser{}
	// 如果是第三方插件的话需要修改一些东西
	util.LogP("\nIs the api from third party plugin? If yes, input y, alse others")
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		input := scanner.Text()
		if input == "y" {
			parser.IsUsePlugin = true
		}
		break
	}
	return parser
}

func (b BilibiliParser) GetVideoListAndSavePathForChrome(videoInfo model.SingleVideoInJson) ([]model.Video, string) {
	var videos []model.Video

	v, e := ParseJsonApi(videoInfo.ApiUrl,
		map[string]string{
			"Host":    "bangumi.bilibili.com",
			"Origin":  "https://www.bilibili.com",
			"Referer": videoInfo.WebpageUrl,
			"Cookie":  videoInfo.Cookie})
	if e != nil {
		util.LogFatal("in GetVideoList " + e.Error())
	}
	vs := b.ConvertJsonToVideoModelsForChrome(v, videoInfo)
	videos = append(videos, vs...)

	return videos, videoInfo.SavePath
}

// 将json格式的视频信息转换为model
func (b BilibiliParser) ConvertJsonToVideoModelsForChrome(vj VideoListJson, vinfo model.SingleVideoInJson) []model.Video {
	var bs []model.Video
	for i := 0; i < len(vj.Durl); i++ {
		u, e := url.Parse(vj.Durl[i].Url)
		if e != nil {
			util.LogFatal("Error: in ConvertJsonToVideoModelsForChrome: " + e.Error())
		}
		u.Scheme = "https"
		if b.IsUsePlugin {
			u.RawQuery = strings.Replace(u.RawQuery, "platform=iphone", "platform=pc", -1)
		}
		bs = append(bs, BilibiliVideo{
			Url:      u.String(),
			Origin:   "https://www.bilibili.com",
			Host:     u.Host,
			Refer:    vinfo.WebpageUrl,
			Cookie:   vinfo.Cookie,
			SavePath: path.Join(vinfo.SavePath, strconv.Itoa(i)+".flv")})
	}
	return bs
}

////////////////////// 以下为原方法 ///////////////////////////

// 获取视频列表
func (b BilibiliParser) GetVideoListAndSavePath() ([]model.Video, []string) {
	info := util.ParseConfigFile()
	num := len(info.Data)
	var videos []model.Video
	var pathLists []string
	var wg sync.WaitGroup
	wg.Add(num)
	for i := 0; i < num; i++ {
		go func(info model.Config, videos *[]model.Video, paths *[]string, wg *sync.WaitGroup, idx int) {
			v, e := ParseJsonApi(info.Data[idx].ApiUrl,
				map[string]string{
					"Host":    "bangumi.bilibili.com",
					"Origin":  "https://www.bilibili.com",
					"Referer": info.Data[idx].WebpageUrl,
					"Cookie":  info.Cookie})
			if e != nil {
				util.LogFatal("in GetVideoList " + e.Error())
			}
			vs := b.ConvertJsonToVideoModels(v, &info, idx)
			*videos = append(*videos, vs...)
			*paths = append(*paths, info.Data[idx].SavePath)
			wg.Done()
		}(info, &videos, &pathLists, &wg, i)
	}
	wg.Wait()

	return videos, pathLists
}

// 获取一个视频的分段列表并解析，返回
func ParseJsonApi(url string, header map[string]string) (VideoListJson, error) {
	resp, err := util.MethodGet(url, header)
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
func (b BilibiliParser) ConvertJsonToVideoModels(vj VideoListJson, info *model.Config, idx int) []model.Video {
	var bs []model.Video
	for i := 0; i < len(vj.Durl); i++ {
		u, e := url.Parse(vj.Durl[i].Url)
		if e != nil {
			util.LogFatal("Error: in ConvertJsonToVideoModels: " + e.Error())
		}
		u.Scheme = "https"
		if b.IsUsePlugin {
			u.RawQuery = strings.Replace(u.RawQuery, "platform=iphone", "platform=pc", -1)
		}
		bs = append(bs, BilibiliVideo{
			Url:      u.String(),
			Origin:   "https://www.bilibili.com",
			Host:     u.Host,
			Refer:    info.Data[idx].WebpageUrl,
			Cookie:   info.Cookie,
			SavePath: path.Join(info.Data[idx].SavePath, strconv.Itoa(i)+".flv")})
	}
	return bs
}
