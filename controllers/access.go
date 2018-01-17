package controllers

import (
	_ "github.com/astaxie/beego/session/redis"
	"sso/models"
	"github.com/dgrijalva/jwt-go"
	"github.com/astaxie/beego/logs"
	"sso/common/errors"
	"fmt"
	"sso/common/util/encrypt"
)

const (
	//token过期时间
	DEFAULT_EXPIRATION = 120 // 24*60*60s=86400
	//jwt签名秘钥
	SIGNED_KEY = "ssoSignedKey"
	KEY       = "\n-----BEGIN AAAAB3NzaC1yc2EAAAADAQABAAABAQDEyquqaClVFyPo0RHmE/Eron9/UX9YdECFTl+HqWotPXas9yJLjhNtA+0IIUw+LDa9QIo7NqXx5SXasuD6S+D3lTS79TvAIir2rDg4ZVMZOVD7yI1r3zIXdyJhvxUwXoMXpdCLRqpUT5kokOsD/X+E+phAE15oN77uKVHb8hJOf7nfjpJ1+a5PXiWCSZ9ZxjjfTsT1NuEQVnWxX/m2y/xs7Ssi0cInfYyZyYSmrYk6bmy+0z9Envnv/agskQjAlB4vm9WzVgf3pCcEFUgORzWgv68ufJWO624teRqFsIbfndAQ7epxyzTX8r8dL5S4mczaDQIwY3wSIj8EQA0lM4jp\n-----END-----"
	KEY2       = "-----BEGIN CERTIFICATE-----\nMIIDHDCCAgSgAwIBAgIIbGH5q1SgoHMwDQYJKoZIhvcNAQEFBQAwMTEvMC0GA1UE\nAxMmc2VjdXJldG9rZW4uc3lzdGVtLmdzZXJ2aWNlYWNjb3VudC5jb20wHhcNMTgw\nMTE0MDA0NTI2WhcNMTgwMTE3MDExNTI2WjAxMS8wLQYDVQQDEyZzZWN1cmV0b2tl\nbi5zeXN0ZW0uZ3NlcnZpY2VhY2NvdW50LmNvbTCCASIwDQYJKoZIhvcNAQEBBQAD\nggEPADCCAQoCggEBAJn7CSKGaEC/f5KGPpmOvcXrnR9O5JEHThGWNhkzpAWkZlrS\nKMc3NBpoPMocz6VHq/yVQ06mZOUXfXCKB4mdLp0CUG1OFkfFrNlKCWPlPst7YuvE\n2s0LKGZm8rRsL0ai5y9SXoyI2DCeh7Mzb4vVyjTjMTahm1w6BgrrQlekGU04dLrE\nkLYf7r14X4sNHrP1UqEkI0OSXG1dnkthHv+aH75Cs9WY5HqmS0ciBC7aOdDx0lj4\nrHSMw0PkXPFmixW3JbIbIzHBDGYAzR3vcEIYRhWYSIF1QLJkUxFicL14/aMXftM4\nw8i9AgqyNh4dWDfNsTwKeZP6KVnrI1Ph1aUYdy8CAwEAAaM4MDYwDAYDVR0TAQH/\nBAIwADAOBgNVHQ8BAf8EBAMCB4AwFgYDVR0lAQH/BAwwCgYIKwYBBQUHAwIwDQYJ\nKoZIhvcNAQEFBQADggEBAE5/xchxalj6kQNvqhmr4y8xqr2QzMAPzFiCZ58xYV/Z\nTy/9TjkvogD8XpVL0tUelycSKG0C5q3beKYRbLH2s3JHpF53XI3e64WM1ieWnhau\nCMqJzfSqRuC72Dq/H0f5pIYdcc8xge0Vt+fP0UX8yFHDsVh43IKPO29GPdgTv+2l\nYDZuHxI/aH3bqvtwMFcaNhYrfp7PqdPi67KWXkcXOTFZGZRx7MgDgFQ/X+cMwIZI\nA+ZsSMeEqY11VbyyeK+Xul/hfjfVRSPFDJzirE9LdbD6gt2GhgIAZuJE/Mvg3y6N\nNyEGBfQuruRTRFAjPBypNJ2E8L7wFvTZf9BcBpjGpn4=\n-----END CERTIFICATE-----\n"
)

type AccessController struct {
	BaseController
}

func (a *AccessController) Auth(){
	status := errors.NewStatus(errors.OK, "认证成功，用户已登录")
	a.Ctx.Output.Body(status.ToBytes())
}

