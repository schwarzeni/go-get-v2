package scheduler

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"sync"

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

// 总任务统计
var totalTaskNumber = 0

// 生成视频列表并送到等待池中
func generateVideoList(listNum int, listParser parser.Parser) []parser.Video {

	var wg sync.WaitGroup
	wg.Add(listNum)
	var videos []parser.Video

	for i := 0; i < listNum; i++ {
		go func(i int, wg *sync.WaitGroup, videos *[]parser.Video) {
			u, _ := url.Parse(strconv.Itoa(i))
			list := listParser.GetVideoList(*u)
			for _, v := range list {
				*videos = append(*videos, v)
				totalTaskNumber++
			}
			wg.Done()
		}(i, &wg, &videos)
	}
	wg.Wait()
	// TODO log here
	util.LogP(fmt.Sprintf("download list generate"))
	return videos
}

// 请求等待池
func requestPool(sendRequestSignal chan GetRequestFromPool, chanToDownload chan parser.Video, videos []parser.Video, finishSingle chan int) {
	var videoQueue VideoQueue
	videoQueue.VideoLists = videos
	for {
		select {
		case <-sendRequestSignal:
			var vw parser.Video
			if !videoQueue.IsEmpty() {
				vw = videoQueue.Pop()
				chanToDownload <- vw
			} else {
				finishSingle <- 1
				// TODO log here
				util.LogP(fmt.Sprintf("this pool is finish"))
				return
			}
		}

	}
}

// 下载器
func downloader(sendRequestSignal chan GetRequestFromPool, chanToDownload chan parser.Video, maxWorkerNumber int, finishSingle chan int) {
	workerCount := 0
	isFinish := false
	// TODO just test
	jobcount := 0
	var worker parser.Video = nil
	for {

		if workerCount < maxWorkerNumber {
			if workerCount == 0 && isFinish == true {
				// TODO plause a little time for count to more accurate
				util.SleepAtRandomTime()
				// TODO log here
				util.LogP(fmt.Sprintf("analyse %d", jobcount))
				return
			}
			if isFinish != true {
				sendRequestSignal <- 1
				select {
				case <-finishSingle:
					// TODO log here
					util.LogP(fmt.Sprintf("video pool is finish"))
					isFinish = true
				case worker = <-chanToDownload:
				}
			}
		}
		if worker != nil {
			workerCount++
			// TODO log here
			util.LogP(fmt.Sprintf("%d, get new download task  ", workerCount))
			resp, err := worker.Download()
			if err != nil {
				util.LogE("!!!error!!! " + err.Error())
			}
			go saveToFile(resp, worker.GetSavePath(), func(savePath string) {
				util.LogP(fmt.Sprintf("%d", workerCount) + " save finish " + savePath)
				workerCount--
				jobcount++
			})
			worker = nil
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
	sendRequestSignal := make(chan GetRequestFromPool)
	chanToDownload := make(chan parser.Video)
	finishSignal := make(chan int)

	//TODO: just a test
	listnumber := 3

	lists := generateVideoList(listnumber, yourParser)
	go requestPool(sendRequestSignal, chanToDownload, lists, finishSignal)
	downloader(sendRequestSignal, chanToDownload, total, finishSignal)
}
