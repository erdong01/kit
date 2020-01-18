package model

type QuestionTypeDivision struct {
	BaseModel
	QuestionTypeDivisionId int64 `gorm:"primary_key" json:"question_type_division_id,omitempty"`
	SubjectId              int64 `json:"subject_id,omitempty"`
	QuestionTypeId         int64 `json:"question_type_id,omitempty"`
	QuestionNumber         int   `json:"question_number,omitempty"`
	ProvinceId             int64 `json:"province_id,omitempty"`
	CityId                 int64 `json:"city_id,omitempty"`
	DistrictId             int64 `json:"district_id,omitempty"`
	IsOnline               int8  `json:"district_id,omitempty"`
}

func (QuestionTypeDivision) TableName() string {
	return "rxt_question_type_division"
}
