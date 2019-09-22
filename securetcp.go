package socks

import (
	"io"
	"log"
	"net"
)

const bufSize = 1024

// SecureTCPConn 加密传输的 TCP
type SecureTCPConn struct {
	io.ReadWriteCloser
	Cipher *cipher
}

// DecodeRead 读取并解密
func (s *SecureTCPConn) DecodeRead(bs []byte) (n int, err error) {
	n, err = s.Read(bs)
	if err != nil {
		return
	}
	s.Cipher.decode(bs[:n])
	return
}

// EncodeWrite 解密并写入
func (s *SecureTCPConn) EncodeWrite(bs []byte) (int, error) {
	s.Cipher.encode(bs)
	return s.Write(bs)
}

// EncodeCopy 读取数据, 并加密, 最后写入到 dst
func (s *SecureTCPConn) EncodeCopy(dst io.ReadWriteCloser) error {
	buf := make([]byte, bufSize)
	for {
		// 读取数据
		readCount, errRead := s.Read(buf) // 普通读
		if errRead != nil {
			if errRead == io.EOF {
				return nil
			}
			return errRead
		}

		if readCount > 0 {
			writeCount, errWrite := (&SecureTCPConn{
				ReadWriteCloser: dst,
				Cipher:          s.Cipher,
			}).EncodeWrite(buf[0:readCount]) // 加密写
			if errWrite != nil {
				return errWrite
			}
			if readCount != writeCount {
				return io.ErrShortWrite
			}
		}
	}
}

// DecodeCopy 读取加密的数据, 并解密, 最后写入到 dst
func (s *SecureTCPConn) DecodeCopy(dst io.ReadWriteCloser) error {
	buf := make([]byte, bufSize)
	for {
		readCount, errRead := s.DecodeRead(buf) // 解密读
		if errRead != nil {
			if errRead == io.EOF {
				return nil
			}
			return errRead
		}

		if readCount > 0 {
			writeCount, errWrite := dst.Write(buf[0:readCount]) // 普通写
			if errWrite != nil {
				return errWrite
			}
			if readCount != writeCount {
				return io.ErrShortWrite
			}
		}
	}
}

// DialTCPSecure 使用 net.DialTCP, 生成 SecureTCPConn
func DialTCPSecure(raddr *net.TCPAddr, cipher *cipher) (*SecureTCPConn, error) {
	// 连接远程服务器
	remoteConn, err := net.DialTCP("tcp", nil, raddr)
	if err != nil {
		return nil, err
	}

	return &SecureTCPConn{
		ReadWriteCloser: remoteConn,
		Cipher:          cipher,
	}, nil
}

// ListenSecureTCP 使用 net.ListenTCP 进行本地监听
func ListenSecureTCP(laddr *net.TCPAddr, cipher *cipher, handlerConn func(localConn *SecureTCPConn), didListen func(listenAddr net.Addr)) error {
	// 本地监听
	listener, err := net.ListenTCP("tcp", laddr)
	if err != nil {
		return err
	}

	defer listener.Close()

	if didListen != nil {
		didListen(listener.Addr())
	}

	for {
		localConn, err := listener.AcceptTCP()
		if err != nil {
			log.Println(err)
			continue
		}
		localConn.SetLinger(0)
		// 处理数据
		go handlerConn(&SecureTCPConn{
			ReadWriteCloser: localConn,
			Cipher:          cipher,
		})
	}
}
