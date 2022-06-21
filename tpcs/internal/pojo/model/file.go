package model

import "github.com/jinzhu/gorm"

type File struct {
	Id *int `gorm:"column:ID;primaryKey" json:"id"`
	*File4Add
	GmtCreate *string `gorm:"column:GMT_CREATE" json:"gmtCreate"`
}

type File4Add struct {
	UserId    *int    `gorm:"column:USER_ID" json:"userId"`
	TruthName *string `gorm:"column:TRUTH_NAME" json:"truthName"`
	FileName  *string `gorm:"column:FILE_NAME" json:"fileName"`
}

func (f *File) TableName() string {
	return "file_info"
}

func (f4a *File4Add) TableName() string {
	return "file_info"
}

func (f *File) Create(db *gorm.DB) error {
	return db.Create(&f).Error
}

func (f4a *File4Add) Create(db *gorm.DB) error {
	return db.Create(&f4a).Error
}

//type UploadStatus struct {
//	// 已读数据
//	BytesRead *int64 `json:"bytesRead"`
//	// 文件总数据
//	ContentLength *int64 `json:"contentLength"`
//	// 第几个文件
//	Items *int64 `json:"items"`
//	// 开始时间
//	StartTime *int64 `json:"startTime"`
//	// 已用时间
//	UseTime *int64 `json:"useTime"`
//	// 完成百分比
//	Percent *int `json:"percent"`
//}
