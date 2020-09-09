package es

import (
	"github.com/olivere/elastic"
	"testing"
)

var question *es

func TestMain(m *testing.M) {
	client, _ := elastic.NewClient(
		elastic.SetURL("http://127.0.0.1:9200", "http://127.0.0.1:9200", "http://127.0.0.1:9200"),
		elastic.SetSniff(false),
	)

	mapping := `{
		"settings":{
			"number_of_shards":1,
			"number_of_replicas":0
		},
		"mappings":{
			"tweet":{
				"properties":{
					"id":{
						"type":"keyword"
					},
					"name":{
						"type":"keyword"
					}
				}
			}
		}
	}`

	question, _ = New(&Config{
		"question2", "question2", client, mapping,
	})

	m.Run()

	//question.IndexDelete()
}

func TestCreateIndex(t *testing.T) {
	b, err := question.IndexExists()
	if err != nil {
		t.Log(err)
	}
	if b {
		// 删除
		question.IndexDelete()
	}
	r, err := question.IndexCreate()
	if err != nil {
		t.Log(err)
	}
	if !r {
		t.Log(r)
	}
}

func TestId(t *testing.T) {
	r, err := question.Id("454")
	//if err != nil {
	//	t.Log(err)
	//}
	if err == nil {

		if r.Found {
			s, err := r.Source.MarshalJSON()
			if err != nil {
				t.Log(err)
			}
			t.Log(string(s))
		}

	}
}

func TestIds(t *testing.T) {
	r, err := question.Ids([]string{"456", "454"})
	if err != nil {
		t.Log(err)
	}
	if 2 != len(r) {
		t.Log("数量错误")
	}
}

func TestCreate(t *testing.T) {
	_, err := question.Create(&DocParam{
		Id: "2",
		Body: JSON{
			"id":   "2",
			"name": "asdasd",
			//"age":  123,
		},
	})
	if err != nil {
		t.Log(err)
	}
}

func TestUpdate(t *testing.T) {
	_, err := question.Update(&DocParam{
		Id: "100",
		Body: JSON{
			"id":   "100",
			"name": "telangp",
		},
	})

	if err != nil {
		t.Log(err)
	}
}

func TestDelete(t *testing.T) {
	_, err := question.Delete(&DocParam{
		Id: "100",
	})
	if err != nil {
		t.Log(err)
	}
}

func TestSearch(t *testing.T) {

	sql := elastic.NewSearchSource()

	scope1 := elastic.NewBoolQuery()
	scope2 := elastic.NewBoolQuery()

	scope2.Must(elastic.NewTermQuery("id", "2"))

	a := []interface{}{"1", "2", "3"}
	scope1.Must(elastic.NewTermsQuery("question_no", a...))

	scope2.Filter(scope1)
	//
	sql.Query(scope2)
	//
	sql.Size(10)
	// 搜索字段
	sql.FetchSourceIncludeExclude([]string{"id"}, []string{})
	sql.From(1)

	DumpSql(sql)

	r, err := question.Search(sql)
	if err != nil {
		t.Log(err)
		Dump(r)
	}
}

func TestSearchRow(t *testing.T) {
	//
	sql1 := JSON{
		"_source": []string{"id"},
		"size":    1,
		"query": JSON{
			"bool": JSON{
				"must": JSON{
					"term": JSON{
						"id": 2,
					},
				},
			},
		},
	}
	r, err := question.SearchRow(sql1)
	if err != nil {
		t.Log(err)
		Dump(r)
	}

	sql2 := JSON{
		"_source": []string{"id"},
		"size":    1,
		"query": JSON{
			"bool": JSON{
				"must": []JSON{
					{
						"term": JSON{
							"id": 2,
						},
					},
				},
			},
		},
	}

	//
	r, err = question.SearchRow(sql2)
	if err != nil {
		t.Log(err)
		Dump(sql2)
	}

	sql3 := `
		{
			"_source" : ["id"],
			"query" : {
				"bool" : {
					"must" : [
						{
							"term": {
								"id" : 3
							}
						}
					]
				}
			}
		}
`

	r, err = question.SearchRow(sql3)
	if err != nil {
		t.Log(err)
	}

	var sql4 = struct {
		Source []string `json:"_source"`
		Size   int      `json:"size"`
	}{
		Source: []string{"id"},
		Size:   1,
	}

	r, err = question.SearchRow(sql4)
	if err != nil {
		t.Log(err)
	}
	Dump(r)

}
