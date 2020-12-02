package util

import (
	"bufio"
	"fmt"
	"github.com/emirpasic/gods/sets/hashset"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

/**
 * 替换nginx配置文件内容并刷新
 * @params ips ip集合
 * @params nginxPath nginx目录（根目录）
 */
func NginxRefresh(ips *hashset.Set, nginxPath string) {
	if ips.Size() <= 0 {
		return
	}
	//转换为正确的系统文件路径
	nginxConf := filepath.FromSlash(nginxPath + "/conf/nginx.conf")
	//读取nginx配置
	configContent := ReadNginxConf(nginxConf)

	template := "        server  %v;\n"
	//替换upstream内容
	var servers string
	for _, ip := range ips.Values() {
		servers += fmt.Sprintf(template, ip.(string))
	}
	servers += "        least_conn;\n"
	configContent = fmt.Sprintf(configContent, servers)
	//重新写入到ngixn配置文件
	writeNginx(nginxConf, configContent)

	nginxReload(nginxPath)
}

/**
 * 读取nginx配置内容，并替换upstream内容
 * @params confPath 文件路径
 * @params ipSet upstream要替换的ip值
 * @return string nginx配置文件内容
 */
func ReadNginxConf(confPath string) string {
	confFile, err := os.OpenFile(confPath, os.O_RDONLY, 0666)
	defer confFile.Close()
	if err != nil {
		log.Fatalln("读取nginx配置失败", err)
	}
	rd := bufio.NewReader(confFile)
	//配置内容
	var configContent string
	//是否已进入语法所属区域
	var isArea = false
	//是否操作完成
	var succ = false
	for {
		lineStr, err := rd.ReadString('\n')
		if err != nil && err == io.EOF {
			//结尾
			configContent += lineStr
			break
		}
		if err != nil {
			log.Printf("读取nginx配置失败: %v", err)
			break
		}
		//防止对含有upstream的所有配置都被影响
		lineOther := strings.Trim(lineStr, " ")
		upstreamIdx := strings.Index(lineOther, "upstream")
		isUpstream := upstreamIdx == 0
		//非语法区域的配置才拼接
		if !isArea {
			configContent += lineStr
		}
		if !isUpstream && !isArea {
			continue
		} else if isUpstream && !isArea {
			//证明开始进入 upstream 语法区
			isArea = true
			continue
		} else if isArea && strings.Contains(lineStr, "}") {
			//离开语法所属区域
			isArea = false
			configContent += lineStr
			continue
		} else if isArea && !succ {
			//upstream语法区使用使用通配符代替，待会直接替换ip即可
			configContent += "%v\n"
			succ = true
		}
	}
	return configContent
}

//写入到nginx.conf文件
func writeNginx(nginxConf, configContent string) {
	//重新写入文件
	fileWrite, err := os.Create(nginxConf)
	if err != nil {
		log.Fatalln("创建配置失败", err)
	}
	defer fileWrite.Close()
	writer := bufio.NewWriter(fileWrite)
	_, wrErr := writer.WriteString(configContent)
	if wrErr != nil {
		log.Fatalln("写入nginx配置失败", err)
	}
	_ = fileWrite.Sync()
	_ = writer.Flush()
}

//重新载入nginx配置
func nginxReload(nginxPath string) {
	//若是win系统，执行命令reload ngxin配置
	if runtime.GOOS == "windows" {
		reload := exec.Command("cmd", "/c", "nginx.exe -s reload")
		//切换调用nginx目录（默认是当前文件的目录执行命令，会报错路径不正确）
		reload.Dir = filepath.FromSlash(nginxPath)
		result, err := reload.Output()
		if err != nil {
			panic(err)
		}
		fmt.Println(string(result))
		log.Println("win 刷新完成！")
	}
}
