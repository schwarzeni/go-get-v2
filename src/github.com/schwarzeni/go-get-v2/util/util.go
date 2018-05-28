package util

import (
	"log"
	"math/rand"
	"os"
	"time"

	"os/exec"

	"sync"

	"net/http"

	"io"
	"io/ioutil"

	"bufio"
	"os/user"
	"path"

	"encoding/json"

	"github.com/schwarzeni/go-get-v2/parser/model"
)

// 只是模拟http延迟
func SleepAtRandomTime() {
	r := rand.Intn(2500)
	rateLimit := time.Tick(time.Duration(r) * time.Millisecond)
	<-rateLimit
}

// 生成文件路径列表
func GenerateFilePathList(videos []model.Video) []string {
	m := make(map[string]string)
	var paths []string
	for _, v := range videos {
		if m[v.GetSavePath()] == "" {
			m[v.GetSavePath()] = "--"
			paths = append(paths, v.GetSavePath())
		}
	}
	return paths
}

// 链接文件
func ConcatFiles(pathLists []string) {
	LogP("Begin to concat files ...")
	var wg sync.WaitGroup
	num := len(pathLists)
	wg.Add(num)
	for i := 0; i < num; i++ {
		go func(dirpath string, wg *sync.WaitGroup) {
			LogP("Concating files in dir " + dirpath + " ...")
			cmd := exec.Command("bash", "concat_file.sh", dirpath)
			err := cmd.Run()
			if err != nil {
				LogE("When concating dir: " + dirpath + ", " + err.Error())
			} else {
				LogP("Concat files in dir " + dirpath + " finish")
			}
			wg.Done()
		}(pathLists[i], &wg)
	}
	wg.Wait()
	LogP("Finish")
}

// 检测一些环境配置
func CheckEnv() {
	// 检测环境是否有支持bash
	LogP("Check if bash is available ...")
	cmd := exec.Command("bash", "--version")
	err := cmd.Run()
	if err != nil {
		LogFatal("Fatal error: " + err.Error() + "\n>> Tips: Please check if bash available")
	} else {
		LogP("OK")
	}

	// 检测环境是否有ffmpeg
	LogP("Check if env has ffmpeg ...")
	cmd = exec.Command("ffmpeg", "--help")
	err = cmd.Run()
	if err != nil {
		LogFatal("Fatal error: " + err.Error() + "\n>> Tips: Please check if ffmpeg is available in your command line environment")
	} else {
		LogP("OK")
	}

	// 检测concat_file.sh 文件是否存在
	LogP("Check if has concat_file.sh ...")
	cmd = exec.Command("ls", "concat_file.sh")
	err = cmd.Run()
	if err != nil {
		LogFatal("Fatal error: " + err.Error() + "\n>> Tips: Please check if there is concat_file.sh exists")
	} else {
		LogP("OK")
	}
}

// 发送http请求
func MethodGet(url string, headers map[string]string) (*http.Response, error) {
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// 将http请求返回的body的内容转为字符串
func ResponseBodyToString(body io.ReadCloser) (string, error) {
	bodyBytes, err2 := ioutil.ReadAll(body)
	if err2 != nil {
		return "", err2
	}
	return string(bodyBytes), nil
}

// 处理用户输入的保存路径，并将其转化为绝对路径
func CheckAndParseSaveFilePath(config model.Config) model.Config {
	usr, err := user.Current()
	if err != nil {
		LogFatal(err)
	}
	// ~/Desktop/temp
	// aaa/22
	for idx, _ := range config.Data {
		// 如果是相对路径则进行拼接
		if string([]rune(config.Data[idx].SavePath)[0]) == "~" {
			config.Data[idx].SavePath = path.Join(usr.HomeDir, string([]rune(config.Data[idx].SavePath)[1:]))
		}
	}
	LogP("Is the following path right? Input enter to continue, input other to quit program and edit config file")
	for _, data := range config.Data {
		LogP(data.SavePath)
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
		LogP("Tips: Save path example: ~/Desktop/video1/01 or /Users/nizhenyang/Desktop/video1/01")
		os.Exit(1)
	}
	return config
}

// 解析配置文件
func ParseConfigFile() model.Config {
	file, e := ioutil.ReadFile("./config/data.json")
	if e != nil {
		LogFatal("in util.ParseConfigFile" + e.Error())
	}
	var jsontype model.Config
	json.Unmarshal(file, &jsontype)
	// 检查同时修改路径
	jsontype = CheckAndParseSaveFilePath(jsontype)
	return jsontype
}

func LogP(v interface{}) {
	log.Println(">> ", v)
}

func LogE(v interface{}) {
	log.Println(">> *** something went wrong ***: ", v)
}

func LogFatal(v interface{}) {
	log.Fatal(">>  *** something went badly wrong ***: ", v)
	os.Exit(1)
}
