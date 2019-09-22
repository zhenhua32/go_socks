package socks

import (
	"sort"
	"testing"
)

func (p *password) Len() int {
	return pwdLength
}

func (p *password) Less(i, j int) bool {
	return p[i] < p[j]
}

func (p *password) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func TestRandPassword(t *testing.T) {
	pwd := RandPassword()
	t.Logf("%v", pwd)
	bsPwd, err := parsePwd(pwd)
	t.Logf("%v", *bsPwd)
	if err != nil {
		t.Error(err)
	}
	sort.Sort(bsPwd)
	for i := 0; i < pwdLength; i++ {
		if bsPwd[i] != byte(i) {
			t.Error("不能出现重复的位数, 应该由 0-255 组成, 且只包含一次")
		}
	}
}
