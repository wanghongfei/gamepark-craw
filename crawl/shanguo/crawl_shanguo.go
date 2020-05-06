package shanguo

import (
	"bytes"
	"context"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/chromedp"
	"github.com/wanghongfei/gamepark-craw/crawl"
	"log"
	"strconv"
	"strings"
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
		gameHtmlNodes := doc.Find(".search-results li")
		gameHtmlNodes.Each(func(i int, selection *goquery.Selection) {
			imgNode := selection.Find(".listed-game-img img")
			gameName, _ := imgNode.Attr("title")

			discountStr := selection.Find(".game-discount").Text()
			oriPriceStr := selection.Find(".game-origin-price").Text()
			nowPriceStr := selection.Find(".game-sale-price").Text()

			log.Printf("%s\t%s\t%s\t%s\n", gameName, discountStr, oriPriceStr, nowPriceStr)
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

func cc() {
	//ctx := context.Background()
	//options := []chromedp.ExecAllocatorOption{
	//	chromedp.UserAgent(`Mozilla/5.0 (Windows NT 6.3; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/73.0.3683.103 Safari/537.36`),
	//}
	//options = append(options, chromedp.DefaultExecAllocatorOptions[:]...)
	//
	//c, cc := chromedp.NewExecAllocator(ctx, options...)
	//defer cc()
	// create context
	ctx, cancel := initChromeContext()
	defer cancel()

	// run task list
	var res string
	//var ua string
	// var res2 string
	err := chromedp.Run(ctx,
		chromedp.Navigate(`http://www.sonkwo.hk/store/search`),
		// chromedp.WaitVisible(`#content-wrapper > div > div.SK-store-search-container > div.search-block > div.search-left > ul`),
		chromedp.WaitVisible(`#content-wrapper > div > div.SK-store-search-container > div.search-block > div.search-left > ul > div > li:nth-child(1) > a > div.listed-game-content > p.tags > span:nth-child(1)`),
		//chromedp.InnerHTML(`#content-wrapper > div > div.SK-store-search-container > div.search-block > div.search-left > ul > div > li:nth-child(1) > a > div.listed-game-content > p.tags > span:nth-child(1)`, &res),
		// chromedp.Sleep(3 * time.Second),
		// chromedp.OuterHTML(`#content-wrapper > div > div.SK-store-search-container > div.search-block > div.search-left > ul`, &res),
		// chromedp.Sleep(10 * time.Second),
		chromedp.OuterHTML(`document.querySelector("body")`, &res, chromedp.ByJSPath),
	)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(strings.TrimSpace(res))
	//log.Println(strings.TrimSpace(ua))
	// log.Println(strings.TrimSpace(res2))
}
