package test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"runtime"
	"path/filepath"
	_ "sso/routers"

	"github.com/astaxie/beego"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/dgrijalva/jwt-go"
	"fmt"
)

func init() {
	_, file, _, _ := runtime.Caller(1)
	apppath, _ := filepath.Abs(filepath.Dir(filepath.Join(file, ".." + string(filepath.Separator))))
	beego.TestBeegoInit(apppath)
}

// TestGet is a sample to run an endpoint test
func TestGet(t *testing.T) {
	r, _ := http.NewRequest("GET", "/v1/object", nil)
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	beego.Trace("testing", "TestGet", "Code[%d]\n%s", w.Code, w.Body.String())

	Convey("Subject: Test Station Endpoint\n", t, func() {
	        Convey("Status Code Should Be 200", func() {
	                So(w.Code, ShouldEqual, 200)
	        })
	        Convey("The Result Should Not Be Empty", func() {
	                So(w.Body.Len(), ShouldBeGreaterThan, 0)
	        })
	})
}

func TestToken(t *testing.T){
	claims := jwt.MapClaims{
		////过期时间，单位：秒
		//"exp":int64(time.Now().Unix() + DEFAULT_EXPIRATION),
		////签发时间
		//"iat":int64(time.Now().Unix()),
		//"UserId":uid,
		"sessionId": "111",
		"username":  "zxy",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	//todo 现在不支持[]byte()字节片段，key的类型必须是rsa.PublicKey
	//ss, err := token.SignedString([]byte("key"))
	//todo 不用key临时测试
	//ss, err := token.SigningString()
	//fmt.Println("ss", ss)
	s2, err := token.SignedString([]byte("key"))
	if err!=nil{
		fmt.Println("SignedString error,err=",err)
	}
	fmt.Print("token",s2)
	token,err =jwt.Parse(s2,func(token *jwt.Token)(interface{},error){
		return []byte("key"),nil
	})
	if err != nil {
		fmt.Errorf("parsr error,err=%v",err)
	}
	fmt.Println("token", token.Claims)
}

type SContrller struct{
	beego.Controller
}

func TestSession(t *testing.T){
	s:=SContrller{}
	s:=s.Controller.GetSession("39aa0fa6-8e28-43e0-9e06-4f67d0db63b2")
	fmt.Println("s,",s)
}
