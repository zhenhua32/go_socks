package socks

import (
	"math/rand"
	"reflect"
	"testing"
)

const MB = 1024 * 1024

func TestCipher(t *testing.T) {
	pwd := RandPassword()
	p, _ := parsePwd(pwd)
	cipher := newCipher(p)

	// 元数据
	org := make([]byte, pwdLength)
	for i := 0; i < pwdLength; i++ {
		org[i] = byte(i)
	}
	temp := make([]byte, pwdLength)
	copy(temp, org)

	// 加密
	cipher.encode(temp)
	// 解密
	cipher.decode(temp)

	if !reflect.DeepEqual(org, temp) {
		t.Error("解码编码数据后无法还原数据")
	}
}

func BenchmarkEncode(b *testing.B) {
	pwd := RandPassword()
	p, _ := parsePwd(pwd)
	cipher := newCipher(p)

	bs := make([]byte, MB)
	b.ResetTimer()
	rand.Read(bs) // 随机化数据
	cipher.encode(bs)
}

func BenchmarkDecode(b *testing.B) {
	pwd := RandPassword()
	p, _ := parsePwd(pwd)
	cipher := newCipher(p)

	bs := make([]byte, MB)
	b.ResetTimer()
	rand.Read(bs) // 随机化数据
	cipher.decode(bs)
}
