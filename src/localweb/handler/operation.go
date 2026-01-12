package handler

import (
	"fmt"
	"net"
	"net/http"
	"os"
	pathpkg "path"
	"path/filepath"
	"strings"
	"time"
)

type MyHandler struct {
	Path              string
	w                 http.ResponseWriter
	r                 *http.Request
	Subpath           string
	DetectContentType bool
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
	urlRelPath := strings.TrimPrefix(this.r.URL.Path, this.Subpath)
	urlRelPath = strings.ReplaceAll(urlRelPath, "\\", "/")
	urlRelPath = strings.TrimPrefix(urlRelPath, "/")

	rootAbs, err := filepath.Abs(this.Path)
	if err != nil {
		http.Error(this.w, "invalid root path", http.StatusInternalServerError)
		return
	}

	fsRelPath := filepath.FromSlash(urlRelPath)
	CompletePath := filepath.Join(rootAbs, fsRelPath)
	CompletePathAbs, err := filepath.Abs(CompletePath)
	if err != nil {
		http.Error(this.w, "invalid path", http.StatusBadRequest)
		return
	}

	rootCmp := rootAbs
	pathCmp := CompletePathAbs
	if os.PathSeparator == '\\' {
		rootCmp = strings.ToLower(rootCmp)
		pathCmp = strings.ToLower(pathCmp)
	}
	prefix := rootCmp
	if !strings.HasSuffix(prefix, string(os.PathSeparator)) {
		prefix += string(os.PathSeparator)
	}
	if pathCmp != rootCmp && !strings.HasPrefix(pathCmp, prefix) {
		http.Error(this.w, "invalid path", http.StatusBadRequest)
		return
	}
	// 检查文件是否存在
	var fileCheck CheckFileExists
	fileCheck = CheckFileExists{}
	// 如果文件不存在 抛出异常
	if !fileCheck.Exists(CompletePath) {
		http.NotFound(this.w, this.r)
		return
	}
	var FileInfo os.FileInfo = fileCheck.Info

	// 判断是文件或者目录 目录则继续打开返回结果 文件则直接读取
	if FileInfo.IsDir() {
		this.DisplayAllFile(CompletePath, urlRelPath)
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
		displayName := fname
		if f.IsDir() {
			displayName += "/"
		}
		href := this.Subpath + "/" + pathpkg.Join(path, fname)
		this.w.Write([]byte("<a href=\"" + href + "\">" + displayName + "</a><br>"))
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
	if this.DetectContentType {
		sniffLen := 512
		if len(content) < sniffLen {
			sniffLen = len(content)
		}
		this.w.Header().Set("Content-Type", http.DetectContentType(content[:sniffLen]))
	} else {
		this.w.Header().Set("Content-Type", "application/octet-stream")
		this.w.Header().Set("Content-Disposition", "attachment; filename="+filepath.Base(CompletePath))
	}
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
		ips := strings.Split(forwarded, ",")
		return strings.TrimSpace(ips[0])
	}
	host, _, err := net.SplitHostPort(this.r.RemoteAddr)
	if err == nil {
		return host
	}
	return this.r.RemoteAddr
}
