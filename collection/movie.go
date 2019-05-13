package collection


type Movie struct {
	Id,Name,Alias, UpdateStatus,Num,Nav1, Nav2,Img,   Director, Actor, Type, Area, Language, Time, UpdateTime ,Introduce,Content,Enable,CreateTime string
}


func (m *Movie) InitNav1(){
	switch m.Nav2{
	case "动作片","喜剧片","爱情片","科幻片","恐怖片", "剧情片","战争片","纪录片","微电影":
		m.Nav1="电影片"
	case "国产剧","香港剧","欧美剧","日本剧","韩国剧","台湾剧","海外剧":
		m.Nav1="连续剧"
	case "综艺片":
		m.Nav1="综艺片"
	case "动漫片":
		m.Nav1="动漫片"
	case "福利片","伦理片":
		m.Nav1="福利片"
	}
}