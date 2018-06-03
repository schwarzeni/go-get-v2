package downloader

import (
	"fmt"
	"net/http"

	"io"
	"os"
	"path"

	"time"

	"path/filepath"

	"github.com/schwarzeni/go-get-v2/core/model"
	parserModel "github.com/schwarzeni/go-get-v2/parser/model"
	"github.com/schwarzeni/go-get-v2/util"
)

func Downloadfunc(sendRequestSignal chan model.GetRequestFromPool, chanToDownload chan parserModel.Video, maxWorkerNumber int, finishSingle chan int, watchSpeedTaskFinishSignal model.SpeedWatchTaskFinishSignal) {
	finishOneWork := make(chan int)
	workerCount := 0
	isFinish := false

	start := time.Now()
	// TODO just test
	jobcount := 0
	var worker parserModel.Video = nil
	for {
		// TODO 添加阻塞模块，防止程序莫名其妙挂掉 否则运行到2min30s作用就不会进行下载了
		if workerCount == maxWorkerNumber || isFinish == true {
			<-finishOneWork
			// TODO delete later
			fmt.Println(".... <-finishOneWork ....")
		}

		// 请求池中还有请求等待同时goroutine的个数未达到上限
		if workerCount < maxWorkerNumber && isFinish == false {
			sendRequestSignal <- 1
			select {
			case <-finishSingle:
				// TODO log here
				util.LogP(fmt.Sprintf("video pool is finish"))
				isFinish = true
			case worker = <-chanToDownload:
				workerCount++
				// TODO log here
				util.LogP(fmt.Sprintf("get new download task"))
				resp, err := worker.Download()
				if err != nil {
					util.LogE("!!!error!!! " + err.Error())
				}
				go SaveToFile(resp, worker.GetSavePath(), func(savePath string) {
					workerCount--
					finishOneWork <- 1
					jobcount++
				})
			}
			// 请求池中已无请求同时已经没有goroutine在工作
		} else if workerCount == 0 && isFinish == true {
			// 结束监控数据速度的任务
			watchSpeedTaskFinishSignal <- 1
			// TODO plause a little time for count to more accurate
			util.SleepAtRandomTime()
			// TODO log here
			util.LogP(fmt.Sprintf("analyse finish %d works, use %s", jobcount, time.Since(start)))
			return
		}
	}
}

// 获取下载内容并将其保存至文件中
func SaveToFile(resp *http.Response, filePath string, finish func(savePath string)) {
	//TODO: save file
	defer resp.Body.Close()

	err := os.MkdirAll(path.Dir(filePath), 0777)
	if err != nil {
		util.LogFatal(err)
	}

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		out, err := os.Create(filePath)
		if err != nil {
			util.LogFatal(err)
		}
		defer out.Close()

		io.Copy(out, resp.Body)

		util.LogP("saved file " + filePath)
	} else {
		util.LogFatal(filePath + " already exist")
	}
	finish(filePath)
}

// 监控网速
func OverwatchNetworkSpeed(pathList []string, taskFinish model.SpeedWatchTaskFinishSignal) {
	go func() {
		var prevFileSize float64 = 0.0
		for {
			select {
			case <-taskFinish:
				return
			default:
				time.Sleep(time.Second / 2)
				currentFileSize := 0.0
				for _, val := range pathList {
					s, _ := DirSize(val)
					currentFileSize = currentFileSize + float64(s)/1000.0
				}
				increaseSize := currentFileSize - prevFileSize
				increaseSize = increaseSize / 2
				if increaseSize > 1000.0 {
					util.LogP(fmt.Sprintf(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>  %f MB/s total: %f MB", increaseSize/1000.0, currentFileSize/1000.0))
				} else {
					util.LogP(fmt.Sprintf(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>  %f KB/s total: %f MB", increaseSize, currentFileSize/1000.0))
				}
				prevFileSize = currentFileSize
			}
		}
	}()
	//<-taskFinish
}

func DirSize(path string) (int64, error) {
	var size int64
	err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if info != nil && !info.IsDir() {
			size += info.Size()
		}
		return err
	})
	return size, err
}
