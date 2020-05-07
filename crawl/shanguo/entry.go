package shanguo

import (
	"fmt"
	"github.com/wanghongfei/gamepark-craw/crawl"
	"github.com/wanghongfei/gamepark-craw/model"
	"log"
	"os"
)

func CrawlShanguo(outputFileName string, startPage int, concurrentPage int) {
	// 打开结果输出文件
	log.Printf("send data to %s, start page %d, max concurrency page count %d, \n", outputFileName, startPage, concurrentPage)
	file, err := os.OpenFile(outputFileName, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if nil != err {
		log.Fatal(err)
	}
	defer file.Close()

	// 定义回调函数
	// 成功函数
	onSuccess := func(info model.GameInfo) {
		// 输出到文件
		line := fmt.Sprintf("%s\t%s\t%d\t%d\t%d\t%s\n", info.Name, info.NameCn, info.SgPrice, info.SgOriPrice, info.SgDiscount, info.SgLink)
		_, werr := file.WriteString(line)
		if nil != werr {
			log.Printf("failed to write data to file: %+v", werr)
			panic(werr)
		}
	}

	// 创建爬虫
	var crawler crawl.GameCrawl
	crawler = &Crawler{
		withEngName:   true,
	}
	// 启动爬虫
	err = crawler.CrawlGameInfo(startPage, concurrentPage, onSuccess, nil)
	if nil != err {
		log.Printf("%v", err)
	}

}
