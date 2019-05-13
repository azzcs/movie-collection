package collection

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"github.com/gocolly/colly"
	"log"
	"movie-collection/config"
	"regexp"
	"strings"
	"time"
	"strconv"
)

func MakeDetailCollection(listChannel <-chan string) <-chan *Movie{
	detailChannel := make(chan *Movie,100)
	channel := make(chan bool)
	go func(listChannel <-chan string,detailChannel chan<- *Movie,channel chan bool) {
		for i:=0;i<config.ThreadNum;i++{
			go analysisDetailChannel(listChannel,detailChannel,channel)
		}
		for i:=0;i<config.ThreadNum;i++{
			<-channel
		}
		close(channel)
		close(detailChannel)
	}(listChannel,detailChannel,channel)
	return detailChannel
}


func analysisDetailChannel(listChannel <-chan string,detailChannel chan<- *Movie,channel chan<- bool){
	for url:= range listChannel{
		detailChannel<-detailCollection(url)
	}
	channel<-true
}
func detailCollection(url string) *Movie{
	cDetail := colly.NewCollector(
		colly.AllowURLRevisit(),
		colly.URLFilters(
			regexp.MustCompile(`http://www.okzyw.com/\?m=vod-detail-id-\d+.html`),
		),
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36 Edge/16.16299"),
		colly.ParseHTTPErrorResponse(),
	)
	cDetail.WithTransport(config.Transport)

	var movie Movie
	cDetail.OnError(func(r *colly.Response, err error) {
		// TODO 增加失败重试次数限制
		log.Println("详情页错误URL:", r.Request.URL,err)
		//time.Sleep(time.Second*5)
		//r.Request.Visit(r.Request.URL.Path)
		//log.Println("重新采集详情页错误URL:", r.Request.URL,err)
	})
	cDetail.OnHTML("div.nvc", func(e *colly.HTMLElement) { //详情页-分类
		typeName := e.DOM.Find("dd>a").Eq(1).Text()
		movie = Movie{
			Nav2: typeName,
		}
		movie.InitNav1()
	})
	cDetail.OnHTML("div.vodBox", func(e *colly.HTMLElement) { //详情页-分类
		name := e.ChildText("h2")
		updateStatus := e.ChildText("div.vodh>span")
		img := e.ChildAttr("div.vodImg>img", "src")
		movie.Name = name
		movie.UpdateStatus = updateStatus
		movie.Img = img
		e.ForEach("div.vodinfobox>ul>li", func(i int, element *colly.HTMLElement) {
			key := element.Text
			value := element.ChildText("span")
			if strings.HasPrefix(key, "别名：") {
				movie.Alias = value
			} else if strings.HasPrefix(key, "导演：") {
				movie.Director = value
			} else if strings.HasPrefix(key, "主演：") {
				movie.Actor = value
			} else if strings.HasPrefix(key, "类型：") {
				movie.Type = value
			} else if strings.HasPrefix(key, "地区：") {
				movie.Area = value
			} else if strings.HasPrefix(key, "语言：") {
				movie.Language = value
			} else if strings.HasPrefix(key, "上映：") {
				b,err := strconv.Atoi(value)
				year:=time.Now().Year()
				if err==nil && b>1500 && b<year {
					movie.Time = value
				} 
			} else if strings.HasPrefix(key, "更新：") {
				movie.UpdateTime = value
			}
		})
	})
	cDetail.OnHTML("div.warp", func(e *colly.HTMLElement) { //详情页-简介
		e.ForEach("div.playBox", func(i int, element *colly.HTMLElement) {
			textContent := element.Text
			if strings.Contains(textContent, "剧情介绍：") {
				movie.Introduce = element.ChildText("div.vodplayinfo")
			} else if strings.Contains(textContent, "播放类型：") {
				var content = map[string][]string{}
				kuyun := []string{}
				ckm3u8 := []string{}
				element.ForEach("div#1>ul>li", func(i int, eKuyun *colly.HTMLElement) {
					kuyunContent := eKuyun.Text
					if kuyunContent != "" {
						kuyun = append(kuyun, kuyunContent)
					}
				})
				element.ForEach("div#2>ul>li", func(i int, eKuyun *colly.HTMLElement) {
					kuyunContent := eKuyun.Text
					if kuyunContent != "" {
						ckm3u8 = append(ckm3u8, kuyunContent)
					}
				})
				content["ckm3u8"] = ckm3u8
				content["kuyun"] = kuyun
				mjson, _ := json.Marshal(content)
				movie.Content = string(mjson)
			}
		})

		h := md5.New()
		h.Write([]byte(movie.Name))
		cipherStr := h.Sum(nil)
		id := hex.EncodeToString(cipherStr)
		movie.Id = id
	})
	cDetail.Visit(url)
	return &movie
}
