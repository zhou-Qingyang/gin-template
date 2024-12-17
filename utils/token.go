package utils

import (
	"errors"
	"time"

	jwt "github.com/golang-jwt/jwt/v4"
)

type JWT struct {
	SigningKey []byte
}

type NumericDate int64

var (
	TokenExpired     = errors.New("Token is expired")            //token过期
	TokenNotValidYet = errors.New("Token not active yet")        //token无效
	TokenMalformed   = errors.New("That's not even a token")     //根本不是token
	TokenInvalid     = errors.New("Couldn't handle this token:") //不能处理这个token
)

type BaseClaims struct {
	StudentPrimaryId int64  `json:"studentPrimaryId"`
	StudentId        string `json:"studentId"`
	StudentName      string `json:"studentName"`
	UserId           int64  `json:"userId"`
	IsAdmin          int8   `json:"isAdmin"`
}
type CustomClaims struct {
	BaseClaims
	BufferTime int64 `json:"bufferTime"`
	jwt.RegisteredClaims
}

func NewJWT() *JWT {
	return &JWT{
		[]byte("1234567890"),
	}
}

func (j *JWT) CreateClaims(baseClaims BaseClaims) CustomClaims {
	claims := CustomClaims{
		BaseClaims: baseClaims,
		BufferTime: int64(24 * time.Hour / time.Second), // 缓冲时间1天 缓冲时间内会获得新的token刷新令牌 此时一个用户会存在两个有效令牌 但是前端只留一个 另一个会丢失
		RegisteredClaims: jwt.RegisteredClaims{
			NotBefore: jwt.NewNumericDate(time.Now().Add(-1000)),              // 签名生效时间
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)), // 过期时间 7天配置文件
			Issuer:    "tz-gin",                                               // 签名的发行者
		},
	}
	return claims
}

func (j *JWT) CreateToken(claims CustomClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.SigningKey)
}

// 解析token
func (j *JWT) ParseToken(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (i interface{}, e error) {
		return j.SigningKey, nil
	})
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, TokenMalformed
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return nil, TokenExpired
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, TokenNotValidYet
			} else {
				return nil, TokenInvalid
			}
		}
	}
	if token != nil {
		if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
			return claims, nil
		}
		return nil, TokenInvalid
	} else {
		return nil, TokenInvalid
	}
}
