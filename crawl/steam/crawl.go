package steam

import (
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/wanghongfei/gamepark-craw/crawl"
	"github.com/wanghongfei/gamepark-craw/model"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"
)


type Crawler struct {

}

var urlTemplate = "https://store.steampowered.com/search/?sort_by=Released_DESC&page=%d"
var hotLink = "https://store.steampowered.com/stats/"

var ERR_ACCESS_DENIED = fmt.Errorf("ACCESS DENIED!!!")


// 爬取steam搜索页面全部游戏数据;
// 可以指定从startPage页开始爬取，startPage从1开始计算；
// 可以指定最多同时爬取concurrentPageAmount个页面
// onGameInfoFunc: 当爬取到一条完整的游戏信息时回调次函数
func (cl *Crawler) CrawlGameInfo(startPage int, concurrentPageAmount int, onGameInfoFunc crawl.OnGameInfo, onGameError crawl.OnGameError) error {
	// 获取最大页码
	lastPage, err := fetchMaxPage()
	if nil != err {
		return err
	}

	// 任务队列
	taskChan := make(chan int, concurrentPageAmount)

	// 初始化协程池
	wg := new(sync.WaitGroup)
	for ix := 0; ix < concurrentPageAmount; ix++ {
		wg.Add(1)
		go crawlTaskRoutine(wg, ix + 1, taskChan, onGameInfoFunc, onGameError)
	}

	// 开始投递任务
	for page := startPage; page <= lastPage ; page++ {
		taskChan <- page
	}

	// 关闭任务队列
	close(taskChan)
	// 等待任务全部结束
	wg.Wait()

	return nil
}

// 爬取热门游戏
// 返回游戏名称列表, 按热度排序
func (cl *Crawler) CrawlHotGames() ([]string, error) {
	log.Printf("getting host game page")

	bodyReader, err := crawl.GetWithRetry(hotLink, 2)
	if nil != err {
		return nil, fmt.Errorf("failed to request hot game page: %w", err)
	}

	doc, err := goquery.NewDocumentFromReader(bodyReader)
	if err != nil {
		return nil, fmt.Errorf("failed to parse steam host page html: %w", err)
	}

	// 找到游戏列表
	gameNodes := doc.Find("#detailStats .gameLink")
	if 0 == gameNodes.Length() {
		return nil, errors.New("no hot game elements found")
	}

	names := make([]string, 0, gameNodes.Length())
	gameNodes.Each(func(i int, selection *goquery.Selection) {
		name := selection.Text()
		names = append(names, name)
	})

	log.Printf("%d hot games found\n", gameNodes.Length())
	return names, nil
}

func crawlTaskRoutine(wg *sync.WaitGroup, id int, taskPageChan chan int, onGameInfoFunc crawl.OnGameInfo, onError crawl.OnGameError) {
	log.Printf("routine %d started\n", id)

	// 等待任务队列里的新任务
	for page := range taskPageChan {
		err := crawlLink(page, onGameInfoFunc)
		if nil != err {
			log.Printf("failed to crawl page %d, %v skip\n", page, err)
			if ERR_ACCESS_DENIED == err {
				// 被封禁了, 退出
				panic(err)
			}

			// 回调
			onError(fmt.Sprintf(urlTemplate, page), err)
		}
	}

	log.Printf("routine %d exited\n", id)
	wg.Done()
}

func fetchMaxPage() (int, error) {
	log.Printf("getting max page")

	link := fmt.Sprintf(urlTemplate, 1)
	bodyReader, err := crawl.GetWithRetry(link, 2)
	if nil != err {
		return -1, fmt.Errorf("failed to request list page: %w", err)
	}

	doc, err := goquery.NewDocumentFromReader(bodyReader)
	if err != nil {
		return -1, fmt.Errorf("failed to parse steam list page html: %w", err)
	}

	// 找出最大页码
	lastPageStr := doc.Find("#search_result_container .search_pagination_right").Children().Last().Prev().Text()
	log.Println("max page number is " + lastPageStr)
	maxPage, err := strconv.Atoi(lastPageStr)
	if nil != err {
		return -1, fmt.Errorf("invalid last page string: %w", err)
	}

	return maxPage, nil

}


