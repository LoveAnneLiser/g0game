package main

import (
	"common/config"
	"common/metrics"
	"flag"
	"fmt"
)

var configFile = flag.String("config", "application.yml", "config file")

func main() {
	// 1. 加载配置
	flag.Parse()
	config.InitConfig(*configFile)
	fmt.Println(config.Conf)
	// 2. 启动监控
	go func() {
		err := metrics.Serve(fmt.Sprintf("0.0.0.0:%d", config.Conf.MetricPort))
		if err != nil {
			panic(err)
		}
	}()
	// 3. 启动grpc服务
	select {}
}
