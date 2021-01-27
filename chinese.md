# GolangWebFramework
GolangWebFramework 是一个非常简单的Go语言编写的网页开发框架. 借鉴了Springboot, Django等流行框架的设计思想.
## 快速开始
执行命令: 

go get -u github.com/gaoze1998/GolangWebFramework

GolangWebFramework create api testproject

cd testproject

go run main.go

browse http://localhost:8082/gwft
## 支持
1. RESTful理念
2. 对象关系映射框架(ORM)
3. 命令行项目脚手架
4. 内存会话管理
5. 加密/解密算法库
6. 分布式处理
## 更多细节
Check testproject(created above) for details.

Check Controller for HTTP handler example.

Check Config for config file example. Config file contains configurations of listening address, port and etc. 

Check Model for application model which is used to serialize and deserialize.

Check Router for routes information.
