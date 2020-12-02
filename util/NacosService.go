package util

import (
	"encoding/json"
	"fmt"
	"github.com/emirpasic/gods/sets/hashset"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/model"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"log"
)

//nacos客户端配置
func NacosConfig() naming_client.INamingClient {
	//client客户端请求配置
	clientConfig := constant.ClientConfig{
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              "D:\\tmp\\log",
		CacheDir:            "D:\\tmp\\cache",
		RotateTime:          "1h",
		MaxAge:              3,
		LogLevel:            "debug",
		NamespaceId:         "3146d3eb-2422-4439-a063-a9a0df197c5e",
	}

	// nacos注册中心服务地址配置（可多个）
	serverConfigs := []constant.ServerConfig{
		{
			IpAddr:      "localhost",
			ContextPath: "/nacos",
			Port:        8848,
			Scheme:      "http",
		},
	}

	// 创建服务发现客户端
	client, err := clients.CreateNamingClient(map[string]interface{}{
		"serverConfigs": serverConfigs,
		"clientConfig":  clientConfig,
	})
	if err != nil {
		log.Fatalf("init client error ：%v\n", err)
	}
	return client
}

/**
 * 监听nacos的服务并把最新的服务ip配置到nginx
 * @params namingClient nacos客户端
 * @params nginxPath nginx目录（根目录）
 */
func RefershNginxListener(namingClient naming_client.INamingClient, nginxPath string) {
	//监听参数配置
	var subParam = vo.SubscribeParam{
		ServiceName: "demo",
		GroupName:   "GZ", // 默认值DEFAULT_GROUP
		SubscribeCallback: func(services []model.SubscribeService, err error) {
			ips := hashset.New()
			for _, service := range services {
				jsonByte, _ := json.Marshal(service)
				log.Printf("service : %v\n", string(jsonByte))
				if !service.Enable {
					continue
				}
				ip := fmt.Sprintf("%v:%d", service.Ip, service.Port)
				ips.Add(ip)
			}
			//刷新nginx配置
			NginxRefresh(ips, nginxPath)
		},
	}

	subErr := namingClient.Subscribe(&subParam)
	if subErr != nil {
		log.Fatalf("Subscribe nacos service error ：%v\n", subErr)
	}
}