func (a *AccessController) Login() {
	user := new(models.User)
	status := a.getReqBody(user)
	if status != nil {
		logs.Error("登录失败，username=%s, msg=%s", user.Username, status.Message)
		a.Ctx.Output.Body(status.ToBytes())
		return
	}

	user.Password = encrypt.NewEncryption(user.Password).String()
	//检查用户名密码是否正确
	user,status=models.QueryUserByNameAndPwd(user.Username,user.Password)
	if status != nil {
		logs.Error("登录失败，用户名或密码错误")
		a.Ctx.Output.Body(status.ToBytes())
		return
	}

	sessionId := a.createSession(user.Id, user.Username)
	//token, status := createToken(sessionId, user.Username)
	//if status != nil {
	//	logs.Error("登录失败，创建token失败")
	//	a.Ctx.Output.Body(status.ToBytes())
	//	return
	//}

	logs.Info("登录成功，username=%s, sessionId=%s",user.Username,sessionId)
	status=errors.NewStatus(errors.OK,"登录成功")
	status.Data=sessionId
	a.Ctx.Output.Body(status.ToBytes())
}

//todo
func (a *AccessController) Logout() {
	fmt.Println("enter logout")

	//tokenString := a.Ctx.Input.Header("token")
	//token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
	//	return SIGNED_KEY, nil
	//})
	//if err != nil {
	//	logs.Error("退出登录时，token解析失败")
	//	status := errors.NewStatus(errors.TOKEN_ERR, "退出登录时，token解析失败")
	//	a.Ctx.Output.Body(status.ToBytes())
	//	return
	//}
	//claims, ok := token.Claims.(jwt.MapClaims)
	//if !ok {
	//	logs.Error("退出登录时，获取jwt claims失败")
	//	status := errors.NewStatus(errors.TOKEN_ERR, "退出登录时，jwt claimts获取失败")
	//	a.Ctx.Output.Body(status.ToBytes())
	//	return
	//}
	//sessionId := claims["sessionId"].(string)
	//username := claims["username"].(string)

	sessionId := a.Ctx.Input.Header("sessionId")
	a.deleteSession(sessionId)

	//logs.Info("用户退出成功，username=%s", username)
	logs.Info("用户退出成功，sessionId=%s", sessionId)
	status := errors.NewStatus(errors.OK, "用户退出成功")
	a.Ctx.Output.Body(status.ToBytes())
}

func (a *AccessController) createSession(uid int64, username string) string {
	s := models.NewSession(uid, username)
	a.SetSession(s.Id, s)
	return s.Id
}

func (a *AccessController) deleteSession(sessionId string) {
	s:=a.GetSession(sessionId)
	if s!=nil{
		fmt.Println("deltession: session=%v",s)
		a.DelSession(sessionId)
	}else{
		fmt.Println("a.DelSession(sessionId): session not exist")
	}
}

func createToken(sessionId, username string) (string, *errors.Status) {
	/*certPEM := KEY
	//certPEM = strings.Replace(certPEM, "\\n", "\n", -1)
	//certPEM = strings.Replace(certPEM, "\"", "", -1)
	block, _ := pem.Decode([]byte(certPEM))
	fmt.Println("block", block)

	cert, _ := x509.ParseCertificate(block.Bytes)
	fmt.Println("cert", cert)

	rsaPublicKey := cert.PublicKey.(*rsa.PublicKey)
	fmt.Println("rsaPublicKey", rsaPublicKey)
	//pubKey,err:=jwt.ParseRSAPublicKeyFromPEM([]byte(SIGNED_KEY))
	//*/

	claims := jwt.MapClaims{
		////过期时间，单位：秒
		//"exp":int64(time.Now().Unix() + DEFAULT_EXPIRATION),
		////签发时间
		//"iat":int64(time.Now().Unix()),
		//"UserId":uid,
		"sessionId": sessionId,
		"username":  username,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	//todo 现在不支持[]byte()字节片段，key的类型必须是rsa.PublicKey
	//ss, err := token.SignedString(rsaPublicKey)
	//todo 不用key临时测试
	ss, err := token.SigningString()
	fmt.Println("ss", ss)
	s2, err := token.SignedString(ss)
	fmt.Println("s2",s2)
	token,err =jwt.Parse(s2,func(token *jwt.Token)(interface{},error){
		return ss,nil
	})
	if err != nil {
		logs.Error("token生成签名错误,err=%v", err)
		status := errors.NewStatus(errors.TOKEN_ERR, "token生成签名错误")
		return "", status
	}
	fmt.Println("token", token.Claims)
	return ss, nil
}
