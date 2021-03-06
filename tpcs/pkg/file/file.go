package file

import (
	"os"
	"path"
	"strings"
	"tpcs/global"
)

type FileType int

const (
	TypeImage FileType = iota + 1
	TypeOther
)

// GetFileName 获取文件名称，
// 先是通过获取文件后缀并筛出原始文件名进行 MD5 加密，最后返回经过加密处理后的文件名
func GetFileName(name string) string {
	ext := GetFileExt(name)
	fileName := strings.TrimSuffix(name, ext)

	return fileName + ext
}

// GetFileExt 获取文件后缀，
// 主要是通过调用 path.Ext 方法进行循环查找”.“符号，最后通过切片索引返回对应的文化后缀名称
func GetFileExt(name string) string {
	return path.Ext(name)
}

// GetSavePath 获取文件保存地址，
// 这里直接返回配置中的文件保存目录即可，也便于后续的调整
func GetSavePath() string {
	return global.AppSetting.FileUploadPath
}

// SavePathExists 检查保存目录是否存在，
// 通过调用 os.Stat 方法获取文件的描述信息 FileInfo，
// 并调用 os.IsNotExist 方法进行判断，
// 其原理是利用 os.Stat 方法所返回的 error 值与系统中所定义的 oserror.ErrNotExist 进行判断，
// 以此达到校验效果
func SavePathExists(dst string) bool {
	_, err := os.Stat(global.AppSetting.UploadDir + dst)
	return !os.IsNotExist(err)
}

// CheckContainExt 检查文件后缀是否包含在约定的后缀配置项中，
// 需要的是所上传的文件的后缀有可能是大写、小写、大小写等，
// 因此我们需要调用 strings.ToUpper 方法统一转为大写（固定的格式）来进行匹配
//func CheckContainExt(name string) bool {
//	ext := GetFileExt(name)
//	ext = strings.ToUpper(ext)
//
//
//	return false
//}

// CheckMaxSize 检查文件大小是否超出最大大小限制
//func CheckMaxSize(t FileType, f multipart.File) bool {
//	content, _ := ioutil.ReadAll(f)
//	size := len(content)
//	switch t {
//	case TypeImage:
//		if size >= global.AppSetting.UploadImageMaxSize*1024*1024 {
//			return true
//		}
//	}
//
//	return false
//}

// CheckPermission 检查文件权限是否足够，
// 与 SavePathExists 方法原理一致，是利用 oserror.ErrPermission 进行判断
func CheckPermission(dst string) bool {
	_, err := os.Stat(dst)
	return os.IsPermission(err)
}

// CreateSavePath 创建在上传文件时所使用的保存目录，
// 在方法内部调用的 os.MkdirAll 方法，
// 该方法将会以传入的 os.FileMode 权限位去递归创建所需的所有目录结构，
// 若涉及的目录均已存在，则不会进行任何操作，直接返回 nil
func CreateSavePath(dst string, perm os.FileMode) error {
	err := os.MkdirAll(dst, perm)
	if err != nil {
		return err
	}

	return nil
}

// DeleteFile 删除文件
func DeleteFile(filePath string) error {
	err := os.Remove(filePath)
	if err != nil {
		return err
	}

	return nil
}

// CheckNotExist 检查文件是否存在
func CheckNotExist(src string) bool {
	_, err := os.Stat(src)

	return os.IsNotExist(err)
}

// IsNotExistMkDir 如果不存在则新建文件夹
func IsNotExistMkDir(src string) error {
	if notExist := CheckNotExist(src); notExist == true {
		if err := MkDir(src); err != nil {
			return err
		}
	}

	return nil
}

// MkDir 新建文件夹
func MkDir(src string) error {
	err := os.MkdirAll(src, os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}

// Open 打开文件
func Open(name string, flag int, perm os.FileMode) (*os.File, error) {
	f, err := os.OpenFile(name, flag, perm)
	if err != nil {
		return nil, err
	}

	return f, nil
}
