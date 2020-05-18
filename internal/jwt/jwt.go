package jwt

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/erDong01/micro-kit/internal/config"
	"strconv"
	"strings"
	"time"
)

type I interface {
	GenerateToken(sub int64) (string, error)
}

type Claims struct {
	jwt.StandardClaims
	Prv string `json:"prv,omitempty"`
}

// 生产token
func GenerateToken(user User) (string, error) {
	jwtCnf := config.GetJwtCnf()
	config.New().Get(&jwtCnf, "jwt")
	var nowTime time.Time
	nowTime = time.Now()
	a := time.Duration(jwtCnf.Expire) * time.Hour
	expireTime := nowTime.Add(a)
	SubStr := strconv.FormatInt(user.Sub, 10)

	claims := Claims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			Issuer:    jwtCnf.Issuer,
			Subject:   SubStr,
			IssuedAt:  nowTime.Unix(),
			NotBefore: nowTime.Unix(),
		},
		Prv: Sha1Encode(user.Prv),
	}
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString([]byte(jwtCnf.Secret))
	return token, err
}

// 验证token
func ValidateToken(signedToken string) (claims Claims, err error) {
	jwtCnf := config.GetJwtCnf()
	jwt.ParseWithClaims(signedToken, &claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jwtCnf.Secret), nil
	})
	token, err := jwt.Parse(signedToken, secret())
	if err != nil {
		err = errors.New("cannot convert claim to mapclaim")
		return
	}
	//验证token，如果token被修改过则为false
	if !token.Valid {
		err = errors.New("token is invalid")
		return
	}
	fmt.Println(claims.ExpiresAt, time.Now().Unix())
	if claims.ExpiresAt < time.Now().Unix() {
		err = errors.New("token已过期")
		return
	}
	return claims, nil
}

// 获取前端token
func GetToken(c *gin.Context) (token string, err error) {
	bearerLength := len("Bearer ")
	hToken := c.GetHeader("Authorization")
	if len(hToken) < bearerLength {
		return token, errors.New("header Authorization has not Bearer token")
	}
	token = strings.TrimSpace(hToken[bearerLength:])
	return token, nil
}
