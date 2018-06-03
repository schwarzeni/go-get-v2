package engine

import (
	"bufio"
	"os"

	"strconv"

	"fmt"

	"github.com/schwarzeni/go-get-v2/core/downloader"
	"github.com/schwarzeni/go-get-v2/core/model"
	"github.com/schwarzeni/go-get-v2/core/scheduler"
	"github.com/schwarzeni/go-get-v2/parser/_dispatcher"
	"github.com/schwarzeni/go-get-v2/parser/bilibili"
	"github.com/schwarzeni/go-get-v2/parser/iqiyi"
	parserModel "github.com/schwarzeni/go-get-v2/parser/model"
	"github.com/schwarzeni/go-get-v2/parser/tencent"
	"github.com/schwarzeni/go-get-v2/parser/youku"
	"github.com/schwarzeni/go-get-v2/util"
)

var defaultMaxWorkerNum = 20

type EngineInfo struct {
	Parser       parserModel.Parser
	maxWorkerNum int
}

type UserConfigInfo struct {
	MaxWorerkNum      int
	ConfigFilePath    string
	IsNeedToCheckPath bool
}

var defaultWorkerNum = 4

func GenerateEngineAndRun() {
	// 检测环境配置
	util.CheckEnv()

	info := generateEngineInfo()
	run(info.maxWorkerNum, info.Parser)
}

// 根据用户的输入生成解析器
func generateEngineInfo() EngineInfo {
	engineInfo := EngineInfo{}

	// 获取parser类型
	util.LogP("" +
		"\nChoose a parser with it's ID\n" +
		"[ID: 0]Bilibili\n" +
		"[ID: 1]Youku\n" +
		"[ID: 2]Tencent\n" +
		"[ID: 3]Iqiyi")
	scanner := bufio.NewScanner(os.Stdin)
	var id string
	var parser parserModel.Parser
	for scanner.Scan() {
		id = scanner.Text()
		break
	}
	switch id {
	case "0":
		parser = bilibili.BilibiliParser{}.BuildParser()
		util.LogP("Generate Parser for Bilibili")
	case "1":
		parser = youku.YouKuParser{}.BuildParser()
		util.LogP("Generate Parser for Youku")
	case "2":
		parser = tencent.TencentParser{}.BuildParser()
		util.LogP("Generate Parser for Tencent")
	case "3":
		parser = iqiyi.IqiyiParser{}.BuildParser()
		util.LogP("Generate Parser for Iqiyi")
	default:
		util.LogP("Please choose the right ID")
		os.Exit(1)
	}
	engineInfo.Parser = parser

	// 获取最大并发数
	var num string
	util.LogP("\nSet max worker num, enter return then choose default 20")
	for scanner.Scan() {
		num = scanner.Text()
		break
	}
	if num == "" {
		engineInfo.maxWorkerNum = defaultMaxWorkerNum
	} else {
		r, e := strconv.Atoi(num)
		if e != nil {
			util.LogP("Please input a right number")
			os.Exit(1)
		}
		engineInfo.maxWorkerNum = r
	}
	return engineInfo
}

// 启动Engine
func run(total int, yourParser parserModel.Parser) {

	sendRequestSignal := make(chan model.GetRequestFromPool)
	chanToDownload := make(chan parserModel.Video)
	finishSignal := make(chan int)

	// 解析获取视频列表与保存路径
	lists, pathsList := yourParser.GetVideoListAndSavePath()

	// 想请求池中添加request
	go scheduler.RequestPool(sendRequestSignal, chanToDownload, lists, finishSignal)

	// 开始下载
	downloader.Downloadfunc(sendRequestSignal, chanToDownload, total, finishSignal, nil)

	// 下载结束后链接文件
	util.ConcatFiles(pathsList)
}

func RunForChrome() {
	config := ReadUserConfig()
	util.CheckEnv()
	sendRequestSignal := make(chan model.GetRequestFromPool)
	chanToDownload := make(chan parserModel.Video)
	overwatchSpeedFinishSignal := make(model.SpeedWatchTaskFinishSignal)
	finishSignal := make(chan int)

	lists, pathsList := _dispatcher.GetVideoListAndSavePath(config.ConfigFilePath, config.IsNeedToCheckPath)

	// 想请求池中添加request
	go scheduler.RequestPool(sendRequestSignal, chanToDownload, lists, finishSignal)
	go downloader.OverwatchNetworkSpeed(pathsList, overwatchSpeedFinishSignal)
	// 开始下载
	downloader.Downloadfunc(sendRequestSignal, chanToDownload, config.MaxWorerkNum, finishSignal, overwatchSpeedFinishSignal)

	// 下载结束后链接文件
	util.ConcatFiles(pathsList)

	// 清除下载目录
	util.ClearWorkingDir(pathsList)
}

func ReadUserConfig() UserConfigInfo {
	config := UserConfigInfo{MaxWorerkNum: defaultMaxWorkerNum, ConfigFilePath: "", IsNeedToCheckPath: true}
	num := len(os.Args)
	for i := 1; i < num; i++ {
		if os.Args[i] == "-w" && i+1 < num {
			workerNum, err := strconv.Atoi(os.Args[i+1])
			if err != nil {
				util.LogFatal("in ReadUserConfig: " + err.Error())
			}
			config.MaxWorerkNum = workerNum
		}
		if os.Args[i] == "-p" && i+1 < num {
			config.ConfigFilePath = os.Args[i+1]
		}
		if os.Args[i] == "-y" {
			config.IsNeedToCheckPath = false
		}
		if os.Args[i] == "-h" {
			fmt.Println("-w [workernum] -p [config file path] -y [dont check the save file path]")
			os.Exit(0)
		}
	}
	return config
}
