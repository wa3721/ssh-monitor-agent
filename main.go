package main

import (
	httpserver "kubeauth/pkg/http"
	scriptinit "kubeauth/pkg/script-init"
)

//1.初始化文件到/etc/profile.d 中
//2.设置一个API接收ssh命令行及ip和用户执行目录等信息
//3.将接收到的这些信息存储在一个k8s crd中 取名 sshaudit(crd结构需要设计)
//4.设计另一个crd，负责管理这些sshaudit记录，比如多长时间清理一次，或者多少数据清理一次，排除哪些ip

func main() {
	scriptinit.NewChecklist().RunAll()
	scriptinit.Exec()
	httpserver.NewServer().StartServer()
}
