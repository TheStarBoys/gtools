package crypto

import (
	"testing"
	"fmt"
)

func TestGeneratePwd(t *testing.T) {
	salt := "123456"
	pwd := GeneratePwd("抖音", salt, 12)
	fmt.Println(pwd)
}
