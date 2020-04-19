package shanguo

//func cc() {
//	ctx := context.Background()
//	options := []chromedp.ExecAllocatorOption{
//		chromedp.UserAgent(`Mozilla/5.0 (Windows NT 6.3; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/73.0.3683.103 Safari/537.36`),
//	}
//	options = append(options, chromedp.DefaultExecAllocatorOptions[:]...)
//
//	c, cc := chromedp.NewExecAllocator(ctx, options...)
//	defer cc()
//	// create context
//	ctx, cancel := chromedp.NewContext(c)
//	defer cancel()
//
//	// run task list
//	var res string
//	//var ua string
//	// var res2 string
//	err := chromedp.Run(ctx,
//		chromedp.Navigate(`http://www.sonkwo.hk/store/search`),
//		// chromedp.WaitVisible(`#content-wrapper > div > div.SK-store-search-container > div.search-block > div.search-left > ul`),
//		chromedp.WaitVisible(`#content-wrapper > div > div.SK-store-search-container > div.search-block > div.search-left > ul > div > li:nth-child(1) > a > div.listed-game-content > p.tags > span:nth-child(1)`),
//		//chromedp.InnerHTML(`#content-wrapper > div > div.SK-store-search-container > div.search-block > div.search-left > ul > div > li:nth-child(1) > a > div.listed-game-content > p.tags > span:nth-child(1)`, &res),
//		// chromedp.Sleep(3 * time.Second),
//		// chromedp.OuterHTML(`#content-wrapper > div > div.SK-store-search-container > div.search-block > div.search-left > ul`, &res),
//		// chromedp.Sleep(10 * time.Second),
//		chromedp.OuterHTML(`document.querySelector("body")`, &res, chromedp.ByJSPath),
//	)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	log.Println(strings.TrimSpace(res))
//	//log.Println(strings.TrimSpace(ua))
//	// log.Println(strings.TrimSpace(res2))
//}
//
//func cc2() {
//	var ua string
//
//	ctx := context.Background()
//	options := []chromedp.ExecAllocatorOption{
//		chromedp.UserAgent(`Mozilla/5.0 (Windows NT 6.3; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/73.0.3683.103 Safari/537.36`),
//	}
//	//options = append(options, chromedp.DefaultExecAllocatorOptions[:]...)
//
//	c, cc := chromedp.NewExecAllocator(ctx, options...)
//	defer cc()
//	// create context
//	ctx, cancel := chromedp.NewContext(c)
//	defer cancel()
//
//	err := chromedp.Run(ctx,
//		chromedp.Navigate(`https://www.whatsmyua.info/?a`),
//		chromedp.WaitVisible(`#custom-ua-string`),
//		chromedp.Text(`#custom-ua-string`, &ua),
//		chromedp.Sleep(10 * time.Second),
//	)
//	if err != nil {
//		log.Fatal(err)
//	}
//	log.Printf("user agent: %s", ua)
//
//}
