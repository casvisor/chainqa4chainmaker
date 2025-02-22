package main

import (
	"fmt"
	"os"
	"tencent-chainmaker/routers"
	"tencent-chainmaker/setting"
)

const defaultConfFile = "./conf/config.ini"

func main() {
	confFile := defaultConfFile
	// 判断是否有指定配置文件
	if len(os.Args) > 2 {
		fmt.Println("[Tencent ChainMaker]使用配置文件： ", os.Args[1])
		confFile = os.Args[1]
	} else {
		fmt.Println("[Tencent ChainMaker]使用默认配置文件： ", defaultConfFile)
	}
	// 加载配置文件
	if err := setting.Init(confFile); err != nil {
		fmt.Printf("[Tencent ChainMaker]加载配置文件失败，错误信息:%v\n", err)
		return
	}

	// 注册路由
	r := routers.SetupRouter()
	if err := r.Run(fmt.Sprintf(":%d", setting.Conf.Port)); err != nil {
		fmt.Printf("[Tencent ChainMaker]启动服务失败，错误信息:%v\n", err)
	}
}
