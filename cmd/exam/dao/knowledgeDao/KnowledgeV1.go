package knowledgeDao

import (
	"github.com/jinzhu/gorm"
	"rxt/cmd/exam/dao/base"
	"rxt/cmd/exam/model"
)

type V1 struct {
	base.Dao
}
type Review struct {
	model.Knowledge
	model.KnowledgeReview
}

type ChapterKnowledge struct {
	model.ChapterKnowledge
	KnowledgeTeachingMinute float64 `json:"knowledge_teaching_minute"`
	KnowledgeParentId       int64   `json:"knowledge_parent_id"`
	KnowledgeId             int64   `json:"knowledge_id"`
	KnowledgeName           string  `json:"knowledge_name"`
}

func (c V1) EditionKnowledge(param Param) ([]Review, []ChapterKnowledge, map[int64]int) {
	var review []Review
	var chapterKnowledge []ChapterKnowledge
	var chapterKnowledgeNew []ChapterKnowledge
	var chapterKnowledgeSort map[int64]int
	chapterKnowledgeSort = make(map[int64]int)
	if param.IsReview == 2 {
		c.Db.
			Select("rxt_chapter_knowledge.chapter_no,chapter_knowledge_sort,chapter_knowledge_id,"+
				"rxt_knowledge_info.knowledge_teaching_minute,"+
				"rxt_knowledge.knowledge_no,rxt_knowledge.knowledge_id,rxt_knowledge.knowledge_parent_id,rxt_knowledge.knowledge_name").
			Joins("INNER JOIN rxt_knowledge ON rxt_knowledge.knowledge_no = rxt_chapter_knowledge.knowledge_no").
			Joins("LEFT JOIN rxt_knowledge_info ON rxt_knowledge.knowledge_no = rxt_knowledge_info.knowledge_no").
			Preload("TeachingPlanQuestionMany", func(db *gorm.DB) *gorm.DB {
				return db.Where("teaching_plan_question_type = ?", 1).
					Where("teaching_plan_question_level_type = ?", 1)
			}).
			Where("rxt_chapter_knowledge.edition_no = ?", param.EditionNo).
			Order("chapter_knowledge_sort asc").
			Find(&chapterKnowledge)
		for chapterKnowledgeK, chapterKnowledgeV := range chapterKnowledge {
			chapterKnowledgeSort[chapterKnowledgeV.KnowledgeNo] = chapterKnowledgeK

		}
		chpaterTree := c.EditionChapterTree(param.EditionNo, 0, 0)

		for _, chpaterTreeV := range chpaterTree {
			for _, chapterKnowledgeV := range chapterKnowledge {
				if chpaterTreeV.ChapterNo == chapterKnowledgeV.ChapterNo {
					chapterKnowledgeNew = append(chapterKnowledgeNew, chapterKnowledgeV)
				}
			}
		}
	} else {
		c.Db.Table("rxt_knowledge").
			Select("rxt_knowledge_review.*,rxt_knowledge.*,rxt_knowledge_info.*").
			Joins("LEFT JOIN rxt_knowledge_review ON rxt_knowledge.knowledge_no = rxt_knowledge_review.knowledge_no AND rxt_knowledge_review.edition_no = ? AND  rxt_knowledge_review.is_review_knowledge = 1", param.EditionNo).
			Joins("LEFT JOIN rxt_knowledge_info ON rxt_knowledge_info.knowledge_no = rxt_knowledge.knowledge_no").
			Where("subject_id = ?", param.SubjectId).
			Scan(&review)

		for reviewK, reviewV := range review {
			chapterKnowledgeSort[reviewV.Knowledge.KnowledgeNo] = reviewK
		}
	}
	return review, chapterKnowledgeNew, chapterKnowledgeSort
}

type preKnowledge struct {
	model.Knowledge
	model.KnowledgeReview
}

func (c V1) ListBySubjectId(subjectId int64, isReview int8, editionNo int64) []preKnowledge {

	var preKnowledgeArr []preKnowledge
	c.Db.Where("rxt_knowledge.subject_id = ?", subjectId).
		Table("rxt_knowledge").
		Joins("LEFT JOIN rxt_knowledge_review ON rxt_knowledge_review.knowledge_no = rxt_knowledge.knowledge_no"+
			" AND rxt_knowledge_review.is_review_knowledge = ? "+
			"AND rxt_knowledge_review.edition_no = ?", isReview, editionNo).
		Scan(&preKnowledgeArr)
	return preKnowledgeArr
}

func (c V1) GetPaperAnalysis(whereMap map[string]interface{}) (model.PaperAnalysis, bool) {
	var paperAnalysis model.PaperAnalysis
	res := c.Db.Where(whereMap).
		First(&paperAnalysis)

	return paperAnalysis, res.RecordNotFound()
}

type BookChapter struct {
	model.Book
	Chapter []model.Chapter
}

func (c V1) EditionChapterTree(editionNo int64, gradeId int64, bookNo int64) []model.Chapter {

	var bookArr []BookChapter
	bookDb := c.Db
	bookDb = bookDb.Table("rxt_book").Where("edition_no = ?", editionNo)
	if bookNo > 0 {
		bookDb = bookDb.Where("book_no = ?", bookNo)
	}
	if gradeId > 0 {
		bookDb = bookDb.Where("grade_id = ?", gradeId)
	}
	bookDb.Order("book_sort ASC").
		Order("book_id ASC").
		Scan(&bookArr)
	var bookNoArr []int64
	for _, bookV := range bookArr {
		bookNoArr = append(bookNoArr, bookV.BookNo)
	}

	var chapterArr []model.Chapter
	chapterByBookMap := map[int64][]model.Chapter{}
	c.Db.Where("book_no IN (?)", bookNoArr).
		Preload("ChapterKnowledge").
		Order("chapter_sort ASC").
		Order("chapter_id ASC").
		Find(&chapterArr)
	for _, chapterDbV := range chapterArr {
		chapterByBookMap[chapterDbV.BookNo] = append(chapterByBookMap[chapterDbV.BookNo], chapterDbV)
	}
	var chapterArrNew []model.Chapter
	for _, bookArrV := range bookArr {
		if _, ok := chapterByBookMap[bookArrV.BookNo]; ok {
			chapterArrNew = append(chapterArrNew, ChapterTree(chapterByBookMap[bookArrV.BookNo], 0)...)
		}
	}
	return chapterArrNew
}

type chapterArr struct {
	model.Chapter
	Children []chapterArr
}

func ChapterTree(chapter []model.Chapter, pid int64) []model.Chapter {
	task := []int64{}
	task = append(task, 0)
	res := []model.Chapter{}
	hasChild := false
	var parent int64 = pid
	for len(task) > 0 {
		hasChild = false
		for k, v := range chapter {
			if v.ChapterId == -1 {
				continue
			}
			if parent == v.ChapterParentId {
				res = append(res, v)
				task = append(task, v.ChapterId)
				chapter[k].ChapterId = -1
				hasChild = true
				parent = v.ChapterId
				break
			}
		}
		if !hasChild {
			end := len(task) - 1
			task = task[0:end] //将该数据删除
			//继续找它上级得其他子类
			if len(task) > 0 {
				end = end - 1
				parent = task[end]
			}

		}

	}
	return res
}
