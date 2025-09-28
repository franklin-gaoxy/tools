package server

import (
	"github.com/gin-gonic/gin"
	"k8s.io/klog"
	"net/http"
)

/*
This is the function that implements all interfaces of gin
*/

func initEnvironment(str string) {
	klog.Info(str)
}

func Test(c *gin.Context) {
	klog.Info(c.Request.RequestURI)
}

func Help(c *gin.Context) {
	klog.Info(c.Request.RequestURI)
	c.String(http.StatusOK, "<h1>help</h1>"+
		"<h5>/test:</h5>test interface."+
		"<h5>/version:</h5>print version."+
		"<h5>/all:</h5>response all info."+
		"<h5>/cpu:</h5>response cpu info."+
		"<h5>/memory:</h5>response memory info."+
		"<h5>/disk:</h5>response host disk info."+
		"<h5>/network:</h5>response host network info."+
		"<h5>/node:</h5>response host node and os info.")
}

func GetAll(c *gin.Context) {
	klog.Info(c.Request.RequestURI)
	c.JSON(http.StatusOK, gin.H{
		"node info":    GetInfoNode(),
		"cpu info":     GetInfoCPU(),
		"mem info":     GetInfoMemory(),
		"disk info":    GetInfoDisk(),
		"network info": GetInfoNetwork(),
	})
}

func GetCPU(c *gin.Context) {
	klog.Info(c.Request.RequestURI)
	c.JSON(http.StatusOK, GetInfoCPU())
	return
}

func GetMemory(c *gin.Context) {
	klog.Info(c.Request.RequestURI)
	c.JSON(http.StatusOK, GetInfoMemory())
	return
}

func GetDisk(c *gin.Context) {
	klog.Info(c.Request.RequestURI)
	c.JSON(http.StatusOK, GetInfoDisk())
	return
}

func GetNetwork(c *gin.Context) {
	klog.Info(c.Request.RequestURI)
	c.JSON(http.StatusOK, GetInfoNetwork())
	return
}

func GetNode(c *gin.Context) {
	klog.Info(c.Request.RequestURI)
	c.JSON(http.StatusOK, GetInfoNode())
	return
}
