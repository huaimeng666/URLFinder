package main

import (
	"github.com/huaimeng666/URLFinder/cmd"
	"github.com/huaimeng666/URLFinder/config"
	"github.com/huaimeng666/URLFinder/crawler"
	"io"
	"log"
)

func main() {
	log.SetOutput(io.Discard)
	cmd.Parse()
	config.Init()
	crawler.Run()
}
