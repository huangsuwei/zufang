package utils

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	//"github.com/tedcy/fdfs_client"
	fdfs_client "github.com/sanxia/gfs"
)

/* 将url加上 http://IP:PROT/  前缀 */
//http:// + 127.0.0.1 + ：+ 8080 + 请求

func AddDomain2Url(url string) (domain_url string) {
	domain_url = "http://" + G_fastdfs_addr + ":" + G_fastdfs_port + "/" + url

	return domain_url
}

// 密码加密函数
func Md5String(str string) string {
	// 创建一个md5对象
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

// 上传二进制文件到fdfs
func UploadByBuffer(buffer []byte, fileExt string) (fileId string, err error) {
	client, err := fdfs_client.NewFdfsClient([]string{G_fastdfs_addr + ":" + G_fastdfs_port})
	if err != nil {
		fmt.Println("fdfs_client 创建句柄失败：", err.Error())
		fileId = ""
		return
	}
	rsp, err := client.UploadByBuffer(buffer, fileExt)
	if err != nil {
		fmt.Println("上传失败：", err.Error())
		fileId = ""
		return
	}
	fileId = rsp.FileId
	return fileId, nil
}

//func UploadByBuffer(buffer []byte, fileExt string) (fileId string, err error) {
//	client, err := fdfs_client.NewClientWithConfig("/home/guan/myproject/sss/IhomeWeb/conf/client.conf")
//	if err != nil {
//		fmt.Println("fdfs_client 创建句柄失败：", err.Error())
//		fileId = ""
//		return
//	}
//	defer client.Destory()
//
//	fileId, err = client.UploadByBuffer(buffer, fileExt)
//	if err != nil {
//		fmt.Println("上传失败：", err.Error())
//		fileId = ""
//		return
//	}
//	fmt.Println("上传成功，fileId:", fileId)
//
//	return fileId, nil
//}
