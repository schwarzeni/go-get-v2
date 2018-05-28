package iqiyi

import (
	"net/url"
	"strconv"
	"time"
)

// 资源链接请求
type IqiyiVideoUrlQuest struct {
	// 链接的url
	Url     *url.URL
	Host    string
	Origin  string
	Pragma  string
	Referer string
	// 请求获得此链接的url
	FromUrl *url.URL
	// 保存路径
	SavePath string
}

func (i IqiyiVideoUrlQuest) GenerateHttpRequestHeader() map[string]string {
	return map[string]string{
		"Connection": "keep-alive",
		"Host":       i.Host,
		"Origin":     i.Origin,
		"Pragma":     "no-cache",
		"Referer":    i.Referer,
		"User-Agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/66.0.3359.181 Safari/537.36"}
}

func (i *IqiyiVideoUrlQuest) SelfConstruct() {

	i.Host = i.Url.Host
	i.Origin = "http://www.iqiyi.com"
	i.Pragma = "no-cache"

	q, _ := url.ParseQuery(i.Url.RawQuery)
	q.Add("cross-domain", "1")
	q.Add("qyid", i.FromUrl.Query().Get("k_uid"))
	q.Add("qypid", i.FromUrl.Query().Get("tvid")+"_"+i.FromUrl.Query().Get("src"))
	q.Add("qypid", i.FromUrl.Query().Get("tvid")+"_"+i.FromUrl.Query().Get("src"))
	q.Add("rn", strconv.Itoa(int(time.Now().Unix())))
	q.Add("pv", "0.1")
	i.Url.RawQuery = q.Encode()
}
