package scheduler

import (
	"fmt"

	"net/url"
	"strconv"

	"net/http"

	"github.com/schwarzeni/go-get-v2/parser"
	"github.com/schwarzeni/go-get-v2/util"
)

type VideoQueue struct {
	VideoLists []parser.Video
}

type FetchVideoFinish int

type GetRequestFromPool int

func (q VideoQueue) IsEmpty() bool {
	return len(q.VideoLists) == 0
}

func (q *VideoQueue) Pop() parser.Video {
	if q.IsEmpty() {
		return nil
	}
	v := q.VideoLists[0]
	q.VideoLists = q.VideoLists[1:]
	return v
}

func (q *VideoQueue) Push(v parser.Video) {
	q.VideoLists = append(q.VideoLists, v)
}

func download(url chan string) {
	for v := range url {
		fmt.Println(v)
	}
}

func fetchStrs(url chan string, num int) {
	for i := 0; i < num; i++ {
		go func(i int) {
			util.SleepAtRandomTime()
			url <- fmt.Sprintf("[%d] generate", i)
		}(i)
	}
}

// 生成视频列表并送到等待池中
func generateVideoList(listNum int, listParser parser.Parser, newVideo chan parser.VideoWrapper) {
	for i := 0; i < listNum; i++ {
		go func(i int) {
			u, _ := url.Parse(strconv.Itoa(i))
			list := listParser.GetVideoList(*u)
			for _, v := range list {
				vw := parser.VideoWrapper{v}
				newVideo <- vw
			}
		}(i)
	}
}

// 请求等待池
func requestPool(sendRequestSignal chan GetRequestFromPool, chanToDownload chan parser.VideoWrapper, newVideo chan parser.VideoWrapper, finish chan FetchVideoFinish) {
	var videoQueue VideoQueue
	for {
		select {
		case <-sendRequestSignal:
			var vw parser.VideoWrapper
			if !videoQueue.IsEmpty() {
				vw = parser.VideoWrapper{videoQueue.Pop()}
			} else {
				vw = parser.VideoWrapper{Video: nil}
			}
			chanToDownload <- vw
		case tmpVideo := <-newVideo:
			// TODO log here
			util.LogP("get new video " + tmpVideo.Video.GetUrlString())
			videoQueue.Push(tmpVideo.Video)
		case <-finish:
			return
		}

	}
}

// 下载器
func downloader(sendRequestSignal chan GetRequestFromPool, chanToDownload chan parser.VideoWrapper, maxWorkerNumber int) {
	count := 0
	for {

		if count <= maxWorkerNumber {
			// TODO log here
			util.LogP("ask for new download task")
			sendRequestSignal <- 1
		}
		worker := <-chanToDownload
		// TODO log here
		util.LogP(fmt.Sprintf("count: %d\n", count))
		if worker.Video != nil {
			// TODO log here
			util.LogP("get new download task " + worker.Video.GetUrlString())
			count++
			resp, err := worker.Video.Download()
			if err != nil {
				util.LogE("!!!error!!! " + err.Error())
			}
			go saveToFile(resp, worker.Video.GetSavePath(), func(savePath string) {
				util.LogP(fmt.Sprintf("%d", count) + " save finish " + savePath)
				count--
			})
		}
	}
}

// 获取下载内容并将其保存至文件中
func saveToFile(resp http.Response, savePath string, finish func(savePath string)) {
	//TODO: save file
	util.SleepAtRandomTime()
	finish(savePath)
}

func Engine(total int, yourParser parser.Parser) {
	video := make(chan parser.VideoWrapper)
	sendRequestSignal := make(chan GetRequestFromPool)
	chanToDownload := make(chan parser.VideoWrapper)
	finish := make(chan FetchVideoFinish)

	//TODO: just a test
	listnumber := 10

	go requestPool(sendRequestSignal, chanToDownload, video, finish)
	generateVideoList(listnumber, yourParser, video)
	downloader(sendRequestSignal, chanToDownload, total)
}
