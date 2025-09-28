package handler

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type MyHandler struct {
	Path    string
	w       http.ResponseWriter
	r       *http.Request
	Subpath string
}

func (this *MyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// 主函数
	this.w = w
	this.r = r
	clientIP := this.GetClientIP()
	this.PrintLog("received request from " + clientIP)
	this.Actuator()
	this.PrintLog("client IP " + clientIP + " corresponding completion.")
}
func (this MyHandler) Actuator() {
	// 获取访问的目录 或者 文件 的路径
	path := strings.TrimPrefix(this.r.URL.Path, this.Subpath)
	// 路径拼接
	CompletePath := filepath.Join(this.Path, path)
	// 检查文件是否存在
	var fileCheck CheckFileExists
	fileCheck = CheckFileExists{}
	// 如果文件不存在 抛出异常
	if !fileCheck.Exists(CompletePath) {
		panic("File or directory not found!")
		return
	}
	var FileInfo os.FileInfo = fileCheck.Info

	// 判断是文件或者目录 目录则继续打开返回结果 文件则直接读取
	if FileInfo.IsDir() {
		this.DisplayAllFile(CompletePath, path)
	} else {
		this.ReadFile(CompletePath)
	}

}

func (this MyHandler) DisplayAllFile(CompletePath string, path string) {
	/*
		展示目录中的所有文件
		this.w.Write 写入返回给客户端的数据
	*/
	this.PrintLog("open directory " + CompletePath)
	files, _ := os.ReadDir(CompletePath)
	this.w.Write([]byte("<html><body>"))
	for _, f := range files {
		fname := f.Name()
		if f.IsDir() {
			fname += "/"
		}
		this.w.Write([]byte("<a href=\"" + this.Subpath + "/" + filepath.Join(path, fname) + "\">" + fname + "</a><br>"))
	}
	this.w.Write([]byte("</body></html>"))
}

func (this MyHandler) ReadFile(CompletePath string) {
	// 读取文件并返回给浏览器
	this.PrintLog("download file " + CompletePath)
	content, err := os.ReadFile(CompletePath)
	// 遇到异常返回错误
	if err != nil {
		http.Error(this.w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 设置响应头及返回数据
	this.w.Header().Set("Content-Type", "application/octet-stream")
	this.w.Header().Set("Content-Disposition", "attachment; filename="+filepath.Base(CompletePath))
	this.w.Write(content)
}

func (this MyHandler) PrintLog(log string) {
	now := time.Now()
	hour := now.Hour()
	minute := now.Minute()
	second := now.Second()
	fmt.Printf("%02d:%02d:%02d :", hour, minute, second)
	fmt.Printf("%s\n", log)
}

func (this MyHandler) GetClientIP() string {
	forwarded := this.r.Header.Get("X-Forwarded-For")
	if forwarded != "" {
		ips := strings.Split(forwarded, ", ")
		return ips[0]
	}
	return strings.Split(this.r.RemoteAddr, ":")[0]
}
