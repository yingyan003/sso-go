sso单点登录系统v1版本API
=====================

### 大纲
--------

v1版本支持如下API操作：

#### 登录
---------

    * API: POST /v1/login

#### 退出
---------

    * API: GET /v1/Logout

#### 登录认证
-----------

    * API: GET /v1/auth

#### 用户注册
-----------

    * API: POST /v1/user

#### 修改用户信息
---------------

    * API: PUT /v1/user

#### 获取用户列表
---------------

    * API: GET /v1/user/list

#### 通过uid查询用户
---------------

    * API: GET /v1/user/:uid


### 详细设计
-----------

#### 登录
---------

* 请求URL：POST /v1/login

* 请求头：
    * "token: default"

> 备注：
> 为解决重复登录的问题，约定请求头必须有"token"。token值如下几种情况：
> 1.当值为"default"时，表示首次登录（即用户目前不处于登录状态）
> 2.当值为有效token时（即验证用户已登录的token）,表示用户重复登录。此时将直接返回前端，不允许重复登录。
> 3.若token解析成功但已失效，则放行登录。
> 3.当值不属于如上2者时，如header没有"token",或者token的值为空，或是其他值，将返回错误。

* 请求体：

    * username, password (必填)

*example*
---------
```json
    {
    "username":"p",
    "password":"2"
    }
```

* 返回结果：

{
    code: "ok"表示请求成功，其他表示失败 (如下所有api适用此规则)
    message: 请求状态信息
    data: token
}

*example*
---------
```json
{"code":"OK","message":"登录成功","data":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6ImV5SmhiR2NpT2lKRlV6STFOaUlzSW5SNWNDSTZJa3BYVkNKOS5leUp5WldScGMwdGxlU0k2SW1VelpXSTFObUUzTFdGbU5tUXROR1EyTUMwNE0yUmpMVFV4WkRCaFpUY3hNbUZqTlNKOS5aYkVHQkM0eXFBR2lGbVRKNzQxOWRPZ3RXRXRGbGdkdFprMXNDXzFFZXgxTnRsbkFXWW80YUNaWnA0ekd6bUNtMjlvNlMtNkNlTHlhMEx2TjVYTGF1USIsInVzZXJJZCI6MjYsInVzZXJuYW1lIjoicCJ9.0A__kI23VMpTU1q_ATpeY1WnosxAHPLxBu40jmEw3MU"}
```

#### 退出
---------

* API: GET /v1/logout

* 请求头：
    * "token: 登录成功返回的有效token值"

* 返回结果：

*example*
---------
```json
{"code":"OK","message":"用户退出成功","data":null}
```

#### 登录认证
---------

* API: GET /v1/auth

* 请求头：
    * "token: 登录成功返回的有效token值"

* 返回结果：

    * data: 用户ID

*example*
---------
```json
{"code":"OK","message":"认证成功，用户已登录","data":"26"}
```

#### 用户注册
---------

* API: POST /v1/user

* 请求体：
   * username, password (必填)

   *example*
   ---------
   ```json
       {
       "username":"p",
       "password":"2"
       }
   ```

* 返回结果：

*example*
---------
```json
{"code":"OK","message":"用户创建成功","data":null}
```

#### 修改用户信息
---------

* API: PUT /v1/user

* 请求头：
    * "token: 登录成功返回的有效token值"

* 请求体：
   * username, password (必填)

   *example*
   ---------
   ```json
       {
       "username":"p",
       "password":"2"
       }
   ```

* 返回结果：

*example*
---------
```json
{"code":"OK","message":"用户更新成功","data":null}
```

#### 获取用户列表
---------

* API: GET /v1/user/list

* 返回结果：

*example*
---------
```json
{
"code":"OK",
"message":"查询用户列表成功",
"data":[
    {"id":4,"username":"a","password":"","createtime":"0001-01-01T00:00:00Z","updatetime":"0001-01-01T00:00:00Z"}
    {"id":6,"username":"b","password":"","createtime":"0001-01-01T00:00:00Z","updatetime":"0001-01-01T00:00:00Z"}
    ]
}
```

#### 通过uid查询用户
------------------

* API: GET /v1/user/:uid

* 返回结果：

*example*
---------
```json
{
{"code":"OK","message":"用户存在","data":{"id":35,"username":"r","password":"TzM+MhHAx9/hJ0MLOzhbzSFv5svqnAp//ZPectaAVLw=","createtime":"2018-01-26T15:33:04+08:00","updatetime":"2018-01-26T15:34:03+08:00"}}
}
```