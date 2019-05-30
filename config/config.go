package config

import (
	"net"
	"net/http"
	"time"
)

var SaveNum int = 1000
var ThreadNum int = 100
var DataSourceName string = "root:8084810821@tcp(mysql:3306)/db_movie"
var Transport http.RoundTripper = &http.Transport{
	Dial: func(netw, addr string) (net.Conn, error) {
		conn, err := net.DialTimeout(netw, addr, time.Second*30) //设置建立连接超时
		if err != nil {
			return nil, err
		}
		conn.SetDeadline(time.Now().Add(time.Second * 30)) //设置发送接受数据超时
		return conn, nil
	},
	ResponseHeaderTimeout: time.Second * 30,
}
