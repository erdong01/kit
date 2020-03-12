package model

type Chapter struct {
	BaseModel
	ChapterId       int64  `gorm:"primary_key" json:"chapter_id,omitempty"`
	ChapterParentId int64  `json:"chapter_parent_id,omitempty"`
	ChapterNo       int64  `json:"chapter_no,omitempty"`
	ChapterSort     int    `json:"chapter_sort,omitempty"`
	ChapterName     string `json:"chapter_name,omitempty"`
	BookNo          int64  `json:"book_no,omitempty"`

	ChapterKnowledge []ChapterKnowledge `gorm:"ForeignKey:ChapterNo;AssociationForeignKey:ChapterNo"`
}

func (Chapter) TableName() string {
	return "rxt_chapter"
}
