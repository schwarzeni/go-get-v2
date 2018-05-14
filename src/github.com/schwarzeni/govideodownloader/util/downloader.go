package util

import (
	"fmt"
	"io"
	"log"
	"net/url"
	"os"
	"strings"

	"sync"

	"io/ioutil"
	"path"

	"github.com/schwarzeni/govideodownloader/model"
)

type downloadWorker struct {
	video chan model.SingleVedioApi
	done  func()
}

func downloadWorking(worker downloadWorker) {
	for {
		video := <-worker.video
		log.Printf("start [%s]", getfileName(video.Url))
		downloadSingleVideo(video.Url, map[string]string{})
		worker.done()
	}
}

func createDownloader(wg *sync.WaitGroup) downloadWorker {
	worker := downloadWorker{make(chan model.SingleVedioApi), func() {
		wg.Done()
	}}
	go downloadWorking(worker)
	return worker
}

func Downloader(videoInfo model.Video, dirpath string) {
	videoSectionLen := len(videoInfo.VedioSection)
	var workers []downloadWorker
	var wg sync.WaitGroup
	wg.Add(videoSectionLen)
	for i := 0; i < videoSectionLen; i++ {
		workers = append(workers, createDownloader(&wg))
	}
	for i := 0; i < videoSectionLen; i++ {
		workers[i].video <- videoInfo.VedioSection[i]
	}
	wg.Wait()

	r, _ := url.Parse(videoInfo.VedioSection[0].Url)
	generalId := strings.Split(r.Path, "-")[0]
	moveFiles(dirpath, generalId)
}

func downloadSingleVideo(currentVedioUrl string, header map[string]string) {

	resp, err := MethodGet(currentVedioUrl, header)
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()

	var filename = getfileName(currentVedioUrl)

	if _, err := os.Stat(filename); os.IsNotExist(err) {
		out, err := os.Create(filename)
		if err != nil {
			log.Fatal(err)
		}
		defer out.Close()

		io.Copy(out, resp.Body)

		fmt.Println("saved file ", filename)
	} else {
		fmt.Println(filename, " already exist")
	}

}

func getfileName(s string) string {
	u, _ := url.Parse(s)
	filepath := u.Path
	paths := strings.Split(filepath, "/")
	return paths[len(paths)-1]
}

func moveFiles(dirname string, videoID string) {
	os.Mkdir(dirname, 0777)
	files, err := ioutil.ReadDir("./")
	if err != nil {
		log.Println(err)
	}

	for _, f := range files {
		if !f.IsDir() {
			ff := strings.Split(f.Name(), "-")[0]
			if strings.Compare(ff, videoID) == 0 {
				os.Rename(f.Name(), path.Join(dirname, f.Name()))
			}
		}
	}
}
