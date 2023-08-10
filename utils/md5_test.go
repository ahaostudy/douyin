package utils

import (
	"fmt"
	"testing"
)

func TestMD5(t *testing.T) {
	pwd := MD5("douyin")
	fmt.Println(pwd)
}
