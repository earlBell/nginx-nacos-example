package main

import (
	"nginx-nacos-example/util"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	nacosCof := util.NacosConf{
		IpAddr:          "192.168.50.75",
		Port:            8848,
		ContextPath:     "/nacos",
		NamespaceId:     "3146d3eb-2422-4439-a063-a9a0df197c5e",
		ListenerService: "demo",
		Group:           "GZ",
	}
	nginxRefreshCof := util.NginxRefreshConf{
		NginxPath: "/usr/local/nginx",
		NacosConf: nacosCof,
	}

	go func() {
		client := util.InitNacosConfig(&nacosCof)
		util.RefershNginxListener(client, &nginxRefreshCof)
	}()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGKILL)
	<-sig
}
