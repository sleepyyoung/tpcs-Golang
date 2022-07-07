package file

import (
	"encoding/json"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"tpcs/internal/pojo"
	"tpcs/internal/pojo/model"
	"tpcs/internal/service"
	fileService "tpcs/internal/service/file"
	userService "tpcs/internal/service/user"
	"tpcs/pkg/app"
	"tpcs/pkg/file"
	"tpcs/pkg/logger"
	"tpcs/pkg/upload"
)

type File struct{}

func NewFile() File {
	return File{}
}

// Upload4MdImg markdown图片上传
func (p File) Upload4MdImg(c *gin.Context) {
	fileSvc := fileService.New(c.Request.Context())
	response := app.NewResponse(c)

	imageFile, err := c.FormFile("editormd-image-file")
	if err != nil {
		logger.Errorf("c.FormFile() err: %v\n", err)
		response.ToFailResultResponse(pojo.ResultMsg_TryAgainLater)
		return
	}

	url, err := fileSvc.Upload4MdImg(imageFile)
	if err != nil {
		logger.Errorf("fileSvc.Upload4MdImg() err: %v\n", err)
		response.Ctx.JSON(
			http.StatusRequestTimeout,
			map[string]interface{}{
				"success": 0,
				"message": pojo.ResultMsg_TryAgainLater,
			})
		return
	}

	response.Ctx.JSON(
		http.StatusOK,
		map[string]interface{}{
			"success": 1,
			"message": nil,
			"url":     url,
		})
	return
}

// List 文件列表
func (p File) List(c *gin.Context) {
	userSvc := userService.New(c.Request.Context())
	fileSvc := fileService.New(c.Request.Context())
	response := app.NewResponse(c)

	type innerFile struct {
		Id        *int    `json:"id"`
		Name      *string `json:"name"`
		Path      *string `json:"path"`
		Teacher   *string `json:"teacher"`
		GmtCreate *string `json:"gmtCreate"`
	}

	var request service.ListRequest
	err := c.ShouldBindQuery(&request)
	if err != nil {
		logger.Errorf("c.ShouldBindWith err: %v", err)
		response.ToResponse(pojo.Result{
			Code:  0,
			Msg:   pojo.ResultMsg_FormParseErr,
			Count: 0,
			Data:  nil,
		})
		return
	}

	session := sessions.DefaultMany(c, "user")
	userI := session.Get("user")
	user := userI.(model.User)

	var count int
	var fileList []model.File
	var innerFileList []innerFile

	fileList, count, err = fileSvc.FileList(*user.Id, &request)
	if err != nil {
		logger.Errorf("svc.FileList err: %v", err)
		response.ToResponse(pojo.Result{
			Code:  0,
			Msg:   pojo.ResultMsg_TryAgainLater,
			Count: 0,
			Data:  nil,
		})
		return
	}

	for _, f := range fileList {
		u, err := userSvc.GetUserByUserId(*f.UserId)
		if err != nil {
			logger.Errorf("svc.GetUserByUserId err: %v", err)
			response.ToResponse(pojo.Result{
				Code:  0,
				Msg:   pojo.ResultMsg_TryAgainLater,
				Count: 0,
				Data:  nil,
			})
			return
		}

		time := strings.ReplaceAll(*f.GmtCreate, "T", " ")
		time = strings.ReplaceAll(time, "+08:00", "")

		innerFileList = append(innerFileList, innerFile{
			Id:        f.Id,
			Name:      f.TruthName,
			Path:      f.FileName,
			Teacher:   u.Username,
			GmtCreate: &time,
		})
	}
	response.ToResponse(pojo.Result{
		Code:  0,
		Msg:   nil,
		Count: count,
		Data:  innerFileList,
	})
}

