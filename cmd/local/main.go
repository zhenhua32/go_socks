package main

import (
	"fmt"
	"log"
	"net"

	"tzh.com/socks"
	"tzh.com/socks/cmd"
)

// DefaultListenAddr 本地默认端口
const DefaultListenAddr = ":7448"
const version = "0.1"

func main() {
	log.SetFlags(log.Lshortfile)

	// 默认配置
	config := &cmd.Config{
		ListenAddr: DefaultListenAddr,
	}
	config.ReadConfig()
	config.SaveConfig()

	// 启动 local 服务
	local, err := socks.NewLocal(
		config.Password,
		config.ListenAddr,
		config.RemoteAddr,
	)
	if err != nil {
		log.Fatalln(err)
	}
	log.Fatalln(local.Listen(
		func(listenAddr net.Addr) {
			log.Println("使用配置:")
			log.Println(
				fmt.Sprintf(`
本地监听地址 listen：
%s
远程服务地址 remote：
%s
密码 password：
%s
`, listenAddr, config.RemoteAddr, config.Password))
			log.Printf("local 启动成功: %s 监听在 %s \n", version, listenAddr.String())
		},
	))
}
