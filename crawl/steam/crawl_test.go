package steam

import (
	"fmt"
	"testing"
)

func TestCrawler_CrawlHotGames(t *testing.T) {
	cl := new(Crawler)

	fmt.Println(cl.CrawlHotGames())
}
