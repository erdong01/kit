package model

type Exam struct {
	BaseModel
	ExamId                   int64              `gorm:"primary_key;column:exam_id" json:"exam_id,omitempty"`
	ExamNo                   int64              `json:"exam_no,omitempty"`
	ExamName                 string             `json:"exam_name,omitempty"`
	ExamStatus               int8               `json:"exam_status,omitempty"`
	ExamLevelType            int8               `json:"exam_level_type,omitempty"`
	ExamScore                float64            `json:"exam_score,omitempty"`
	ExamTime                 float32            `json:"exam_time,omitempty"`
	ExamAnswerType           int8               `json:"exam_answer_type,omitempty"`
	FrontUserNo              int64              `json:"front_user_no,omitempty"`
	CompanyNo                int64              `json:"company_no,omitempty"`
	CampusNo                 int64              `json:"campus_no,omitempty"`
	SubjectId                int64              `json:"subject_id,omitempty"`
	GradeId                  int64              `json:"grade_id,omitempty"`
	GradeChildrenId          int64              `json:"grade_children_id,omitempty"`
	SubjectParentId          int64              `json:"subject_parent_id,omitempty"`
	EditionNo                int64              `json:"edition_no,omitempty"`
	IsReview                 int8               `json:"is_review,omitempty"`
	ExamTypeCode             int                `json:"exam_type_code,omitempty"`
	ExamGroupType            int8               `json:"exam_group_type,omitempty"`
	ExamActualScore          float64           `json:"exam_actual_score,omitempty"`
	ExamActualTime           float32            `json:"exam_actual_time,omitempty"`
	ExamWordStatus           int8               `json:"exam_word_status,omitempty"`
	ScoringMethod            int8               `json:"scoring_method,omitempty"`
	ExamQuestionCount        int32              `json:"exam_question_count,omitempty"`
	ExamQuestionCorrectCount int32              `json:"exam_question_correct_count,omitempty"`
	ExamQuestionErrorCount   int32              `json:"exam_question_error_count,omitempty"`
	ExamQuestionType         []ExamQuestionType `gorm:"ForeignKey:exam_no;AssociationForeignKey:exam_no"`
}

func (Exam) TableName() string {
	return "rxt_exam"
}
