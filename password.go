package socks

import (
	"encoding/base64"
	"errors"
	"math/rand"
	"strings"
	"time"
)

const pwdLength = 256

type password [pwdLength]byte

func init() {
	rand.Seed(time.Now().Unix())
}

func (p *password) String() string {
	return base64.StdEncoding.EncodeToString(p[:])
}

// 解密 base64 后的字符串, 还原密码
func parsePwd(text string) (*password, error) {
	bs, err := base64.StdEncoding.DecodeString(strings.TrimSpace(text))
	if err != nil || len(bs) != pwdLength {
		return nil, errors.New("不合法的密码")
	}
	pwd := password{}
	copy(pwd[:], bs)
	bs = nil
	return &pwd, nil
}

// RandPassword 生成一个随机的密码
func RandPassword() string {
	intArr := rand.Perm(pwdLength)
	pwd := &password{}
	for i, v := range intArr {
		pwd[i] = byte(v)
		if i == v {
			return RandPassword()
		}
	}
	return pwd.String()
}
