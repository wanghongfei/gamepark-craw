package shanguo

import (
	"bytes"
	"context"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/chromedp"
	"github.com/wanghongfei/gamepark-craw/crawl"
	"github.com/wanghongfei/gamepark-craw/model"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

const SHANGUO_PAGE = "https://www.sonkwo.com/store/search"

type Crawler struct {
	chromeContext context.Context
	cancelFunc context.CancelFunc
}

func (c *Crawler) CrawlGameInfo(startPage int, concurrentPageAmount int, onInfo crawl.OnGameInfo, onError crawl.OnGameError) error {
	maxPage, err := c.fetchMaxPage()
	if nil != err {
		return fmt.Errorf("failed to fetch max page, %w", err)
	}
	log.Printf("max page number is %d\n", maxPage)

	for page := startPage; page <= maxPage; page++ {
		log.Printf("crawling page %d\n", page)

		// 爬取页面完整html
		link := SHANGUO_PAGE + "?page=" + strconv.Itoa(page)
		pageHtml, err := c.fetchHtml(link)
		if nil != err {
			log.Printf("failed to visit page %s, %v\n", link, err)
			continue
		}


		// 解析页面
		htmlReader := bytes.NewReader([]byte(pageHtml))
		doc, err := goquery.NewDocumentFromReader(htmlReader)
		if nil != err {
			log.Printf("failed to parse search result page %s, %v\n", link, err)
			continue
		}

		// 游戏信息节点
		infos := make([]*model.GameInfo, 0, 20)
		gameHtmlNodes := doc.Find(".search-results li")
		gameHtmlNodes.Each(func(i int, selection *goquery.Selection) {
			// 游戏名
			imgNode := selection.Find(".listed-game-img img")
			gameName, _ := imgNode.Attr("title")

			// 详情页地址
			detailLink, _ := selection.Find(".listed-game-block").Attr("href")

			// 价格
			discountStr := selection.Find(".game-discount").Text()
			oriPriceStr := selection.Find(".game-origin-price").Text()
			nowPriceStr := selection.Find(".game-sale-price").Text()

			oriPrice, err := parsePrice(oriPriceStr)
			if nil != err {
				log.Printf("failed to parse ori price %s, %v\n", oriPriceStr, err)
				return
			}
			nowPrice, err := parsePrice(nowPriceStr)
			if nil != err {
				log.Printf("failed to parse now price %s, %v\n", nowPriceStr, err)
				return
			}
			discount, err := parseDiscount(discountStr)
			if nil != err {
				log.Printf("failed to parse discount %s, %v\n", discountStr, err)
				return
			}


			info := &model.GameInfo{
				GameId:        0,
				Name:          gameName,
				CreateTime:    time.Now(),
				SgPrice:       nowPrice,
				SgOriPrice:    oriPrice,
				SgDiscount:    discount,
				SgLink:        detailLink,
			}

			infos = append(infos, info)

			onInfo(*info)
			// log.Printf("%s\t%s\t%s\t%s\n", gameName, discountStr, oriPriceStr, nowPriceStr)
		})
	}

	return nil
}

func (c *Crawler) fetchMaxPage() (int, error) {
	pageHtml, err := c.fetchHtml(SHANGUO_PAGE)
	if nil != err {
		return 0, err
	}

	htmlReader := bytes.NewReader([]byte(pageHtml))
	doc, err := goquery.NewDocumentFromReader(htmlReader)
	if nil != err {
		return 0, fmt.Errorf("failed to parse shanguo page, %w", err)
	}

	maxPage := 0
	doc.Find(".SK-pagedown-list .item").Last().Each(func(i int, selection *goquery.Selection) {
		maxPage, err = strconv.Atoi(selection.Text())
	})

	if nil != err {
		return 0, fmt.Errorf("failed to find max page number, %w", err)
	}

	return maxPage, nil
}

func parseDiscount(disStr string) (int, error) {
	if "" == disStr {
		return 0, nil
	}

	return strconv.Atoi(disStr[1:len(disStr) - 1])
}

func parsePrice(priceStr string) (int, error) {
	if "" == priceStr {
		return 0, nil
	}

	priceStr = strings.ReplaceAll(priceStr, "￥", "")

	ptIdx := strings.Index(priceStr, ".")
	numStr := priceStr[0:ptIdx]

	return strconv.Atoi(numStr)
}

func initChromeContext() (context.Context, context.CancelFunc) {
	ctx := context.Background()
	options := []chromedp.ExecAllocatorOption{
		chromedp.UserAgent(`Mozilla/5.0 (Windows NT 6.3; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/73.0.3683.103 Safari/537.36`),
	}
	options = append(options, chromedp.DefaultExecAllocatorOptions[:]...)

	c, _ := chromedp.NewExecAllocator(ctx, options...)
	// defer cc()

	// create context
	return chromedp.NewContext(c)
}

func (c *Crawler) fetchHtml(link string) (string, error) {
	if nil == c.chromeContext {
		log.Println("init chrome")
		c.chromeContext, c.cancelFunc = initChromeContext()
	}

	log.Printf("Chrome visit page %s\n", link)

	var htmlContent string
	err := chromedp.Run(c.chromeContext,
		chromedp.Navigate(link),
		chromedp.WaitVisible(`#content-wrapper > div > div.SK-store-search-container > div.search-block > div.search-left > ul > div > li:nth-child(1) > a > div.listed-game-content > p.tags > span:nth-child(1)`),
		chromedp.OuterHTML(`document.querySelector("body")`, &htmlContent, chromedp.ByJSPath),
	)
	if err != nil {
		log.Fatal(err)
	}

	return strings.TrimSpace(htmlContent), nil
}

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
		line := fmt.Sprintf("%s\t%d\t%d\t%d\t%s\n", info.Name, info.SgPrice, info.SgOriPrice, info.SgDiscount, info.SgLink)
		_, werr := file.WriteString(line)
		if nil != werr {
			log.Printf("failed to write data to file: %+v", werr)
			panic(werr)
		}
	}

	// 创建爬虫
	var crawler crawl.GameCrawl
	crawler = new(Crawler)
	// 启动爬虫
	err = crawler.CrawlGameInfo(startPage, concurrentPage, onSuccess, nil)
	if nil != err {
		log.Printf("%v", err)
	}

}
