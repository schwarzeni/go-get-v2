package engine

import (
	"bufio"
	"os"

	"strconv"

	"github.com/schwarzeni/go-get-v2/core/downloader"
	"github.com/schwarzeni/go-get-v2/core/model"
	"github.com/schwarzeni/go-get-v2/core/scheduler"
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

func GenerateEngineAndRun() {
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

	// 检测环境配置
	util.CheckEnv()

	sendRequestSignal := make(chan model.GetRequestFromPool)
	chanToDownload := make(chan parserModel.Video)
	finishSignal := make(chan int)

	// 解析获取视频列表与保存路径
	lists, pathsList := yourParser.GetVideoListAndSavePath()

	// 想请求池中添加request
	go scheduler.RequestPool(sendRequestSignal, chanToDownload, lists, finishSignal)

	// 开始下载
	downloader.Downloadfunc(sendRequestSignal, chanToDownload, total, finishSignal)

	// 下载结束后链接文件
	util.ConcatFiles(pathsList)
}
