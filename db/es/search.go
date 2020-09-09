package es

import (
	"github.com/olivere/elastic"
)

// 通过查询构建器查询
func (e *es) Search(sql *elastic.SearchSource) (*elastic.SearchResult, error) {
	return e.client.Search().Index(e.indexName).Type(e.typeName).SearchSource(sql).Do(todo())
}

// 通过原始语句查询
func (e *es) SearchRow(sql interface{}) (*elastic.SearchResult, error) {
	return e.client.Search().Index(e.indexName).Type(e.typeName).Source(sql).Do(todo())
}
