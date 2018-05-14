package pc

import (
	"log"

	"fmt"
	"io"
	"os"

	"sync"

	"net/url"

	"strconv"

	"path"

	"github.com/schwarzeni/govideodownloader/model"
	"github.com/schwarzeni/govideodownloader/util"
)

type downloadWorker struct {
	video chan model.SingleVedioApi
	idx   int
	done  func()
}

var cookie string

func Engine() {

	tasks := util.ParseConfigFile()
	cookie = tasks.Cookie
	num := len(tasks.Data)

	var wg sync.WaitGroup
	wg.Add(num)
	for i := 0; i < num; i++ {
		api := tasks.Data[i].ApiUrl
		pageUrl := tasks.Data[i].WebpageUrl
		dirpath := tasks.Data[i].SavePath

		go func(api string, pageUrl string, dirpath string, header map[string]string, wg *sync.WaitGroup) {
			v, e := util.GenerateVideoInfoForWeb(api, pageUrl, header)
			if e != nil {
				log.Println(e)
			}
			DownloadForWeb(v, dirpath)
			wg.Done()
		}(api, pageUrl, dirpath, map[string]string{
			"Cookie": tasks.Cookie}, &wg)
	}
	wg.Wait()

}

func downloadWorking(worker downloadWorker, download func(url string, idx int)) {
	for {
		video := <-worker.video
		u, _ := url.Parse(video.Url)
		log.Printf("start [%s]", u.Host+u.Path)
		download(video.Url, worker.idx)
		worker.done()
	}
}

func createDownloader(wg *sync.WaitGroup, i int, download func(url string, idx int)) downloadWorker {
	worker := downloadWorker{make(chan model.SingleVedioApi), i, func() {
		wg.Done()
	}}
	go downloadWorking(worker, download)
	return worker
}

func DownloadForWeb(videoInfo model.Video, dirpath string) {
	videoSectionLen := len(videoInfo.VedioSection)
	var workers []downloadWorker
	var wg sync.WaitGroup
	wg.Add(videoSectionLen)
	for i := 0; i < videoSectionLen; i++ {
		workers = append(workers, createDownloader(&wg, i, func(downloadUrl string, idx int) {

			// TODO：提取函数
			u, _ := url.Parse(downloadUrl)
			host := "https://" + u.Host
			i := strconv.Itoa(idx)
			paths := path.Join(dirpath, i+".flv")

			httpToHttpsStr(u)

			DownloadSingleVedioForWeb(fmt.Sprintf("%s", u), map[string]string{
				// TODO: 修改 refer以及host
				"Referer":         videoInfo.WebPageUrl,
				"Accept-Encoding": "gzip, deflate, br",
				"Connection":      "keep-alive",
				"User-Agent":      "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/65.0.3325.181 Safari/537.36",
				"Accept-Language": "en,zh-CN;q=0.9,zh;q=0.8,zh-TW;q=0.7",
				"Origin":          "https://www.bilibili.com",
				"Host":            host,
				"Range":           "bytes=0-",
				"Cookie":          cookie},
				paths)
		}))
	}
	for i := 0; i < videoSectionLen; i++ {
		workers[i].video <- videoInfo.VedioSection[i]
	}
	wg.Wait()
}

func DownloadSingleVedioForWeb(currentVedioUrl string, header map[string]string, filePath string) {
	resp, err := util.MethodGet(currentVedioUrl, header)
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()

	err = os.MkdirAll(path.Dir(filePath), 0777)
	if err != nil {
		log.Fatal(err)
	}

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		out, err := os.Create(filePath)
		if err != nil {
			log.Fatal(err)
		}
		defer out.Close()

		io.Copy(out, resp.Body)

		fmt.Println("saved file ", filePath)
	} else {
		fmt.Println(filePath, " already exist")
	}
}

func httpToHttpsStr(u *url.URL) {
	u.Scheme = "https"
}
