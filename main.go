package main

import (
	"github.com/smartwalle/dbs"
	"log"
	"movie-collection/collection"
	"movie-collection/config"
	_ "movie-collection/config"
	"movie-collection/models"
	"time"
)

func main() {
	for{
		now := time.Now()
		hour, _, _:= now.Clock()
		if hour==1{
			run(true)
			time.Sleep(time.Hour*6)
		}else {
			run(false)
			time.Sleep(time.Minute*10)
		}
	}
}

func run(isFull bool)  {
	listCollection := collection.MakeListCollection(isFull)
	detailCollection := collection.MakeDetailCollection(listCollection)
	movies := make([]*collection.Movie,0)
	for  movie:= range detailCollection  {
		if len(movies)<config.SaveNum{
			movies = append(movies,movie)
		}else {
			save(movies)
			movies = make([]*collection.Movie,0)
		}
	}
	if len(movies)>0{
		save(movies)
	}
	log.Println("采集完成")
}

func save(movies []*collection.Movie)  {
	var ib = dbs.NewInsertBuilder()
	ib.Table("movie")
	ib.Columns("id","name","alias","update_status","nav1",
		"nav2","img","director","actor","type","area","language",
		"time","update_time","introduce","content")
	for _,m := range movies{
		ib.Values(m.Id,m.Name, m.Alias, m.UpdateStatus,
			m.Nav1, m.Nav2, m.Img, m.Director,
			m.Actor, m.Type, m.Area, m.Language,
			m.Time, m.UpdateTime, m.Introduce, m.Content)
	}
	ib.Suffix(dbs.OnDuplicateKeyUpdate().
		Append("name=VALUES(name)").
		Append("alias=VALUES(alias)").
		Append("update_status=VALUES(update_status)").
		Append("nav1=VALUES(nav1)").
		Append("nav2=VALUES(nav2)").
		Append("img=VALUES(img)").
		Append("director=VALUES(director)").
		Append("actor=VALUES(actor)").
		Append("type=VALUES(type)").
		Append("area=VALUES(area)").
		Append("language=VALUES(language)").
		Append("time=VALUES(time)").
		Append("update_time=VALUES(update_time)").
		Append("introduce=VALUES(introduce)").
		Append("content=VALUES(content)"))
	_, e := ib.Exec(models.Dbm)
	if e !=nil{
		log.Printf("数据插入错误",e)
	}

}
