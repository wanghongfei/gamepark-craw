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
	"strconv"
	"strings"
	"time"
)

const SHANGUO_HOST = "https://www.sonkwo.com"
const SHANGUO_PAGE = "https://www.sonkwo.com/store/search"

type Crawler struct {
	chromeContext context.Context
	cancelFunc context.CancelFunc

	withEngName bool
}

var SEARCH_RESULT_WAIT_EXPRESSION = `#content-wrapper > div > div.SK-store-search-container > div.search-block > div.search-left > ul > div > li:nth-child(1) > a > div.listed-game-content > p.tags > span:nth-child(1)`
var DETAIL_WAIT_EXPRESSION = `#content-wrapper > div > div > div.new-content-container > div`

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
		pageHtml, err := c.fetchHtml(link, SEARCH_RESULT_WAIT_EXPRESSION)
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
			gameName = strings.TrimSpace(gameName)

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

			if !strings.HasPrefix(detailLink, "http") {
				detailLink = SHANGUO_HOST + detailLink
			}

			info := &model.GameInfo{
				GameId:        0,
				// Name:		   '',
				NameCn:        gameName,
				CreateTime:    time.Now(),
				SgPrice:       nowPrice,
				SgOriPrice:    oriPrice,
				SgDiscount:    discount,
				SgLink:        detailLink,
			}

			// 爬取英文名
			if c.withEngName {
				engName, err := c.fetchEngName(detailLink)
				if nil != err {
					log.Printf("failed to fetch eng name in page %s, %v\n", engName, err)
					return
				}

				info.Name = engName
			}



			infos = append(infos, info)

			onInfo(*info)
		})
	}

	c.cancelFunc()
	c.chromeContext = nil
	c.cancelFunc = nil

	return nil
}

func (c *Crawler) fetchEngName(link string) (string, error) {
	detailHtml, err := c.fetchHtml(link, DETAIL_WAIT_EXPRESSION)
	if nil != err {
		return "", fmt.Errorf("failed to fetch eng name in link %s, %w", link, err)
	}

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader([]byte(detailHtml)))
	if nil != err {
		return "", fmt.Errorf("failed to parse detail page, %w", err)
	}

	return doc.Find(".typical-name-2").Text(), nil
}

func (c *Crawler) fetchMaxPage() (int, error) {
	pageHtml, err := c.fetchHtml(SHANGUO_PAGE, SEARCH_RESULT_WAIT_EXPRESSION)
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
	options := []chromedp.ExecAllocatorOption{
		//chromedp.Flag("headless", false), // debug使用
		chromedp.Flag("blink-settings", "imagesEnabled=false"),
		chromedp.UserAgent(`Mozilla/5.0 (Windows NT 6.3; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/73.0.3683.103 Safari/537.36`),
	}
	options = append(chromedp.DefaultExecAllocatorOptions[:], options...)

	c, _ := chromedp.NewExecAllocator(context.Background(), options...)
	// defer allocatorCancel()

	// create context
	chromeCtx, cancel := chromedp.NewContext(c, chromedp.WithLogf(log.Printf))
	// 执行一个空task, 用提前创建Chrome实例
	chromedp.Run(chromeCtx, make([]chromedp.Action, 0, 1)...)

	return chromeCtx, cancel
}

func (c *Crawler) fetchHtml(link string, waitExpression string) (string, error) {
	// 如果是第一次调用
	if nil == c.chromeContext {
		// 初始化chrome
		log.Println("init chrome")
		c.chromeContext, c.cancelFunc = initChromeContext()
		log.Println("done initialization")
	}

	// 给每个页面的爬取设置超时时间
	timeoutCtx, cancel := context.WithTimeout(c.chromeContext, 30 * time.Second)
	defer cancel()


	log.Printf("Chrome visit page %s\n", link)

	var htmlContent string
	err := chromedp.Run(timeoutCtx,
		chromedp.Navigate(link),
		chromedp.WaitVisible(waitExpression),
		chromedp.OuterHTML(`document.querySelector("body")`, &htmlContent, chromedp.ByJSPath),
	)
	if err != nil {
		return "", err
	}

	return htmlContent, nil
}

