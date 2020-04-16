package steam

import (
	"fmt"
	"gamepark-craw/model"
	"github.com/PuerkitoBio/goquery"
	"testing"
)

func TestCrawAllCategory(t *testing.T) {
	res, err := CrawAllCategory()
	if nil != err {
		t.Fatal(err)
	}

	fmt.Println(res)
}

func TestCrawGameInfo(t *testing.T) {
	err := CrawGameInfo(func(info model.GameInfo) {
		fmt.Printf("%s\t%d\n", info.Name, info.SteamPrice)
	})
	if nil != err {
		t.Fatal(err)
	}
}

func TestFindMaxPage(t *testing.T) {
	bodyReader, err := Get("https://store.steampowered.com/tags/zh-cn/%E5%8A%A8%E4%BD%9C/?tab=NewReleases")
	if nil != err {
		t.Fatal(err)
	}

	doc, err := goquery.NewDocumentFromReader(bodyReader)
	if nil != err {
		t.Fatal(err)
	}

	fmt.Println(findMaxPage(doc))
}
