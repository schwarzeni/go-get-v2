package _dispatcher

import (
	"github.com/schwarzeni/go-get-v2/parser/bilibili"
	"github.com/schwarzeni/go-get-v2/parser/iqiyi"
	"github.com/schwarzeni/go-get-v2/parser/model"
	"github.com/schwarzeni/go-get-v2/parser/tencent"
	"github.com/schwarzeni/go-get-v2/parser/youku"
)

////////////////////////

// 解析器代号
var ParserMap = map[string]model.Parser{
	"0": bilibili.BilibiliParser{},
	"1": bilibili.BilibiliParser{IsUsePlugin: true},
	"2": youku.YouKuParser{},
	"3": iqiyi.IqiyiParser{},
	"4": tencent.TencentParser{}}
