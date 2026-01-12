package main

import (
	"fmt"
	"localweb/handler"
	"net/http"
	"path/filepath"
)

func main() {
	// 先获取传参
	var s handler.Start
	s = handler.Start{}
	s.DetermainStartupParametere()

	// 创建事件处理器 同时传入监听地址
	var m *handler.MyHandler
	m = &handler.MyHandler{Path: *s.Path, Subpath: "/file", DetectContentType: *s.DetectContentType}

	absPath, err := filepath.Abs(*s.Path)
	if err != nil {
		absPath = *s.Path
	}
	fmt.Printf("Startup parameters:\n")
	fmt.Printf("- port: %s\n", *s.Port)
	fmt.Printf("- path: %s\n", *s.Path)
	fmt.Printf("- path(abs): %s\n", absPath)
	fmt.Printf("- route prefix: %s\n", m.Subpath)
	if m.DetectContentType {
		fmt.Printf("Enabled features:\n")
		fmt.Printf("- auto Content-Type detection (short: -d, full: -detect-content-type)\n")
	} else {
		fmt.Printf("Enabled features:\n")
		fmt.Printf("- force download (Content-Disposition: attachment, Content-Type: application/octet-stream)\n")
	}
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
