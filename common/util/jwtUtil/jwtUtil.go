package jwtUtil

import (
	"math/big"
	"github.com/dgrijalva/jwt-go"
	"crypto/ecdsa"
	"crypto/elliptic"
	"github.com/astaxie/beego/logs"
	"fmt"
	"sso/common/errors"
	"sso/models"
)

const (
	//ES256 keys
	ECDSAKeyD = "3999161F60FCAE8D34E05D55F7C07ED9C761CCC102EA75C42C1C5483A6FEBCB4"
	ECDSAKeyX = "C2124438109F89DB29063DAF3C66CDBC31A9D4E2292652997127663DB219429E"
	ECDSAKeyY = "3868441F8D52D78EDD5396182C72B58E71F96BE758B28C1775C4813FAC86929E"

	//HS256 signed key
	SIGNED_KEY = "HSSignedKey"
)

//获取签名算法为ES256的token
//该token的内容只有Redis的key,用于保存用户的登录状态
func GetEStoken(redisKey string) (string, *errors.Status) {
	keyD := new(big.Int)
	keyX := new(big.Int)
	keyY := new(big.Int)

	keyD.SetString(ECDSAKeyD, 16)
	keyX.SetString(ECDSAKeyX, 16)
	keyY.SetString(ECDSAKeyY, 16)

	claims := jwt.MapClaims{
		"redisKey": redisKey,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	publicKey := ecdsa.PublicKey{
		Curve: elliptic.P256(),
		X:     keyX,
		Y:     keyY,
	}

	privateKey := ecdsa.PrivateKey{D: keyD, PublicKey: publicKey}

	ss, err := token.SignedString(&privateKey)
	if err != nil {
		logs.Error("ES256的token生成签名错误,err=%v", err)
		status := errors.NewStatus(errors.TOKEN_ERR, "EStoken生成签名错误")
		return "", status
	}
	return ss, nil
}

//获取签名算法为HS256的token
func GetHStoken(tokenFirst string, user *models.User) (string, *errors.Status) {
	claims := jwt.MapClaims{
		"id": tokenFirst,
		//解析时，该变量的类型被转换成float64
		"userId":   user.Id,
		"username": user.Username,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString([]byte(SIGNED_KEY))
	if err != nil {
		logs.Error("token生成签名错误,err=%v", err)
		status := errors.NewStatus(errors.TOKEN_ERR, "token生成签名错误")
		return "", status
	}
	return ss, nil
}

//解析签名算法为ES256的token
func ParseEStoken(tokenES string) (string, *errors.Status) {
	keyX := new(big.Int)
	keyY := new(big.Int)

	keyX.SetString(ECDSAKeyX, 16)
	keyY.SetString(ECDSAKeyY, 16)

	publicKey := ecdsa.PublicKey{
		Curve: elliptic.P256(),
		X:     keyX,
		Y:     keyY,
	}

	token, err := jwt.Parse(tokenES, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodECDSA); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return &publicKey, nil
	})
	if err != nil {
		logs.Error("ES256的token解析错误,err=%v", err)
		status := errors.NewStatus(errors.TOKEN_ERR, "token解析失败")
		return "", status
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		return claims["redisKey"].(string), nil
	}

	logs.Error("ParseEStoken:Claims类型转换失败")
	status := errors.NewStatus(errors.TOKEN_ERR, "token类型转换错误")
	return "", status
}

//解析签名算法为HS256的token
func ParseHStoken(tokenString string) (jwt.MapClaims, *errors.Status) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(SIGNED_KEY), nil
	})
	if err != nil {
		logs.Error("HS256的token解析错误，err:", err)
		status := errors.NewStatus(errors.TOKEN_ERR, "token解析失败")
		return nil, status
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		logs.Error("ParseHStoken:claims类型转换失败")
		status := errors.NewStatus(errors.TOKEN_ERR, "token类型转换错误")
		return nil, status
	}
	return claims, nil
}
