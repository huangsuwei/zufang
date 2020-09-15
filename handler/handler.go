package handler

import (
	"IhomeWeb/models"
	DeleteSessionProto "IhomeWeb/proto/DeleteSession"
	GetAreaProto "IhomeWeb/proto/GetArea"
	GetImageCdProto "IhomeWeb/proto/GetImageCd"
	GetSessionProto "IhomeWeb/proto/GetSession"
	GetSmscdProto "IhomeWeb/proto/GetSmscd"
	GetUserHousesProto "IhomeWeb/proto/GetUserHouses"
	GetUserInfoProto "IhomeWeb/proto/GetUserInfo"
	PostAvatarProto "IhomeWeb/proto/PostAvatar"
	PostLoginProto "IhomeWeb/proto/PostLogin"
	PostRetProto "IhomeWeb/proto/PostRet"
	PostUserAuthProto "IhomeWeb/proto/PostUserAuth"
	PutUserInfoProto "IhomeWeb/proto/PutUserInfo"
	"IhomeWeb/utils"
	"context"
	"encoding/json"
	"fmt"
	"github.com/afocus/captcha"
	"github.com/astaxie/beego"
	"github.com/julienschmidt/httprouter"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/service/grpc"
	"github.com/micro/go-plugins/registry/consul"
	"image"
	"image/png"
	"net/http"
	"reflect"
	"regexp"
	//"time"
)

/*func IhomeWebCall(w http.ResponseWriter, r *http.Request) {
	// decode the incoming request as json
	var request map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// call the backend service
	GetAreaSrvClient := GetAreaSrv.NewGetAreaService("go.micro.srv.IhomeWeb", client.DefaultClient)
	rsp, err := GetAreaSrvClient.Call(context.TODO(), &GetAreaSrv.Request{
		Name: request["name"].(string),
	})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// we want to augment the response
	response := map[string]interface{}{
		"msg": rsp.Msg,
		"ref": time.Now().UnixNano(),
	}

	// encode and write the response as json
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}*/

