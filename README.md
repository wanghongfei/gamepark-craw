# Steam/杉果商店爬虫

爬取Steam商店搜索页面上的全部游戏数据。使用方法：

```shell
Usage of ./gamepark-craw:
  -concurrency int
        page crawl concurrency (default 1)
  -output string
        output file path (default "steam.tsv")
  -start int
        start page (default 1)
  -target string
        target website, steam/shanguo (default "steam")
```



其中`start`表示从第一页开始爬取, `concurrency`表示最多同时爬取2个页面, `output`用来指定输出文件，`target`表示爬取的平台。

游戏数据以`tsv`格式存储，一行有6个token, 每个token用`\t`分割，格式：

```shell
游戏名\t现价\t原价\t打折幅度\t图片\t商店
```

Demo:

```
Atom Zombie Smasher	36	36	0	https://store.steampowered.com/app/55040/Atom_Zombie_Smasher/?snr=1_7_7_230_150_2500	https://media.st.dl.eccdnx.com/steam/apps/55040/header.jpg?t=1586734113
```



