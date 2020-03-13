package model

// Attachment 附件模型
type Attachment struct {
	BaseModel
	AttachmentTypeId        int64  `gorm:"primary_key" json:"attachment_type_id,omitempty"`
	AttachmentTypeTableName string `json:"attachment_type_table_name,omitempty"`
	AttachmentTypeFieldName string `json:"attachment_type_field_name,omitempty"`
	AttachmentTypeName      string `json:"attachment_type_name,omitempty"`
	AttachmentTypeCnName    string `json:"attachment_type_cn_name,omitempty"`
	QiniuPrefix             string `json:"qiniu_prefix,omitempty"`
}

// TableName 表名
func (Attachment) TableName() string {
	return "rxt_attachment"
}
