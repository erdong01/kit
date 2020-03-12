package model

type PaperAnalysis struct {
	BaseModel
	PaperAnalysisId   int64 `gorm:"primary_key" json:"paper_analysis_id,omitempty"`
	PaperAnalysisNo   int64 `json:"paper_analysis_no,omitempty"`
	ProvinceId        int64 `json:"province_id,omitempty"`
	CityId            int64 `json:"city_id,omitempty"`
	DistrictId        int64 `json:"district_id,omitempty"`
	GradeChildrenId   int64 `json:"grade_children_id,omitempty"`
	SubjectId         int64 `json:"subject_id,omitempty"`
	PaperAnalysisType int8  `json:"paper_analysis_type,omitempty"`
}

func (PaperAnalysis) TableName() string {
	return "rxt_paper_analysis"
}
