package models

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"movie-collection/config"
)

var Dbm *sql.DB

func init(){
	var err error
	//Dbm,err = sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/db_movie")
	Dbm,err = sql.Open("mysql", config.DataSourceName)
	if err != nil{
		panic(err)
	}
}