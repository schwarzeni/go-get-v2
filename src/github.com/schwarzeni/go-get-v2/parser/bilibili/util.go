package bilibili

import (
	"net/http"
)

func SetRequestForListHeader(video BilibiliVideo, req *http.Request) {
	m := map[string]string{
		"Referer":         video.Refer,
		"Accept-Encoding": "gzip, deflate, br",
		"Connection":      "keep-alive",
		"User-Agent":      "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/65.0.3325.181 Safari/537.36",
		"Accept-Language": "en,zh-CN;q=0.9,zh;q=0.8,zh-TW;q=0.7",
		"Origin":          video.Origin,
		"Host":            video.Host,
		"Range":           "bytes=0-",
		"Cookie":          video.Cookie}

	for k, v := range m {
		req.Header.Set(k, v)
	}
}

func MethodGet(url string, headers map[string]string) (*http.Response, error) {
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	return res, nil
}
