package main

import (
	"fmt"
	"localweb/handler"
	"net/http"
)

func main() {
	// 先获取传参
	var s handler.Start
	s = handler.Start{}
	s.DetermainStartupParametere()

	// 创建事件处理器 同时传入监听地址
	var m *handler.MyHandler
	m = &handler.MyHandler{Path: *s.Path, Subpath: "/file"}
	// 启动服务器
	fmt.Printf("Server listening on port %s...\n", *s.Port)
	http.Handle("/", m)
	server := http.Server{
		Addr:                         *s.Port,
		Handler:                      nil,
		DisableGeneralOptionsHandler: false,
		TLSConfig:                    nil,
		ReadTimeout:                  0,
		ReadHeaderTimeout:            0,
		WriteTimeout:                 0,
		IdleTimeout:                  0,
		MaxHeaderBytes:               0,
		TLSNextProto:                 nil,
		ConnState:                    nil,
		ErrorLog:                     nil,
		BaseContext:                  nil,
		ConnContext:                  nil,
	}
	server.ListenAndServe()
}
