package steam


import (
	"fmt"
	"github.com/wanghongfei/gamepark-craw/crawl"
	"github.com/wanghongfei/gamepark-craw/model"
	"log"
	"os"
)

func CrawlSteam(outputFileName string, startPage int, concurrentPage int) {
	// 打开结果输出文件
	log.Printf("send data to %s, start page %d, max concurrency page count %d, \n", outputFileName, startPage, concurrentPage)
	file, err := os.OpenFile(outputFileName, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if nil != err {
		log.Fatal(err)
	}
	defer file.Close()

	// 打开失败输出文件
	errFile, err := os.OpenFile("failed.txt", os.O_CREATE|os.O_RDWR, 0666)
	if nil != err {
		log.Fatal(err)
	}
	defer errFile.Close()

	// 标题行
	headLine := "游戏名\t现价\t原价\t打折幅度\t图片\t商店\n"
	file.WriteString(headLine)

	// 定义回调函数
	// 成功函数
	onSuccess := func(info model.GameInfo) {
		// 输出到文件
		line := fmt.Sprintf("%s\t%d\t%d\t%d\t%s\t%s\n", info.Name, info.SteamPrice, info.SteamOriPrice, info.SteamDiscount, info.SteamLink, info.SteamImgLink)
		_, werr := file.WriteString(line)
		if nil != werr {
			log.Printf("failed to write data to file: %+v", werr)
			panic(werr)
		}
	}
	// 失败函数
	onFailed := func(link string, err error) {
		line := fmt.Sprintf("%s\t%v\n", link, err)
		errFile.WriteString(line)
	}

	// 创建爬虫
	var crawler crawl.GameCrawl
	crawler = new(Crawler)
	// 启动爬虫
	err = crawler.CrawlGameInfo(startPage, concurrentPage, onSuccess, onFailed)

	if nil != err {
		log.Printf("%v", err)
	}

	// 爬取热门游戏列表
	err = saveHotGames(crawler, outputFileName)
	if nil != err {
		log.Printf("%v", err)
	}

}

func saveHotGames(crawler crawl.GameCrawl, hotFileName string) error {
	// 爬取热门游戏列表
	hotGameCrawler, ok := crawler.(crawl.HotGameCrawl)
	if !ok {
		return nil
	}
	hotGames, err := hotGameCrawler.CrawlHotGames()
	if nil != err {
		return fmt.Errorf("failed to crawl hot games, %v\n", err)
	}

	// 写入到热门文件
	hotFiles, err := os.OpenFile(hotFileName + ".hot", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if nil != err {
		return fmt.Errorf("failed to open hot file, %v", err)
	}
	defer hotFiles.Close()

	for _, name := range hotGames {
		hotFiles.WriteString(name + "\n")
	}

	return nil
}

