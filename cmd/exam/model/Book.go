package model

type Book struct {
	BaseModel
	BookId    int64  `gorm:"primary_key" json:"book_id,omitempty"`
	BookName  string `json:"book_name,omitempty"`
	BookNo    int64  `json:"book_no,omitempty"`
	BookSort  int    `json:"book_sort,omitempty"`
	BookImg   string `json:"book_img,omitempty"`
	EditionNo int64  `json:"edition_no,omitempty"`
	GradeId   int64  `json:"grade_id,omitempty"`
	PhaseId   int64  `json:"phase_id,omitempty"`
	SubjectId int64  `json:"subject_id,omitempty"`
	BookVol   int64  `json:"book_vol,omitempty"`
}

func (Book) TableName() string {
	return "rxt_book"
}
