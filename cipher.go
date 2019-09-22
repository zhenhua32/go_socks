package socks

type cipher struct {
	encodePwd *password
	decodePwd *password
}

// 加密
func (c *cipher) encode(bs []byte) {
	for i, v := range bs {
		bs[i] = c.encodePwd[v]
	}
}

// 解密
func (c *cipher) decode(bs []byte) {
	for i, v := range bs {
		bs[i] = c.decodePwd[v]
	}
}

// 创建一个编码器, 包含加解密对
func newCipher(encodePwd *password) *cipher {
	decodePwd := &password{}
	// 对称
	for i, v := range encodePwd {
		encodePwd[i] = v
		decodePwd[v] = byte(i)
	}
	return &cipher{
		encodePwd: encodePwd,
		decodePwd: decodePwd,
	}
}
