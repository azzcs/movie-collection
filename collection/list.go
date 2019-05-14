package collection

import (
	"github.com/gocolly/colly"
	"log"
	"movie-collection/config"
	"movie-collection/models"
	"regexp"
)

// 采集并放入listChannel中
func MakeListCollection(isFull bool)  <-chan string{
	listChannel := make(chan string)
	go func(listChannel chan<- string,isFull bool) {
		listCollection(listChannel,isFull)
	}(listChannel,isFull)
	return listChannel

}
func listCollection(listChannel chan<- string,isFull bool){
	//列表页
	cList := colly.NewCollector(
		colly.AllowURLRevisit(),
		colly.URLFilters(
			regexp.MustCompile(`http://www.okzyw.com/\?m=vod-index-pg-\d+.html`),
		),
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36 Edge/16.16299"),
		colly.ParseHTTPErrorResponse(),
	)
	isReturn := false
	cList.WithTransport(config.Transport)
	cList.OnError(func(r *colly.Response, err error) {//失败重试
		// TODO 增加失败重试次数限制
		log.Println("列表页错误URL:", r.Request.URL,err)
		close(listChannel)
	})

	cList.OnHTML("div.xing_vb>ul>li", func(e *colly.HTMLElement) { //列表
		//name := strings.Split(e.ChildText("span.xing_vb4")," ")[0]
		url := e.ChildAttr("span.xing_vb4>a", "href")
		//typeN := e.ChildText("span.xing_vb5")
		updateTime := e.ChildText("span.xing_vb7")
		updateTime2 := e.ChildText("span.xing_vb6")
		if updateTime == "" {
			updateTime = updateTime2
		}
		if url != "" {
			if isFull {
				listChannel<-"http://www.okzyw.com"+url
			}else {
				row := models.Dbm.QueryRow(" SELECT id from movie WHERE update_time = ? ",updateTime);
				id := "";
				row.Scan(&id);
				if id != ""{
					isReturn = true
				}
				listChannel<-"http://www.okzyw.com"+url
				
			}
		}

	})
	cList.OnHTML("a.pagelink_a", func(e *colly.HTMLElement) { //翻页 下一页
		name := e.Text
		url := e.Attr("href")
		if name == "下一页" {
			if isReturn{
				close(listChannel)
				return
			}
			log.Println("\n下一页:", url)
			e.Request.Visit(url)
		}
	})
	cList.OnHTML("div.pages>em", func(e *colly.HTMLElement) { //翻页 下一页
		name := e.Text
		if name == "下一页" {
			//log.Println("列表页采集结束")
			close(listChannel)
		}
	})
	cList.Visit("http://www.okzyw.com/?m=vod-index-pg-1.html")
}