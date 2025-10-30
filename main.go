package main

import (
	"context"
	"flag"
	"sshmonitor/config"
	httpserver "sshmonitor/pkg/http"
	scriptinit "sshmonitor/pkg/script-init"
)

//1.初始化文件到/etc/profile.d 中
//2.设置一个API接收ssh命令行及ip和用户执行目录等信息
//3.将接收到的这些信息存储在一个k8s crd中 取名 sshaudit(crd结构需要设计)
//4.设计另一个crd，负责管理这些sshaudit记录，比如多长时间清理一次，或者多少数据清理一次，排除哪些ip

func main() {
	var loglevel string
	var logOutput string
	var prod bool
	var catcherLength int
	var controllerAddr string
	var consumerNum int
	flag.StringVar(&loglevel, "loglevel", "debug", "Set log level , Optional: debug, info, warn, error, default: info")
	flag.StringVar(&logOutput, "logoutput", "stdout", "Set log output path ,Optional: console, file, double, default: console")
	flag.BoolVar(&prod, "prod", false, "Set deployment mode or prod mode Optional: false, true default: false")
	flag.IntVar(&catcherLength, "catcherlength", 1000, "Set ssh command provider channel length, default: 1000")
	flag.StringVar(&controllerAddr, "controlleraddr", "127.0.0.1:8081", "Set controller address to link kubernetes , default: 0.0.0.0")
	flag.IntVar(&consumerNum, "consumernum", 3, "Set consumer number , default: 3")
	flag.Parse()

	err := config.InitLogger(loglevel, logOutput, prod)
	if err != nil {
		panic(err)
	}
	scriptinit.NewChecklist().RunAll()
	scriptinit.Exec()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	for i := 0; i < consumerNum; i++ {
		go httpserver.NewConsumer().Consume(ctx, controllerAddr)
	}
	httpserver.NewServer().StartServer(catcherLength)
}
