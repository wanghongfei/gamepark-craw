package steam

import (
	"fmt"
	"gamepark-craw/crawl"
	"gamepark-craw/model"
	"github.com/PuerkitoBio/goquery"
	"log"
	"strconv"
	"strings"
	"time"
)

type Category struct {
	CategoryName string
	Link string
}

type Crawler struct {

}


func (cl *Crawler) CrawlGameInfo(startPage int, onGameInfoFunc crawl.OnGameInfo) error {
	urlTemplate := "https://store.steampowered.com/search/?sort_by=Released_DESC&page=%d"
	lastPage := -1
	for page := startPage; ; page++ {
		startTime := time.Now()

		log.Printf("current page: %d\n", page)
		if page == lastPage {
			break
		}

		link := fmt.Sprintf(urlTemplate, page)
		bodyReader, err := crawl.GetWithRetry(link, 2)
		if nil != err {
			return fmt.Errorf("failed to request list page: %w", err)
		}

		doc, err := goquery.NewDocumentFromReader(bodyReader)
		if err != nil {
			return fmt.Errorf("failed to parse steam list page html: %w", err)
		}

		if -1 == lastPage {
			// 找出最大页码
			lastPageStr := doc.Find("#search_result_container .search_pagination_right").Children().Last().Prev().Text()
			log.Println("max page number is " + lastPageStr)
			lastPage, err = strconv.Atoi(lastPageStr)
			if nil != err {
				return fmt.Errorf("invalid last page string: %w", err)
			}
		}

		// 找到游戏节点
		gameCounter := 0
		doc.Find("#search_result_container #search_resultsRows .responsive_search_name_combined").Each(func(i int, s *goquery.Selection) {
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
				log.Printf("invalid game price %s for %s: %+v\n", oriPriceStr, name, err)
			}
			discountPrice, err := parsePrice(discountPriceStr)
			if nil != err {
				log.Printf("invalid game price %s for %s: %+v\n", discountPriceStr, name, err)
			}
			discountPercent, err := parsePrice(discountPercentStr)
			if nil != err {
				log.Printf("invalid game price %s for %s: %+v\n", discountPercentStr, name, err)
			}

			// 抓取详情页面链接
			detailLink, _ := s.Parent().Attr("href")
			imgLink := extractBigImage(detailLink)

			info := &model.GameInfo{
				Name:       name,
				CreateTime: time.Time{},
				UpdateTime: time.Time{},
				SteamPrice: discountPrice,
				SteamOriPrice: oriPrice,
				SteamDiscount: discountPercent,
				SteamLink: detailLink,
				SteamImgLink: imgLink,
				EpicPrice:  0,
			}

			onGameInfoFunc(*info)
			gameCounter++

		})

		endTime := time.Now()
		diffTime := endTime.Sub(startTime).Milliseconds()
		log.Printf("done parsing %s, game count = %d, cost = %dms\n", link, gameCounter, diffTime)

	}


	return nil
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
	price, err := strconv.Atoi(priceToken)
	if nil != err {
		return -2, fmt.Errorf("invalid price string: %w", err)
	}

	return price, nil
}

