package model

type QuestionDivision struct {
	BaseModel
	QuestionDivisionId int64 `gorm:"primary_key" json:"question_division_id,omitempty"`
	QuestionNo         int64 `json:"question_no,omitempty"`
	ProvinceId         int64 `json:"province_id,omitempty"`
	CityId             int64 `json:"city_id,omitempty"`
	DistrictId         int64 `json:"district_id,omitempty"`
}

func (QuestionDivision) TableName() string {
	return "rxt_question_difficulty"
}
