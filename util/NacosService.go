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
	"io/ioutil"
	"log"
)

type NacosConf struct {
	//nacos注册中心地址
	IpAddr string
	//端口
	Port uint64
	//地址前缀
	ContextPath string
	//namespace
	NamespaceId string
	//监听的服务名称
	ListenerService string
	Group           string
}

//nginx刷新配置
type NginxRefreshConf struct {
	NacosConf NacosConf
	//nginx 安装根目录
	NginxPath string
}

//nacos客户端配置
func InitNacosConfig(nacosCof *NacosConf) naming_client.INamingClient {
	logPath, _ := ioutil.TempDir("", "nacos-log")
	cachePath, _ := ioutil.TempDir("", "nacos-cache")
	//client客户端请求配置
	clientConfig := constant.ClientConfig{
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              logPath,
		CacheDir:            cachePath,
		RotateTime:          "1h",
		MaxAge:              3,
		LogLevel:            "debug",
		NamespaceId:         nacosCof.NamespaceId,
	}

	// nacos注册中心服务地址配置（可多个）
	serverConfigs := []constant.ServerConfig{
		{
			IpAddr:      nacosCof.IpAddr,
			ContextPath: nacosCof.ContextPath,
			Port:        nacosCof.Port,
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
 * @params NginxRefreshConf
 */
func RefershNginxListener(namingClient naming_client.INamingClient, conf *NginxRefreshConf) {
	//监听参数配置
	var subParam = vo.SubscribeParam{
		ServiceName: conf.NacosConf.ListenerService,
		GroupName:   conf.NacosConf.Group, // 默认值DEFAULT_GROUP
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
			NginxRefresh(ips, conf.NginxPath)
		},
	}

	subErr := namingClient.Subscribe(&subParam)
	if subErr != nil {
		log.Fatalf("Subscribe nacos service error ：%v\n", subErr)
	}
}
