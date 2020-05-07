package shanguo

import (
	"fmt"
	"github.com/wanghongfei/gamepark-craw/model"
	"log"
	"testing"
)

func TestCrawler_FetchHtml(t *testing.T) {
	cl := new(Crawler)
	content, err := cl.fetchHtml("https://www.sonkwo.com/store/search", SEARCH_RESULT_WAIT_EXPRESSION)
	if nil != err {
		t.Fatal(err)
	}
	fmt.Println(content)

	content, err = cl.fetchHtml("http://www.sonkwo.hk/sku/2925", DETAIL_WAIT_EXPRESSION)
	if nil != err {
		t.Fatal(err)
	}
	fmt.Println(content)
}

func TestFetchMaxPage(t *testing.T) {
	cl := new(Crawler)
	page, err := cl.fetchMaxPage()
	if nil != err {
		t.Fatal(err)
	}

	fmt.Println(page)
}

func TestCrawler_CrawlGameInfo(t *testing.T) {
	cl := new(Crawler)
	cl.CrawlGameInfo(1, 1, func(info model.GameInfo) {
		log.Printf("%v\t%v\t%v\t%v\n", info.Name, info.SgDiscount, info.SgOriPrice, info.SgPrice)
	}, nil)
}

