package model

import "time"

const PLATFORM_STEAM = 1
const PLATFORM_EPIC = 2

type GameInfo struct {
	GameId int
	Name string
	CreateTime time.Time
	UpdateTime time.Time

	SteamPrice int
	SteamOriPrice int
	SteamDiscount int
	SteamLink string
	SteamImgLink string

	EpicPrice int
}

type GamePrice struct {
	PriceId int
	GameId int
	CreateTime time.Time

	Price int
	Platform int
}