// 获取地区信息
func GetArea(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// 创建grpc服务
	reg := consul.NewRegistry()

	server := grpc.NewService(
		//micro.Name("go.micro.srv.GetArea"),
		micro.Version("latest"),
		//micro.Address(":8081"),
		micro.Registry(reg),
	)
	/*server := grpc.NewService()*/
	server.Init()

	// 请求服务
	getAreaClient := GetAreaProto.NewGetAreaService("go.micro.srv.GetArea", server.Client())
	//callOpts := func(opts *client.CallOptions) {
	//	opts.RequestTimeout = time.Second * 30
	//	opts.DialTimeout = time.Second * 30
	//}
	rsp, err := getAreaClient.Call(context.TODO(), &GetAreaProto.Request{} /*callOpts*/)
	if err != nil {
		http.Error(w, err.Error(), 500)
		fmt.Println("------------------错粗了--------------", err.Error())
		return
	}

	// 接收数据
	var areaList []models.Area
	// 循环接收数据
	for _, value := range rsp.Data {
		tmp := models.Area{Id: int(value.Aid), Name: value.Aname}
		print(tmp.Name)
		areaList = append(areaList, tmp)
	}
	fmt.Println("------------------------成功获得返回数据---------------------------------")

	// 返回给前端的map
	response := map[string]interface{}{
		"errno":  rsp.Error,
		"errmsg": rsp.Errmsg,
		"data":   areaList,
	}

	//会传数据的时候三直接发送过去的并没有设置数据格式
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

// 获取session
func GetSession(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Println("------------------------GetSession---------------------------------")

	// 读取cookie
	cookie, err := r.Cookie("userlogin")
	if err != nil {
		// 没有获取到cookie
		// 返回给前端的map
		response := map[string]interface{}{
			"errno":  utils.RECODE_SESSIONERR,
			"errmsg": utils.RecodeText(utils.RECODE_SESSIONERR),
			"data":   "",
		}
		//会传数据的时候三直接发送过去的并没有设置数据格式
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		return
	}

	// 创建服务
	reg := consul.NewRegistry()
	server := grpc.NewService(
		micro.Name("go.micro.srv.GetSession"),
		micro.Version("latest"),
		micro.Registry(reg),
	)
	server.Init()

	getSessionService := GetSessionProto.NewGetSessionService("go.micro.srv.GetSession", server.Client())
	rsp, err := getSessionService.Call(context.TODO(), &GetSessionProto.Request{
		Sessionid: cookie.Value,
	})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	//var data map[string]string
	data := make(map[string]string)
	data["name"] = rsp.UserName

	response := map[string]interface{}{
		"errno":  rsp.Errno,
		"errmsg": rsp.Errmsg,
		"data":   data,
	}
	// 设置返回数据格式
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

}

// 获取首页轮播图
func GetIndex(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Println("------------------------GetIndex---------------------------------")
	// 返回给前端的map
	response := map[string]interface{}{
		"errno":  utils.RECODE_OK,
		"errmsg": utils.RecodeText(utils.RECODE_OK),
		"data":   "",
	}

	//会传数据的时候三直接发送过去的并没有设置数据格式
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

// 获取图形验证码
func GetImageCd(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	fmt.Println("------------------------获取图形验证码 GetImageCd---------------------------------")
	// 接收web传的参数
	// 创建服务
	reg := consul.NewRegistry()
	server := grpc.NewService(
		micro.Registry(reg),
	)
	// 服务初始化
	server.Init()

	// 连接服务
	GetImageCdService := GetImageCdProto.NewGetImageCdService("go.micro.srv.GetImageCd", server.Client())

	// 接收参数
	uuid := p.ByName("uuid")
	rsp, err := GetImageCdService.Call(context.TODO(), &GetImageCdProto.Request{
		Uuid: uuid,
	})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	// 接收图片信息的图片格式
	var img image.RGBA
	img.Pix = rsp.Pix
	img.Stride = int(rsp.Stride)
	img.Rect.Min.X = int(rsp.Min.X)
	img.Rect.Min.Y = int(rsp.Min.Y)
	img.Rect.Max.X = int(rsp.Max.X)
	img.Rect.Max.Y = int(rsp.Max.Y)

	var image captcha.Image
	image.RGBA = &img

	// 讲图片发送给浏览器
	png.Encode(w, image)

}

// 获取短信验证码
func GetSmsCd(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	fmt.Println("------------------------获取短信验证码 GetSmsCd---------------------------------")
	// 接收web传的参数
	// 创建服务
	reg := consul.NewRegistry()
	server := grpc.NewService(
		micro.Registry(reg),
	)
	// 服务初始化
	server.Init()

	// 连接服务
	GetSmscdService := GetSmscdProto.NewGetSmscdService("go.micro.srv.GetSmscd", server.Client())

	// 接收参数
	mobile := p.ByName("mobile")
	// 获取url上的传参
	id := r.URL.Query()["id"][0]
	text := r.URL.Query()["text"][0]
	beego.Info("id:", id)
	beego.Info("text:", text)

	// 正则验证手机号
	// 创建正则句柄
	mobileReg := regexp.MustCompile(`0?(13|14|15|17|18|19)[0-9]{9}`)
	matchRes := mobileReg.MatchString(mobile)

	// 设置返回数据格式
	w.Header().Set("Content-Type", "application/json")
	if matchRes == false {
		response := map[string]interface{}{
			"error":  utils.RECODE_MOBILEERR,
			"errmsg": utils.RecodeText(utils.RECODE_MOBILEERR),
		}
		// 发送数据到前端
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
	}

	rsp, err := GetSmscdService.Call(context.TODO(), &GetSmscdProto.Request{
		Mobile: mobile,
		Id:     id,
		Text:   text,
	})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	// 返回数据
	response := map[string]interface{}{
		"errno":  rsp.Errno,
		"errmsg": rsp.Errmsg,
	}
	if err = json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

}

// 用户注册
func PostRet(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Println("------------------------用户注册---------------------------------")
	// 创建服务
	reg := consul.NewRegistry()
	server := grpc.NewService(
		micro.Registry(reg),
	)
	// 服务初始化
	server.Init()

	// 接收post发过来的数据
	// decode the incoming request as json
	var request map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	fmt.Println("mobile类型是：", reflect.TypeOf(request["mobile"]))
	fmt.Println("password类型是：", reflect.TypeOf(request["password"]))
	fmt.Println("sms_code类型是：", reflect.TypeOf(request["sms_code"]))
	if request["mobile"] == "" || request["password"] == "" || request["sms_code"] == "" {
		// 准备回传数据
		response := map[string]interface{}{
			"errno":  utils.RECODE_DATAERR,
			"errmsg": utils.RecodeText(utils.RECODE_DATAERR),
		}
		// 设置返回数据的格式
		w.Header().Set("Content-Type", "application/json")
		// 将数据返给前端
		err := json.NewEncoder(w).Encode(response)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		return
	}

	// 连接服务
	PostRetService := PostRetProto.NewPostRetService("go.micro.srv.PostRet", server.Client())
	rsp, err := PostRetService.Call(context.TODO(), &PostRetProto.Request{
		Mobile:   request["mobile"].(string),
		Password: request["password"].(string),
		SmsCode:  request["sms_code"].(string),
	})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// 读取cookie
	cookie, err := r.Cookie("userlogin")
	if err != nil || cookie.Value == "" {
		// 未获取到cookie或者cookie值为空
		// 创建cookie
		cookie := http.Cookie{
			Name:   "userlogin",
			Value:  rsp.SessionID,
			Path:   "/",
			MaxAge: 3600,
		}
		// 对浏览器进行cookie设置
		http.SetCookie(w, &cookie)
	}

	// 准备回传数据
	response := map[string]interface{}{
		"errno":  rsp.Errno,
		"errmsg": rsp.Errmsg,
	}
	// 设置返回数据的格式
	w.Header().Set("Content-Type", "application/json")
	// 将数据返给前端
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

}

// 用户注册
func PostLogin(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Println("------------------------用户登录 login---------------------------------")
	// 创建服务
	reg := consul.NewRegistry()
	server := grpc.NewService(
		micro.Registry(reg),
	)
	// 服务初始化
	server.Init()

	// 接收post发过来的数据
	// decode the incoming request as json
	var request map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	fmt.Println("mobile类型是：", reflect.TypeOf(request["mobile"]))
	fmt.Println("password类型是：", reflect.TypeOf(request["password"]))
	//fmt.Println("sms_code类型是：",reflect.TypeOf(request["sms_code"]))
	if request["mobile"] == "" || request["password"] == "" {
		// 准备回传数据
		response := map[string]interface{}{
			"errno":  utils.RECODE_DATAERR,
			"errmsg": utils.RecodeText(utils.RECODE_DATAERR),
		}
		// 设置返回数据的格式
		w.Header().Set("Content-Type", "application/json")
		// 将数据返给前端
		err := json.NewEncoder(w).Encode(response)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		return
	}

	PostLoginService := PostLoginProto.NewPostLoginService("go.micro.srv.PostLogin", server.Client())
	rsp, err := PostLoginService.Call(context.TODO(), &PostLoginProto.Request{
		Mobile:   request["mobile"].(string),
		Password: request["password"].(string),
	})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// 设置cookie
	cookie, err := r.Cookie("userlogin")
	if err != nil || cookie.Value == "" {
		// cookie 不存在，重新设置
		cookie := http.Cookie{
			Name:   "userlogin",
			Value:  rsp.Sessionid,
			Path:   "/",
			MaxAge: 600,
		}
		http.SetCookie(w, &cookie)
	}

	// 返回数据
	response := map[string]interface{}{
		"errno":  rsp.Errno,
		"errmsg": rsp.Errmsg,
	}

	// 设置返回数据格式
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

}

// 退出登录
func DeleteSession(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	beego.Info("请求退出登录，deletesession ...")
	// 创建服务
	reg := consul.NewRegistry()
	server := grpc.NewService(
		micro.Registry(reg),
	)

	// 服务初始化
	server.Init()

	DeleteSessionService := DeleteSessionProto.NewDeleteSessionService("go.micro.srv.DeleteSession", server.Client())

	// 获取cookie
	cookie, err := r.Cookie("userlogin")
	if err != nil || cookie.Value == " " {
		beego.Info("获取cookie失败，或者cookie位空：", err)
		// 没有cookie 本来就就是退出状态，直接返回错误
		response := map[string]interface{}{
			"errno":  utils.RECODE_DBERR,
			"errmsg": utils.RecodeText(utils.RECODE_DBERR),
		}
		// 设置返回格式
		w.Header().Set("Content-Type", "application/json")
		// 返回数据到前端
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		return
	}

	rsp, err := DeleteSessionService.Call(context.TODO(), &DeleteSessionProto.Request{
		Sessionid: cookie.Value,
	})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// 删除cookie存的sessionid
	newCookie := http.Cookie{
		Name:   "userlogin",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}
	http.SetCookie(w, &newCookie)

	// 返回数据
	response := map[string]interface{}{
		"errno":  rsp.Errno,
		"errmsg": rsp.Errmsg,
	}
	// 设置返回格式
	w.Header().Set("Content-Type", "application/json")
	// 返回数据到前端
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

// 获取用户信息
func GetUserInfo(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	beego.Info("请求获取用户信息...")
	// 创建服务
	reg := consul.NewRegistry()
	server := grpc.NewService(
		micro.Registry(reg),
	)
	// 服务初始化
	server.Init()

	// 创建服务客户端
	GetUserInfoService := GetUserInfoProto.NewGetUserInfoService("go.micro.srv.GetUserInfo", server.Client())
	// 读cookie获取sessionid
	cookie, err := r.Cookie("userlogin")
	if err != nil || cookie.Value == "" {
		// 构造返回数据
		response := map[string]interface{}{
			"errno":  utils.RECODE_DATAERR,
			"errmsg": utils.RecodeText(utils.RECODE_DATAERR),
		}
		// 设置返回头
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		return
	}

	// 请求服务端
	rsp, err := GetUserInfoService.Call(context.TODO(), &GetUserInfoProto.Request{
		Sessionid: cookie.Value,
	})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// 拼接返回参数
	data := make(map[string]interface{})
	data["name"] = rsp.Name
	data["mobile"] = rsp.Mobile
	data["user_id"] = rsp.Userid
	data["real_name"] = rsp.Realname
	data["avatar_url"] = utils.AddDomain2Url(rsp.Avatarurl)
	data["id_card"] = rsp.Idcard

	response := map[string]interface{}{
		"errno":  rsp.Errno,
		"errmsg": rsp.Errmsg,
		"data":   data,
	}
	// 设置返回头
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	return
}

// 上传用户头像
func PostAvatar(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	beego.Info("请求上传用户头像...")
	// 获取请求参数文件信息
	file, fileHeader, err := r.FormFile("avatar")
	if err != nil {
		// 准备回传数据
		response := map[string]interface{}{
			"errno":  utils.RECODE_DATAERR,
			"errmsg": utils.RecodeText(utils.RECODE_DATAERR),
		}
		// 设置返回数据的格式
		w.Header().Set("Content-Type", "application/json")
		if err = json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		return
	}
	beego.Info("文件大小：", fileHeader.Size)
	beego.Info("文件名：", fileHeader.Filename)

	// 创建一个文件大小的切片
	fileBuf := make([]byte, fileHeader.Size)
	// 将file数据读取到filebuf中
	_, err = file.Read(fileBuf)
	if err != nil {
		beego.Info("没有读取到图片数据流")
		// 准备回传数据
		response := map[string]interface{}{
			"errno":  utils.RECODE_DATAERR,
			"errmsg": utils.RecodeText(utils.RECODE_DATAERR),
		}
		// 设置返回数据的格式
		w.Header().Set("Content-Type", "application/json")
		if err = json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		return
	}

	// 通过cookie获取sessionid
	cookie, err := r.Cookie("userlogin")
	if err != nil || cookie.Value == "" {
		beego.Info("没有读取到cookie")
		// 准备回传数据
		response := map[string]interface{}{
			"errno":  utils.RECODE_DATAERR,
			"errmsg": utils.RecodeText(utils.RECODE_DATAERR),
		}
		// 设置返回数据的格式
		w.Header().Set("Content-Type", "application/json")
		if err = json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		return
	}

	// 创建服务
	reg := consul.NewRegistry()
	service := grpc.NewService(
		micro.Registry(reg),
	)
	// 初始化
	service.Init()

	// 获取客户端句柄
	PostAvatarService := PostAvatarProto.NewPostAvatarService("go.micro.srv.PostAvatar", service.Client())
	// 发起请求
	rsp, err := PostAvatarService.Call(context.TODO(), &PostAvatarProto.Request{
		Avatar:    fileBuf,
		FileSize:  fileHeader.Size,
		FileExt:   fileHeader.Filename,
		SessionId: cookie.Value,
	})

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// 处理返回数据，返回给浏览器
	data := make(map[string]string)
	data["avatar_url"] = utils.AddDomain2Url(rsp.AvatarUrl)
	response := map[string]interface{}{
		"errno":  rsp.Errno,
		"errmsg": rsp.Errmsg,
		"data":   data,
	}
	// 设置返回数据的格式
	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	return
}

// 更改用户名
func PutUserInfo(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	beego.Info("请求更改用户名...")
	w.Header().Set("Content-Type", "application/json")
	// 获取请求的参数name
	var request map[string]string
	response := make(map[string]interface{})
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		response["errno"] = utils.RECODE_PARAMERR
		response["errmsg"] = utils.RecodeText(utils.RECODE_PARAMERR)
		if err = json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		return
	}
	beego.Info("接收到参数name:", request["name"])
	if request["name"] == " " {
		beego.Info("参数name是空")
		response["errno"] = utils.RECODE_PARAMERR
		response["errmsg"] = utils.RecodeText(utils.RECODE_PARAMERR)
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		return
	}

	// 通过cookie获取sessionid
	cookie, err := r.Cookie("userlogin")
	if err != nil || cookie.Value == "" {
		beego.Info("没有sessionid")
		response["errno"] = utils.RECODE_PARAMERR
		response["errmsg"] = utils.RecodeText(utils.RECODE_PARAMERR)
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		return
	}

	// 创建并初始化服务
	reg := consul.NewRegistry()
	service := grpc.NewService(
		micro.Registry(reg),
	)
	service.Init()

	// 创建客户端服务
	putUserInfoService := PutUserInfoProto.NewPutUserInfoService("go.micro.srv.PutUserInfo", service.Client())
	// 请求服务端
	rsp, err := putUserInfoService.Call(context.TODO(), &PutUserInfoProto.Request{
		SessionId: cookie.Value,
		UserName:  request["name"],
	})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	data := make(map[string]string)
	data["name"] = rsp.UserName
	response["data"] = data
	response["errno"] = rsp.Errno
	response["errmsg"] = rsp.Errmsg
	if err = json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	return
}

// 查看实名认证
func GetUserAuth(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	beego.Info("请求获取用户信息...")
	// 创建服务
	reg := consul.NewRegistry()
	server := grpc.NewService(
		micro.Registry(reg),
	)
	// 服务初始化
	server.Init()

	// 创建服务客户端
	GetUserInfoService := GetUserInfoProto.NewGetUserInfoService("go.micro.srv.GetUserInfo", server.Client())
	// 读cookie获取sessionid
	cookie, err := r.Cookie("userlogin")
	if err != nil || cookie.Value == "" {
		// 构造返回数据
		response := map[string]interface{}{
			"errno":  utils.RECODE_DATAERR,
			"errmsg": utils.RecodeText(utils.RECODE_DATAERR),
		}
		// 设置返回头
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		return
	}

	// 请求服务端
	rsp, err := GetUserInfoService.Call(context.TODO(), &GetUserInfoProto.Request{
		Sessionid: cookie.Value,
	})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// 拼接返回参数
	data := make(map[string]interface{})
	data["name"] = rsp.Name
	data["mobile"] = rsp.Mobile
	data["user_id"] = rsp.Userid
	data["real_name"] = rsp.Realname
	data["avatar_url"] = utils.AddDomain2Url(rsp.Avatarurl)
	data["id_card"] = rsp.Idcard

	response := map[string]interface{}{
		"errno":  rsp.Errno,
		"errmsg": rsp.Errmsg,
		"data":   data,
	}
	// 设置返回头
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	return
}

// 提交认证信息
func PostUserAuth(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	beego.Info("请求处理认证信息...")
	// 接收参数
	request := make(map[string]string)
	response := make(map[string]string)
	// 设置返回输出格式
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		response["errno"] = utils.RECODE_PARAMERR
		response["errmsg"] = utils.RecodeText(utils.RECODE_PARAMERR)
		if err = json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		return
	}
	// 创建并初始化服务
	reg := consul.NewRegistry()
	service := grpc.NewService(
		micro.Registry(reg),
	)
	service.Init()
	// 创建客户端服务
	postUserAuthService := PostUserAuthProto.NewPostUserAuthService("go.micro.srv.PostUserAuth", service.Client())
	// 获取sessionid
	cookie, err := r.Cookie("userlogin")
	if err != nil || cookie.Value == "" {
		beego.Info("没有sessionid")
		response["errno"] = utils.RECODE_PARAMERR
		response["errmsg"] = utils.RecodeText(utils.RECODE_PARAMERR)
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		return
	}
	sessionId := cookie.Value
	// 请求服务端
	rsp, err := postUserAuthService.Call(context.TODO(), &PostUserAuthProto.Request{
		SessionId: sessionId,
		RealName:  request["real_name"],
		IdCard:    request["id_card"],
	})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	response["errno"] = rsp.Errno
	response["errmsg"] = rsp.Errmsg
	if err = json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	return
}

// 获取已发布房源
func GetUserHouses(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	beego.Info("请求获取已发布房源...")
	cookie, err := r.Cookie("userlogin")
	response := make(map[string]interface{})
	if err != nil || cookie.Value == "" {
		beego.Info("没有sessionid")
		response["errno"] = utils.RECODE_PARAMERR
		response["errmsg"] = utils.RecodeText(utils.RECODE_PARAMERR)
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		return
	}
	sessionId := cookie.Value
	// 创建并初始化服务
	reg := consul.NewRegistry()
	service := grpc.NewService(
		micro.Registry(reg),
	)
	service.Init()
	// 创建客户端服务
	getUserHousesService := GetUserHousesProto.NewGetUserHousesService("go.micro.srv.GetUserHouses", service.Client())
	// 请求服务端
	rsp, err := getUserHousesService.Call(context.TODO(), &GetUserHousesProto.Request{
		SessionId: sessionId,
	})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	// 将服务端返回的二进制数据流解码到切片中
	var houseList []models.House
	_ = json.Unmarshal(rsp.Mix, &houseList)
	var houses []interface{}
	// 遍历返回完整房屋信息
	for _, value := range houseList {
		// 获取到有用的添加到切片当中
		houses = append(houses, value.To_house_info())
	}
	// 创建一个data的map
	data := make(map[string]interface{})
	data["houses"] = houses
	// 整理返回参数
	response["errno"] = rsp.Errno
	response["errmsg"] = rsp.Errmsg
	response["data"] = data
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

}
