package jwt

import (
	"crypto/sha1"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"io"
	"rxt/internal/config"
)

// 哈希加密
// laravel解密兼容示例
func Sha1Encode(data string) string {
	if data == "" {
		return ""
	}
	t := sha1.New()
	io.WriteString(t, data)
	return fmt.Sprintf("%x", t.Sum(nil))
}

// 生成签名
func secret() jwt.Keyfunc {
	jwtCnf := config.GetJwtCnf()
	config.New().Get(&jwtCnf, "jwt")
	return func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		hmacSampleSecret := []byte(jwtCnf.Secret)
		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return hmacSampleSecret, nil
	}
}
