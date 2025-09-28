package server

import (
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"k8s.io/klog"
	"os"
	"strings"
)

var port int

func InitStart() bool {
	klogInit()
	return true
}

func klogInit() {
	klog.InitFlags(nil)
	// flag.Set("V", "2")
	// flag.Parse()

	parameterProcessing()
	_ = flag.CommandLine.Parse(nil)
	//klog.Infof("klog init: log event %d\n", tools.LogEvent)
	defer klog.Flush()
}

func version() string {
	return "v1.0"
}

/*
parameterProcessing
Used to handle conflicts between klog framework and cobra framework
Enable the klog framework to correctly receive the parameters of -- v
*/
func parameterProcessing() {
	// 临时存储 os.Args
	args := os.Args[1:]
	remainingArgs := []string{os.Args[0]}

	for _, arg := range args {
		if strings.HasPrefix(arg, "--v=") {
			vValue := strings.TrimPrefix(arg, "--v=")
			fmt.Printf("Handling --v=%s parameter\n", vValue)

			// Force setting the - v parameter of the klog framework
			if err := flag.Set("v", vValue); err != nil {
				fmt.Printf("Failed to set klog -v flag: %v\n", err)
			}
		} else {
			remainingArgs = append(remainingArgs, arg)
		}
	}

	// 重新设置 os.Args 为剩余参数，不包含 --v 参数
	os.Args = remainingArgs
}

/*
cobra相关内容
*/

// configFilePath 此变量用于接受--config参数的内容 然后传递到启动函数里
var configFilePath string

// rootCmd 主命令 也就是不加任何子命令情况 执行此函数
var rootCmd = cobra.Command{
	Use:   "config",
	Short: "input config file address.",
	Run: func(cmd *cobra.Command, args []string) {
		NewStart(configFilePath)

	},
}

// 增加一个新的子命令 version
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "print version.",
	Run: func(cmd *cobra.Command, args []string) {
		klog.Infoln(version())
	},
}

// init cobra框架 将所有的都添加到rootCmd这个主命令下
func init() {
	rootCmd.AddCommand(versionCmd)
	// 添加指定端口
	rootCmd.Flags().IntVarP(&port, "port", "p", 8080, "port to listen on.")
}

// Start Cobra's startup function
func Start() {

	if err := rootCmd.Execute(); err != nil {
		klog.Fatalln("start error! please check databases config!")
	}
}

/*
New Start
*/

// NewStart default command execute
func NewStart(configFilePath string) {
	// start gin server: use command -p
	startGinServer(port)

}

/*
startGinServer
Add the new interface address that needs to be processed here.
The corresponding method is implemented in the server.go file.
*/
func startGinServer(port int) {
	var route *gin.Engine
	route = gin.Default()

	// TODO demo: binding interface
	route.GET("/version", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "normal",
			"version": version(),
		})
	})

	// TODO demo: Test interface
	route.POST("/test", func(c *gin.Context) {
		Test(c)
	})

	route.GET("/help", func(c *gin.Context) {
		Help(c)
	})

	route.GET("/all", func(c *gin.Context) {
		GetAll(c)
	})

	route.GET("/cpu", func(c *gin.Context) {
		GetCPU(c)
	})

	route.GET("/memory", func(c *gin.Context) {
		GetMemory(c)
	})

	route.GET("/disk", func(c *gin.Context) {
		GetDisk(c)
	})

	route.GET("/network", func(c *gin.Context) {
		GetNetwork(c)
	})

	route.GET("/node", func(c *gin.Context) {
		GetNode(c)
	})

	klog.V(1).Infof("start gin server on port %d", port)
	_ = route.Run(fmt.Sprintf(":%d", port))
}
