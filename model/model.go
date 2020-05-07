package model

import "time"

const PLATFORM_STEAM = 1
const PLATFORM_EPIC = 2

type GameInfo struct {
	GameId int
	Name string
	NameCn string
	CreateTime time.Time
	UpdateTime time.Time

	SteamPrice int
	SteamOriPrice int
	SteamDiscount int
	SteamLink string
	SteamImgLink string

	SgPrice int
	SgOriPrice int
	SgDiscount int
	SgLink string
}

