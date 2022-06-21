package file

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"time"
	"tpcs/global"
)

var session sessions.Session

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

func (us *UploadStatus) Write(p []byte) (int, error) {
	n := len(p)
	us.BytesRead += uint64(n)

	us.Percent = int(100 * us.BytesRead / uint64(us.ContentLength))
	session.Delete("upload_status")
	session.Options(sessions.Options{
		Path:     "/",
		Domain:   "",
		MaxAge:   -1,
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
	})
	//_ = session.Save()
	session.Set("upload_status", us)
	session.Options(sessions.Options{
		Path:     "/",
		Domain:   "",
		MaxAge:   0,
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
	})
	_ = session.Save()

	return n, nil
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

// SaveFile 保存所上传的文件，
// 该方法主要是通过调用 os.Create 方法创建目标地址的文件，
// 再通过 file.Open 方法打开源地址的文件，
// 结合 io.Copy 方法实现两者之间的文件内容拷贝
func SaveFile(c *gin.Context, item int, file *multipart.FileHeader, dst string) error {
	session = sessions.DefaultMany(c, "upload_status")
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
