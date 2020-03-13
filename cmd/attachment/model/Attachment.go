package model

// Attachment 附件模型
type Attachment struct {
	BaseModel
	AttachmentId 	 int64 	`gorm:"primary_key" json:"attachment_id,omitempty"`
	AttachmentTypeId int64 	`json:"attachment_type_id,omitempty"`
	DataPrimaryId 	 int64 	`json:"attachment_type_id,omitempty"`
	FileSize 		 int64 	`json:"file_size,omitempty"`
	ExtensionName 	 string `json:"extension_name,omitempty"`
	OriginalFileName string `json:"original_file_name,omitempty"`
	MimeType 		 string `json:"mime_type,omitempty"`
	QiniuPath 		 string `json:"qiniu_path,omitempty"`
	IsCompress 		 int8 	`json:"is_compress,omitempty"`
}

// TableName 表名
func (Attachment) TableName() string {
	return "rxt_attachment"
}