package engine

import (
	"github.com/schwarzeni/go-get-v2/core/downloader"
	"github.com/schwarzeni/go-get-v2/core/model"
	"github.com/schwarzeni/go-get-v2/core/scheduler"
	parserModel "github.com/schwarzeni/go-get-v2/parser/model"
	"github.com/schwarzeni/go-get-v2/util"
)

func Run(total int, yourParser parserModel.Parser) {

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
