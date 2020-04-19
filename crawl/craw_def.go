package crawl

import "github.com/wanghongfei/gamepark-craw/model"

// 爬取到一个完整游戏信息时的回调函数
type OnGameInfo func(info model.GameInfo)

// 爬取发生错误时的回调函数
type OnGameError func(link string, err error)

type GameCrawl interface {
	CrawlGameInfo(int, int, OnGameInfo, OnGameError) error
}

type HotGameCrawl interface {
	CrawlHotGames() ([]string, error)
}
