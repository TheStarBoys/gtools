package crypto

import (
	"crypto/sha256"
	"golang.org/x/crypto/ripemd160"
	"fmt"
)

// 生成密码, n 截前n个字符
func GeneratePwd(rawPwd, salt string, n int) string {
	pwd := generateRawPwd([]byte(rawPwd), []byte(salt))

	return rawPwd2Str(pwd, n)
}

func generateRawPwd(rawPwd, salt []byte) []byte {
	rawPwd = append(rawPwd, salt...)
	firstSHA := sha256.Sum256(rawPwd)
	secondSHA := sha256.Sum256(firstSHA[:])

	RIPEMD160 := ripemd160.New()
	_, err := RIPEMD160.Write(secondSHA[:])
	if err != nil {
		panic(err)
	}

	pwdRIPEMD := RIPEMD160.Sum(nil)

	return pwdRIPEMD
}

//func injectChar(pwd []byte) []byte {
//
//}

func rawPwd2Str(rawPwd []byte, n int) string {
	str := fmt.Sprintf("%x", rawPwd)
	if n >= len(str) || n <= 0 {
		return str
	}

	return str[:n]
}

