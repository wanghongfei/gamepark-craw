package main

import (
	"fmt"
	"gamepark-craw/model"
	"gamepark-craw/steam"
	"log"
	"os"
)

func main() {
	initLog()

	file, err := os.OpenFile("steam-20200415.tsv", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	// file, err := os.Create("steam.txt")
	if nil != err {
		log.Fatal(err)
	}
	defer file.Close()

	headLine := "游戏名\t现价\t原价\t打折幅度\n"
	file.WriteString(headLine)

	err = steam.CrawGameInfo(2073, func(info model.GameInfo) {
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
