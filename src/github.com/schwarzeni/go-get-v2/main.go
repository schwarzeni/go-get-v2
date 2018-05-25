package main

import (
	"github.com/schwarzeni/go-get-v2/core/engine"
	"github.com/schwarzeni/go-get-v2/parser/bilibili"
)

func main() {
	engine.Run(10, bilibili.BilibiliParser{})
}
