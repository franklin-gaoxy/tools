package handler

import (
	"flag"
	"os"
	"runtime"
)

type Start struct {
	// 端口和地址
	Port              *string
	Path              *string
	DetectContentType *bool
}

func (this *Start) DetermainStartupParametere() {
	// 获取输入的变量
	this.Port = flag.String("port", ":8080", "please input a port,such as format :8080")
	this.Path = flag.String("path", "nil", "please input a file path. if there is no input,C drive will be used by default under windows system.")
	this.DetectContentType = new(bool)
	flag.BoolVar(this.DetectContentType, "d", false, "enable auto detection of Content-Type (full name: -detect-content-type); default false uses application/octet-stream with attachment")
	flag.BoolVar(this.DetectContentType, "detect-content-type", false, "enable auto detection of Content-Type (short: -d); default false uses application/octet-stream with attachment")
	flag.Parse()
	// 检查指定的path是否存在
	this.CheckOperationSystem()
	this.CheckPath()
}
func (this *Start) CheckOperationSystem() {
	if *this.Path == "nil" {
		// 如果没有传入path 根据系统自动判断默认值
		if runtime.GOOS == "windows" {
			*this.Path = "D:\\"
		} else if runtime.GOOS == "linux" || runtime.GOOS == "darwin" {
			*this.Path = "./"
		} else {
			panic("unable to determine the current operating system type! or pass in the parameter path.")
		}
	}
}
func (this *Start) CheckPath() {
	// 指定了路径判断路径是否存在
	if *this.Path != "nil" {
		var c CheckFileExists
		c = CheckFileExists{}
		if !c.Exists(*this.Path) {
			panic("the specified path does not exists!")
		}
	}
}

type CheckFileExists struct {
	Info os.FileInfo // 文件元数据信息
}

func (this *CheckFileExists) Exists(path string) bool {
	// 判断给出的文件或者目录是否存在
	var err error
	this.Info, err = os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return true
}
