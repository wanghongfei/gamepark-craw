package shanguo

import (
	"fmt"
	"testing"
)

func TestCC(t *testing.T) {
	cc()
}

func TestCrawler_FetchHtml(t *testing.T) {
	cl := new(Crawler)
	content, err := cl.fetchHtml("https://www.sonkwo.com/store/search")
	if nil != err {
		t.Fatal(err)
	}
	fmt.Println(content)

	content, err = cl.fetchHtml("https://www.sonkwo.com/store/search?page=2")
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
	cl.CrawlGameInfo(1, 1, nil, nil)
}

