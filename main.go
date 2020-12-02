package main

import (
	"nginx-nacos-example/util"
	"os"
	"os/signal"
	"syscall"
)

//ngxin 根路径
const NGINX_PATH = "D:\\util\\java\\nginx-1.18.0"

func main() {
	go func() {
		client := util.NacosConfig()
		util.RefershNginxListener(client, NGINX_PATH)
	}()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGKILL)
	<-sig
}
