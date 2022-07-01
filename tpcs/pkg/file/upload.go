package file

import (
	"github.com/gin-contrib/sessions"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"time"
	"tpcs/global"
)

var Session sessions.Session

type UploadStatus struct {
	// 已读数据
	BytesRead uint64 `json:"bytesRead"`
	// 文件总数据
	ContentLength int64 `json:"contentLength"`
	// 第几个文件
	Items int `json:"items"`
	// 开始时间
	StartTime int64 `json:"startTime"`
	// 已用时间
	UseTime int64 `json:"useTime"`
	// 完成百分比
	Percent int `json:"percent"`
}

// SaveFile4MdImg 保存所上传的文件markdown图片
func SaveFile4MdImg(file *multipart.FileHeader, dst string) error {
	src, err := file.Open()
	if err != nil {
		global.Logger.Errorf("file.Open err: %v\n", err)
		return err
	}
	defer src.Close()

	out, err := os.Create(dst)
	if err != nil {
		global.Logger.Errorf("os.Create err: %v\n", err)
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, src)
	if err != nil {
		global.Logger.Errorf("io.Copy err: %v\n", err)
		return err
	}

	return nil
}

func (us *UploadStatus) Write(p []byte) (int, error) {
	n := len(p)
	us.BytesRead += uint64(n)

	us.Percent = int(100 * us.BytesRead / uint64(us.ContentLength))
	Session.Delete("upload_status")
	Session.Options(sessions.Options{
		Path:     "/",
		Domain:   "",
		MaxAge:   -1,
		Secure:   false,
		HttpOnly: true,
		SameSite: http.SameSiteDefaultMode,
	})
	_ = Session.Save()
	Session.Set("upload_status", us)
	Session.Options(sessions.Options{
		Path:     "/",
		Domain:   "",
		MaxAge:   0,
		Secure:   false,
		HttpOnly: true,
		SameSite: http.SameSiteDefaultMode,
	})
	err := Session.Save()
	if err != nil {
		global.Logger.Errorf("session.Save() error: %v\n", err)
		return 0, err
	}

	return n, nil
}

// SaveFile 保存所上传的文件，
// 该方法主要是通过调用 os.Create 方法创建目标地址的文件，
// 再通过 file.Open 方法打开源地址的文件，
// 结合 io.Copy 方法实现两者之间的文件内容拷贝
func SaveFile(item int, file *multipart.FileHeader, dst string) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	startTime := time.Now().UnixMilli()
	uploadStatus := &UploadStatus{
		ContentLength: file.Size,
		Items:         item,
		StartTime:     startTime,
		UseTime:       time.Now().UnixMilli() - startTime,
	}
	_, err = io.Copy(out, io.TeeReader(src, uploadStatus))
	if err != nil {
		global.Logger.Errorf("io.Copy(out, io.TeeReader(src, uploadStatus)) err: %v", err)
		return err
	}
	return nil
}