// Upload 文件上传
func (p File) Upload(c *gin.Context) {
	upload.Session = sessions.DefaultMany(c, "upload_status")
	fileSvc := fileService.New(c.Request.Context())
	response := app.NewResponse(c)
	w := response.Ctx.Writer
	w.Header().Set("Cache-Control", "max-age=86400")
	w.Header().Set("Content-type", "application/octet-stream")

	session := sessions.DefaultMany(c, "user")
	userI := session.Get("user")
	user := userI.(model.User)

	multipartForm, err := c.MultipartForm()
	if err != nil {
		logger.Errorf("c.MultipartForm() err: %v\n", err)
		response.ToFailResultResponse(pojo.ResultMsg_TryAgainLater)
		return
	}

	for index, fileHeader := range (multipartForm.File)["files"] {
		err := fileSvc.Upload(index+1, fileHeader, *user.Id)
		if err != nil {
			response.ToResponse(map[string]string{"result": "fail"})
			return
		}
	}

	response.ToResponse(map[string]string{"result": "success"})
	return
}

// UploadStatus 文件上传进度
func (p File) UploadStatus(c *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			logger.Errorf("从session中获取文件上传进度失败! : %v\n", r)
		}
	}()

	response := app.NewResponse(c)
	session := sessions.DefaultMany(c, "upload_status")
	uploadStatus := session.Get("upload_status").(upload.UploadStatus)
	response.ToResponse(uploadStatus)
	return
}

// Download 文件下载
func (p File) Download(c *gin.Context) {
	fileSvc := fileService.New(c.Request.Context())
	response := app.NewResponse(c)
	w := response.Ctx.Writer

	fileId, err := strconv.Atoi(c.Param("fileId"))
	if err != nil {
		logger.Errorf("strconv.Atoi(c.Param(\"fileId\")) err: %v\n", err)
		response.Ctx.HTML(http.StatusOK, "404.tmpl", gin.H{
			"message": pojo.ResultMsg_FormParseErr,
		})
		return
	}

	f, err := fileSvc.GetFileById(fileId)
	if err != nil {
		logger.Errorf("svc.GetFileById err: %v\n", err)
		response.Ctx.HTML(http.StatusOK, "404.tmpl", gin.H{
			"message": pojo.ResultMsg_FileNotFound,
		})
		return
	}
	if f == nil {
		response.Ctx.HTML(http.StatusOK, "404.tmpl", gin.H{
			"message": pojo.ResultMsg_FileNotFound,
		})
		return
	}
	if !file.SavePathExists(*f.FileName) {
		response.Ctx.HTML(http.StatusOK, "404.tmpl", gin.H{
			"message": pojo.ResultMsg_FileNotFound,
		})
		return
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("X-Accel-Redirect", *f.FileName)
	w.Header().Set("X-Accel-Charset", "utf-8")
	w.Header().Set("Content-Disposition", "attachment; filename="+*f.TruthName)
	return
}

// Delete 文件删除
func (p File) Delete(c *gin.Context) {
	fileSvc := fileService.New(c.Request.Context())
	response := app.NewResponse(c)

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		logger.Errorf("strconv.Atoi(c.Param(\"id\")) err: %v\n", err)
		response.ToFailResultResponse(pojo.ResultMsg_TryAgainLater)
		return
	}

	var result pojo.Result
	result = fileSvc.DeleteFile(id)
	response.ToResultResponse(&result)
}

// BatchDelete 文件批量删除
func (p File) BatchDelete(c *gin.Context) {
	fileSvc := fileService.New(c.Request.Context())
	response := app.NewResponse(c)

	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		logger.Errorf("ioutil.ReadAll err: %v\n", err)
		response.ToFailResultResponse(pojo.ResultMsg_TryAgainLater)
		return
	}

	ids := make([]int, 15)
	err = json.Unmarshal(body, &ids)
	if err != nil {
		logger.Errorf("json.Unmarshal err: %v\n", err)
		response.ToFailResultResponse(pojo.ResultMsg_TryAgainLater)
		return
	}

	var result pojo.Result
	result = fileSvc.BatchDeleteFile(ids)
	response.ToResultResponse(&result)
}
