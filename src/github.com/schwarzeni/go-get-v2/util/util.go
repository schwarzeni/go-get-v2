package util

import (
	"log"
	"math/rand"
	"os"
	"time"

	"os/exec"

	"sync"

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