func crawlLink(page int, onGameInfoFunc crawl.OnGameInfo) error {
	startTime := time.Now()

	log.Printf("crawl page: %d\n", page)

	link := fmt.Sprintf(urlTemplate, page)
	bodyReader, err := crawl.GetWithRetry(link, 2)
	if nil != err {
		return fmt.Errorf("failed to request list page: %w", err)
	}

	doc, err := goquery.NewDocumentFromReader(bodyReader)
	if err != nil {
		return fmt.Errorf("failed to parse steam list page html: %w", err)
	}

	// 判断有没有被封禁
	title := doc.Find("title").Text()
	if "Access Denied" == title {
		return ERR_ACCESS_DENIED
	}

	// 找到游戏节点
	gameNode := doc.Find("#search_result_container #search_resultsRows .responsive_search_name_combined")
	if len(gameNode.Nodes) == 0 {
		// 没找到游戏节点, 说明页面报错了或者改版了
		return fmt.Errorf("no game found on page %s", link)
	}

	detailParseResultChan := make(chan *model.GameInfo, 3 * gameNode.Length())
	// 启动大图抓取routine
	imageTaskChan := make(chan *model.GameInfo, 3)
	for ix := 0; ix < 3; ix++ {
		go extractBigImageRoutine(imageTaskChan, detailParseResultChan)
	}
	fmt.Printf("%v image routine started\n", 3)

	// 遍历游戏节点
	gameCounter := 0
	var innerErr error
	gameNode.Each(func(i int, s *goquery.Selection) {
		name := s.Find(".search_name .title").Text()
		priceNode := s.Find(".search_price_discount_combined .search_price")

		oriPriceStr := ""
		discountPriceStr := ""
		discountPercentStr := ""

		// 判断是否打折
		exist := priceNode.HasClass("discounted")
		if exist {
			// 原价
			oriPriceStr = priceNode.Find("strike").Text()
			priceNode.Find("span").Remove()
			// 打折后价格
			discountPriceStr = priceNode.Text()
			// 打折幅度
			discountPercentStr = priceNode.Prev().Find("span").Text()

		} else {
			oriPriceStr = priceNode.Text()
			discountPriceStr = oriPriceStr
			discountPercentStr = "0%"
		}


		// 解析价格
		oriPrice, err := parsePrice(oriPriceStr)
		if nil != err {
			innerErr = fmt.Errorf("invalid game price %s for %s: %v\n", oriPriceStr, name, err)
			return
		}
		discountPrice, err := parsePrice(discountPriceStr)
		if nil != err {
			innerErr = fmt.Errorf("invalid game price %s for %s: %v\n", discountPriceStr, name, err)
			return
		}
		discountPercent, err := parsePrice(discountPercentStr)
		if nil != err {
			innerErr = fmt.Errorf("invalid game price %s for %s: %v\n", discountPercentStr, name, err)
			return
		}

		// 详情页面链接
		detailLink, _ := s.Parent().Attr("href")

		info := &model.GameInfo{
			Name:       name,
			CreateTime: time.Time{},
			UpdateTime: time.Time{},
			SteamPrice: discountPrice,
			SteamOriPrice: oriPrice,
			SteamDiscount: discountPercent,
			SteamLink: detailLink,
			// SteamImgLink: imgLink,
		}

		// 投递抓取详情页的大图按任务
		imageTaskChan <- info

		//go func(info *model.GameInfo, link string) {
		//	imgLink := extractBigImage(link)
		//	info.SteamImgLink = imgLink
		//
		//	detailParseResultChan <- info
		//}(info, detailLink)


		gameCounter++
	})
	if nil != innerErr {
		return innerErr
	}

	close(imageTaskChan)

	// 等待并发任务完成
	for ix := 0; ix < gameCounter; ix++ {
		info := <- detailParseResultChan
		onGameInfoFunc(*info)
	}
	close(detailParseResultChan)

	endTime := time.Now()
	diffTime := endTime.Sub(startTime).Milliseconds()
	log.Printf("done parsing %s, game count = %d, cost = %dms\n", link, gameCounter, diffTime)


	return nil
}

func extractBigImageRoutine(taskChan chan *model.GameInfo, resultChan chan *model.GameInfo) {
	for info := range taskChan {
		info.SteamImgLink = extractBigImage(info.SteamLink)
		resultChan <- info
	}
}

func extractBigImage(url string) string {
	bodyReader, err := crawl.GetWithRetry(url, 2)
	if nil != err {
		log.Printf("failed to open detail page:%v", err)
		return ""
	}

	doc, err := goquery.NewDocumentFromReader(bodyReader)
	if err != nil {
		log.Printf("failed to parse detail page html: %v", err)
		return ""
	}

	imgLink, _ := doc.Find(".game_header_image_full").Attr("src")
	return imgLink

}

func parsePrice(str string) (int, error) {
	str = strings.TrimSpace(str)
	str = strings.ReplaceAll(str, ",", "")
	str = strings.ReplaceAll(str, "%", "")
	if "" == str {
		// 没有价格
		return -1, nil
	}

	if strings.Contains(str, "免费") || strings.Contains(str, "Free") || "0" == str {
		// 免费
		return 0, nil
	}

	priceIndex := strings.LastIndex(str, "¥")
	priceToken := strings.TrimSpace(str[priceIndex + 2:])
	price, err := strconv.ParseFloat(priceToken, 64)
	if nil != err {
		return -2, fmt.Errorf("invalid price string: %w", err)
	}

	return int(price), nil
}

