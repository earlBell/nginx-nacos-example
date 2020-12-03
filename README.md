## Nginx + Nacos 动态服务发现例子

### 1.场景 
> *  用于单体应用+Nginx+Nacos整合、微服务网关与Nginx整合，动态刷新Nacos的健康服务到Nginx中

### 2.使用方式

#### 2.1 准备
> * nginx配置文件`nginx.conf` 先加入一个`upstream`(ip随意配置一个服务都行，第一次监听nacos就会替换掉)，项目已经提供一个简单的配置文件例子
> * 启动Nginx
> * 启动nacos，将[服务注册到nacos上][1] 

#### 2.2 修改main.go的配置



```go
   nacosCof := util.NacosConf{
        //nacos地址
		IpAddr:          "192.168.50.75", 
		Port:            8848,
	    //nacos路径
		ContextPath:     "/nacos",
		//nacos的NamespaceId
		NamespaceId:     "3146d3eb-2422-4439-a063-a9a0df197c5e",
		//监听注册到nacos的服务名称
		ListenerService: "demo",
		//nacos的分组名称
		Group:           "GZ",
	}
	nginxRefreshCof := util.NginxRefreshConf{
	    #安装nginx的根目录
		NginxPath: "/usr/local/nginx",
		NacosConf: nacosCof,
	}
``` 

#### 2.3 运行
> 运行`main.go` 或者 build到相应的平台运行





### 3.说明
> 1：目前只针对Nginx单个upstream的情况处理 
>
> 2：目前只针对Nacos指定的服务进行刷新（特定的服务）
>
> 3：极端频繁刷新Nginx配置的情况，还没使用缓存方式减少IO操作
>
> 4：仅支持windows和linux
>
> 5：当所有服务都down掉，不会再刷新Nginx配置

#### 


  [1]: https://github.com/earlBell/nacos-example