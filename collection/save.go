package collection

import (
	"github.com/smartwalle/dbs"
)

var movies []*Movie

func save(movie *Movie) {
	if len(movies)<10{
		movies=append(movies,movie)
	}else {
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
	}
	//fmt.Printf("消费：%s ",movie.Name)
}
