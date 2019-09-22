package socks

import "net"

// Local 保存本地连接信息
type Local struct {
	Cipher     *cipher
	ListenAddr *net.TCPAddr
	RemoteAddr *net.TCPAddr
}

// NewLocal 创建一个新的本地客户端
func NewLocal(pwd string, listenAddr, remoteAddr string) (*Local, error) {
	bsPwd, err := parsePwd(pwd)
	if err != nil {
		return nil, err
	}
	structListenAddr, err := net.ResolveTCPAddr("tcp", listenAddr)
	if err != nil {
		return nil, err
	}
	structRemoteAddr, err := net.ResolveTCPAddr("tcp", remoteAddr)
	if err != nil {
		return nil, err
	}
	return &Local{
		Cipher:     newCipher(bsPwd),
		ListenAddr: structListenAddr,
		RemoteAddr: structRemoteAddr,
	}, nil
}

// Listen 启动本地监听
func (local *Local) Listen(didListen func(listenAddr net.Addr)) error {
	return ListenSecureTCP(local.ListenAddr, local.Cipher, local.handleConn, didListen)
}

// handleConn 本地数据处理
func (local *Local) handleConn(userConn *SecureTCPConn) {
	// userConn 是本地的
	defer userConn.Close()

	// 远程服务器
	proxyServer, err := DialTCPSecure(local.RemoteAddr, local.Cipher)
	if err != nil {
		return
	}
	defer proxyServer.Close()

	// 从 proxyServer 获取加密后的数据  server => local
	go func() {
		err := proxyServer.DecodeCopy(userConn)
		// 在 copy 的过程中可能会存在网络超时等 error 被 return，只要有一个发生了错误就退出本次工作
		if err != nil {
			userConn.Close()
			proxyServer.Close()
		}
	}()
	// 将本地的数据发送到 proxyServer  local => server
	userConn.EncodeCopy(proxyServer)
}
