package main

import (
	"fmt"
	"gamepark-craw/crawl/steam"
	"gamepark-craw/model"
	"log"
	"os"
	"time"
)

func main() {
	initLog()

	now := time.Now().Format("20060102")
	file, err := os.OpenFile("out/steam-" + now + ".tsv2", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if nil != err {
		log.Fatal(err)
	}
	defer file.Close()

	headLine := "游戏名\t现价\t原价\t打折幅度\n"
	file.WriteString(headLine)

	crawler := new(steam.Crawler)
	err = crawler.CrawlGameInfo(1, func(info model.GameInfo) {
		line := fmt.Sprintf("%s\t%d\t%d\t%d\n", info.Name, info.SteamPrice, info.SteamOriPrice, info.SteamDiscount)
		_, werr := file.WriteString(line)
		if nil != werr {
			log.Printf("failed to write data to file: %+v", werr)
			panic(werr)
		}
		// fmt.Print(line)
	})

	if nil != err {
		log.Printf("%v", err)
	}
}

func initLog()  {
	log.SetOutput(os.Stdout)
}
