package main

import (
	"fmt"
	"log"
	"net"

	"github.com/phayes/freeport"
	"tzh.com/socks"
	"tzh.com/socks/cmd"
)

const version = "0.1"

func main() {
	log.SetFlags(log.Lshortfile)

	// 服务端监听端口随机生成
	port, err := freeport.GetFreePort()
	if err != nil {
		port = 7448
	}

	// 默认配置
	config := &cmd.Config{
		ListenAddr: fmt.Sprintf(":%d", port),
		// 密码随机生成
		Password: socks.RandPassword(),
	}
	config.ReadConfig()
	config.SaveConfig()

	// 启动 server 端并监听
	lsServer, err := socks.NewServer(config.Password, config.ListenAddr)
	if err != nil {
		log.Fatalln(err)
	}
	log.Fatalln(lsServer.Listen(
		func(listenAddr net.Addr) {
			log.Println("使用配置:")
			log.Println(fmt.Sprintf(`
本地监听地址 listen：
%s
密码 password：
%s
	`, listenAddr, config.Password))
			log.Printf("lightsocks-server:%s 启动成功 监听在 %s\n", version, listenAddr.String())
		},
	))
}
