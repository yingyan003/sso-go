SSO单点登录系统
=============

# 使用方式：

1. 在本地新建一个目录dir(名字随意)
mkdir dir

2. dir目录下新建目录src(必须是src)
cd dir
mkdir src

3. 在src目录下新建git仓库(前提是你的电脑安装了git)
cd src
git init(新建一个git仓库)

4. 将项目代码拉到本地
git pull git@github.com:yingyan003/sso-go.git master

5. 将dir目录的路径设置为gopath  
export GOPATH=dir目录的全路径  

6. 运行beego项目
cd sso
bee run main.go 

说明：
· 如果你的开发工具是GoLand,运行最好不要使用GoLand的run按钮运行，容易报错，导致编译失败。  
· 一定要在命令行下，cd到main.go目录下，用“bee run main.go”运行程序。


# 部署方式

###  虚拟机部署

1. 编译程序，生成二进制文件main
cd main.go所在目录
go build main.go

（注意交叉编译，如在mac上直接go build生成的二进制只能在mac上执行，不能在linux上执行）

2. 部署虚拟机
· ssh登录到目标虚拟机
· cd 到某个目录，如sso目录
· 将本地编译好的main二进制文件scp到sso目录下
· 将本地项目中的conf文件scp到sso目录下
· 运行程序：./main

成功运行时，服务就对外了。可以通过“虚拟机IP：服务暴露的端口”方式访问。


### 容器云部署（部署到k8s）

> 前言：  
sso项目目录下有2个文件，".gitlab-ci.yml"和"Dockerfile"，这2个文件就是部署容器云需要的。  
因为公司的代码仓库用的是gitlab，而gitlab编译程序需要设置一个runner,其实就是个机器，虚机和物理机都行，物理机太浪费，一般都是虚机。因为go程序在这里编译，所以需要安装go等需要的环境。  

操作：   
当从本地（你自己的电脑）把代码通过git提交到gitlab仓库时，gitlab会自动执行“.gitlab-ci.yml”文件，该文件中的命令其实就是shell命令。把你要做的操作写在.gitlab-ci.yml中，gitlab会依次执行里面指定的操作命令。  
  
我们在.gitlab-ci.yml中进行了以下操作：  
· 如果项目目录不存在时新建，并在项目目录下新建src目录，再把整个sso拷贝到src目录下。
· 设置GOPATH
· 编译项目
· 使用Dockerfile打镜像
· 给镜像打标签tag
· 将镜像推送到远程仓库

Dockerfile说明：
打镜像时需要用到Dockerfile(文件名不可修改)，在.gitlab-ci.yml文件中，  
命令“docker build -t $IMAGE:$TAG . ”就是打镜像的操作，“.”表示使用当前路径下的Dockerfile来打镜像。打镜像的过程就是按照Dockerfile中指定的操作进行.  
在Dockerfile中，我们定义了以下操作：    
· ADD conf /conf ：将当前路径下的conf文件夹拷贝到容器的/路径下，文件夹名为conf  
· ADD Dockerfile / : 将当前路径下的Dockerfile文件拷贝到容器的/路径下，文件夹名为Dockerfile    
· ADD main /sso ：将当前路径下的二进制文件main拷贝到容器的/路径下，文件夹名为sso  
· ENTRYPOINT ["/sso"] ： 程序运行入口是sso,也就是当容器启动时，指定“./sso”

gitlib构建日志：

![](https://github.com/yingyan003/sso-go/blob/master/picture/log1.png)

![](https://github.com/yingyan003/sso-go/blob/master/picture/log2.png)

![](https://github.com/yingyan003/sso-go/blob/master/picture/log3.png)

![](https://github.com/yingyan003/sso-go/blob/master/picture/log4.png)