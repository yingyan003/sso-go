package jwtUtil

import (
	"testing"
	"crypto/ecdsa"
	"crypto/elliptic"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"crypto/rand"
)
var(
	prk *ecdsa.PrivateKey
	puk ecdsa.PublicKey
	curve    elliptic.Curve
)


func TestEcdsaTest(t *testing.T){
	randKey:=rand.Reader
	var err error
	prk, err = ecdsa.GenerateKey(elliptic.P256(),randKey )
	if err!=nil{
		fmt.Println("generate key error",err)
	}
	puk=prk.PublicKey
	fmt.Println("prk",prk," \npbk",puk)
//get
	claims := jwt.MapClaims{
		"redisKey": "mykey",
	}
	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	ss, err := token.SignedString(prk)
	if err!=nil{
		fmt.Println("generate token error",err)
	}
	fmt.Println("tokenString",ss)
//parse
	token, err = jwt.Parse(ss, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodECDSA); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return &puk, nil
	})
	if err!=nil{
		fmt.Println("parse token error",err)
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		fmt.Println("claims:",claims)
	}
}
