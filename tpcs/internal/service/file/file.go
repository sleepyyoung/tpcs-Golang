package file

import (
	"github.com/gin-gonic/gin"
	"mime/multipart"
	"os"
	"strconv"
	"strings"
	"time"
	"tpcs/global"
	"tpcs/internal/pojo"
	"tpcs/internal/pojo/model"
	"tpcs/internal/service"
	"tpcs/pkg/file"
)

// GetFileById 通过id获取文件
func (svc *Service) GetFileById(id int) (*model.File, error) {
	return svc.dao.GetFileById(id)
}

// Upload4MdImg markdown图片上传
func (svc *Service) Upload4MdImg(fileHeader *multipart.FileHeader) (string, error) {
	truthName := file.GetFileName(fileHeader.Filename)
	splits := strings.Split(truthName, ".")
	suffix := splits[len(splits)-1]

	savePath := global.AppSetting.MdImgUploadPath
	if !file.SavePathExists(savePath) {
		if err := file.CreateSavePath(savePath, os.ModePerm); err != nil {
			global.Logger.Errorf("图片上传失败！原因: %v\n", err)
			//return pojo.Result{Success: pojo.ResultSuccess_False, Msg: pojo.ResultMsg_TryAgainLater}
			return "", err
		}
	}

	formatFileName := strconv.FormatInt(time.Now().UnixMilli(), 10) + "." + suffix
	dst := savePath + "/" + formatFileName

	if err := file.SaveFile4MdImg(fileHeader, dst); err != nil {
		global.Logger.Errorf("图片上传失败！原因: %v\n", err)
		//return pojo.Result{Success: pojo.ResultSuccess_False, Msg: pojo.ResultMsg_TryAgainLater}
		return "", err
	}

	//return pojo.Result{Success: pojo.ResultSuccess_True}
	return "/files/" + formatFileName, nil
}

// FileList 获取文件列表
func (svc *Service) FileList(userId int, param *service.ListRequest) ([]model.File, int, error) {
	return svc.dao.FileList(userId, param.Page, param.Limit)
}

// Upload 上传文件
func (svc *Service) Upload(c *gin.Context, item int, fileHeader *multipart.FileHeader, userId int) error {
	truthName := file.GetFileName(fileHeader.Filename)
	splits := strings.Split(truthName, ".")
	suffix := splits[len(splits)-1]

	savePath := global.AppSetting.FileUploadPath
	if !file.SavePathExists(savePath) {
		if err := file.CreateSavePath(savePath, os.ModePerm); err != nil {
			global.Logger.Errorf("文件上传失败！原因: %v\n", err)
			return err
		}
	}

	formatFileName := strconv.FormatInt(time.Now().UnixNano(), 10) + "." + suffix
	dst := savePath + "/" + formatFileName
	if err := file.SaveFile(item, fileHeader, dst); err != nil {
		global.Logger.Errorf("文件上传失败！原因: %v\n", err)
		return err
	}

	fileName := "/uploads/" + formatFileName
	err := svc.dao.AddFile(model.File4Add{
		TruthName: &truthName,
		UserId:    &userId,
		FileName:  &fileName,
	})
	if err != nil {
		global.Logger.Errorf("文件上传失败！原因: %v\n", err)
		return err
	}
	return nil
}

// DeleteFile 删除文件
func (svc *Service) DeleteFile(id int) pojo.Result {
	err := svc.dao.DeleteFile(id)
	if err != nil {
		global.Logger.Errorf("文件删除失败！原因: %v\n", err)
		return pojo.Result{Success: pojo.ResultSuccess_False, Msg: pojo.ResultMsg_TryAgainLater}
	}
	return pojo.Result{Success: pojo.ResultSuccess_True}
}

// BatchDeleteFile 批量删除文件
func (svc *Service) BatchDeleteFile(ids []int) pojo.Result {
	err := svc.dao.BatchDeleteFile(ids)
	if err != nil {
		global.Logger.Errorf("批量删除失败！原因: %v\n", err)
		return pojo.Result{Success: pojo.ResultSuccess_False, Msg: pojo.ResultMsg_TryAgainLater}
	}
	return pojo.Result{Success: pojo.ResultSuccess_True}
}
