package main

import (
	"flag"
	"github.com/wanghongfei/gamepark-craw/crawl/shanguo"
	"github.com/wanghongfei/gamepark-craw/crawl/steam"
	"log"
	"os"
)

func main() {
	initLog()
	log.Println("version 1.3")

	// 命令行参数
	var outputFileName string
	var startPage int
	var concurrentPage int
	var target string
	flag.StringVar(&outputFileName, "output", "steam.tsv", "output file path")
	flag.IntVar(&startPage, "start", 1, "start page")
	flag.IntVar(&concurrentPage, "concurrency", 1, "page crawl concurrency")
	// flag.IntVar(&concurrentPage, "imageconcurrency", 1, "page crawl concurrency")
	flag.StringVar(&target, "target", "steam", "target website, steam/shanguo")
	flag.Parse()

	if "steam" == target {
		steam.CrawlSteam(outputFileName, startPage, concurrentPage)

	} else if "shanguo" == target {
		shanguo.CrawlShanguo(outputFileName, startPage, concurrentPage)

	} else {
		log.Println("invaild target")
		os.Exit(1)
	}
}



func initLog()  {
	log.SetOutput(os.Stdout)
}
