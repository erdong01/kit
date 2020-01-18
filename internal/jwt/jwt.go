package jwt

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"rxt/internal/config"
	"strings"
	"time"
)

type I interface {
	GenerateToken(sub int64) (string, error)
}

type Claims struct {
	jwt.StandardClaims
}

// 生产token
func GenerateToken(user User) (string, error) {
	jwtCnf := config.GetJwtCnf()
	config.New().Get(&jwtCnf, "jwt")
	var nowTime time.Time
	nowTime = time.Now()
	a := time.Duration(jwtCnf.Expire) * time.Hour
	expireTime := nowTime.Add(a)
	claims := jwt.MapClaims{
		"exp": expireTime.Unix(),
		"iss": jwtCnf.Issuer,
		"sub": user.Sub,
		"iat": nowTime.Unix(),
		"nbf": nowTime.Unix(),
		"prv": Sha1Encode(user.Prv),
	}
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString([]byte(jwtCnf.Secret))
	return token, err
}

// 验证token
func ValidateToken(signedToken string) (claims jwt.MapClaims, err error) {
	token, err := jwt.Parse(signedToken, secret())
	if err != nil {
		err = errors.New("cannot convert claim to mapclaim")
		return
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		err = errors.New("cannot convert claim to mapclaim")
		return
	}
	//验证token，如果token被修改过则为false
	if !token.Valid {
		err = errors.New("token is invalid")
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
