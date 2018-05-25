package scheduler

import (
	"fmt"

	"github.com/schwarzeni/go-get-v2/core/model"
	parserModel "github.com/schwarzeni/go-get-v2/parser/model"
	"github.com/schwarzeni/go-get-v2/util"
)

// 请求等待池
func RequestPool(sendRequestSignal chan model.GetRequestFromPool, chanToDownload chan parserModel.Video, videos []parserModel.Video, finishSingle chan int) {
	var videoQueue VideoQueue
	videoQueue.VideoLists = videos
	for {
		select {
		// 下载器请求获取一个任务
		case <-sendRequestSignal:
			var vw parserModel.Video
			// 如果池中还有任务的话就输送给它，否则就发出结束信号表明以无新下载任务
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
