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
	//todo 这3个key不能瞎改，待研究
	//ES256 keys
	ECDSAKeyD = "A181A5E47749A1891E98DC6429F694DA1389E317155B8B45F1762FC2A0860F3C"
	ECDSAKeyX = "F4A9122D5BCF4725F59EEB36D6A4C9E2E3F9C728B03BA8C8654F61ECE102A020"
	ECDSAKeyY = "E7E9D5033A87258CB698E6160E9E01A815675B2916215F516A0B64BA691806AB"

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
