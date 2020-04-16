package crawl

import "gamepark-craw/model"

// 爬取到一游戏信息时的回调函数
type OnGameInfo func(info model.GameInfo)

type GameCrawl interface {
	CrawlGameInfo(startPage int, onGameInfoFunc OnGameInfo) error
}
