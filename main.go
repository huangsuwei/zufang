package main

import (
	"IhomeWeb/handler"
	_ "IhomeWeb/models"
	"github.com/julienschmidt/httprouter"
	"github.com/micro/go-micro/util/log"
	"github.com/micro/go-micro/web"
	"github.com/micro/go-plugins/registry/consul"
	"net/http"
)

func main() {
	//reg := consul.NewRegistry(func(op *registry.Options) {
	//	op.Addrs = []string {
	//		"127.0.0.1:8500",
	//	}
	//})
	reg := consul.NewRegistry()

	// create new web service
	service := web.NewService(
		web.Name("go.micro.web.IhomeWeb"),
		web.Version("latest"),
		web.Address(":8081"),
		web.Registry(reg),
	)

	// initialise service
	if err := service.Init(); err != nil {
		log.Fatal(err)
	}

	rou := httprouter.New()
	rou.NotFound = http.FileServer(http.Dir("html"))
	// 获取地区请求
	rou.GET("/api/v1.0/areas", handler.GetArea)
	// 获取session
	rou.GET("/api/v1.0/session", handler.GetSession)
	// 获取index页面
	rou.GET("/api/v1.0/house/index", handler.GetIndex)
	// 获取图片验证码
	rou.GET("/api/v1.0/imagecode/:uuid", handler.GetImageCd)
	// 获取手机验证码
	rou.GET("/api/v1.0/smscode/:mobile", handler.GetSmsCd)
	// 注册
	rou.POST("/api/v1.0/users", handler.PostRet)
	// 登录
	rou.POST("/api/v1.0/sessions", handler.PostLogin)
	// 退出登录
	rou.DELETE("/api/v1.0/session", handler.DeleteSession)
	// 获取用户信息
	rou.GET("/api/v1.0/user", handler.GetUserInfo)
	// 上传用户头像
	rou.POST("/api/v1.0/user/avatar", handler.PostAvatar)
	// 更改用户名
	rou.PUT("/api/v1.0/user/name", handler.PutUserInfo)
	// 查看用户实名认证
	rou.GET("/api/v1.0/user/auth", handler.GetUserAuth)
	// 提交实名认证
	rou.POST("/api/v1.0/user/auth", handler.PostUserAuth)
	// 获取我的已发布房源
	rou.GET("/api/v1.0/user/houses", handler.GetUserHouses)

	service.Handle("/", rou)

	// register html handler
	//service.Handle("/", http.FileServer(http.Dir("html")))

	// register call handler
	//service.HandleFunc("/IhomeWeb/call", handler.IhomeWebCall)

	// run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
