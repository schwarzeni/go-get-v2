package main

import (
	"github.com/schwarzeni/go-get-v2/parser/bilibili"
	"github.com/schwarzeni/go-get-v2/scheduler"
)

func main() {
	scheduler.Engine(2, bilibili.BilibiliParser{})

}
