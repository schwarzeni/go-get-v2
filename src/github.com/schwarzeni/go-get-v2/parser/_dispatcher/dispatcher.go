package _dispatcher

import (
	"bufio"
	"encoding/json"
	"io/ioutil"
	"os"
	"os/user"
	"path"

	"sync"

	"fmt"

	"github.com/schwarzeni/go-get-v2/parser/model"
	"github.com/schwarzeni/go-get-v2/util"
)

var isNeedToCheckDefault = true

func GetVideoListAndSavePath(jsonfilePath string, isNeedToCheck bool) ([]model.Video, []string) {
	isNeedToCheckDefault = isNeedToCheck
	tasklist := GenerateParsersTasks(jsonfilePath)
	num := len(tasklist)
	var wg sync.WaitGroup
	wg.Add(num)
	//TODO log here
	fmt.Println(num)

	var videos []model.Video
	var paths []string
	// TODO: log
	util.LogP("Begin to generate task list...")
	for idx := 0; idx < num; idx++ {
		go func(wg *sync.WaitGroup, videos *[]model.Video, paths *[]string, idx int) {
			currentTask := tasklist[idx]
			v, p := currentTask.Parser.GetVideoListAndSavePathForChrome(currentTask.VideoInfo)
			*videos = append(*videos, v...)
			*paths = append(*paths, p)
			wg.Done()
		}(&wg, &videos, &paths, idx)
	}

	wg.Wait()
	util.LogP("Generate task list finish...")
	return videos, paths
}

// 生成任务列表
func GenerateParsersTasks(jsonfilePath string) []model.ParserAndVideo {
	config := ParseJsonFile(jsonfilePath)
	var results []model.ParserAndVideo
	for _, val := range config.Videos {
		p := GenerateParser(val.WebId)
		if p != nil {
			results = append(results, model.ParserAndVideo{Parser: p, VideoInfo: val})
		}
	}
	return results
}

// 解析json配置文件
func ParseJsonFile(jsonfilePath string) model.JsonConfigFile {
	file, e := ioutil.ReadFile(jsonfilePath)
	if e != nil {
		util.LogFatal("in util.ParseConfigFile" + e.Error())
	}
	var jsontype model.JsonConfigFile
	json.Unmarshal(file, &jsontype)
	// 检查同时修改路径
	jsontype = CheckAndParseSaveFilePath(jsontype)

	return jsontype
}

// 处理用户输入的保存路径，并将其转化为绝对路径
func CheckAndParseSaveFilePath(config model.JsonConfigFile) model.JsonConfigFile {
	usr, err := user.Current()
	if err != nil {
		util.LogFatal(err)
	}
	// ~/Desktop/temp
	// aaa/22
	for idx, _ := range config.Videos {
		// 如果是相对路径则进行拼接
		if len(config.Videos[idx].SavePath) > 0 && string([]rune(config.Videos[idx].SavePath)[0]) == "~" {
			config.Videos[idx].SavePath = path.Join(usr.HomeDir, string([]rune(config.Videos[idx].SavePath)[1:]))
		}
	}
	if isNeedToCheckDefault == true {
		util.LogP("Is the following path right? Input enter to continue, input other to quit program and edit config file")
		for _, data := range config.Videos {
			util.LogP(data.SavePath)
		}
		scanner := bufio.NewScanner(os.Stdin)
		isRight := false
		for scanner.Scan() {
			if scanner.Text() == "" {
				isRight = true
			}
			break
		}
		if isRight == false {
			util.LogP("Tips: Save path example: ~/Desktop/video1/01 or /Users/nizhenyang/Desktop/video1/01")
			os.Exit(1)
		}
	}
	return config
}

// 根据ID返回parser
func GenerateParser(id string) model.Parser {
	if parser, ok := ParserMap[id]; ok {
		return parser
	}
	return nil
}
