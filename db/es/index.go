package es

import "github.com/olivere/elastic"

// 判断索引是否存在
// bool true 存在 false 不存在
func (e *es) IndexExists() (bool, error) {
	return e.client.IndexExists(e.indexName).Do(todo())
}

// 创建索引
// body mappings 不传|string|json
// bool true 创建成功 false 创建失败
func (e *es) IndexCreate(mapping ...interface{}) (bool, error) {

	var body interface{}

	if 0 == len(mapping) {
		body = e.mapping
	} else {
		body = mapping[0]
	}

	c := e.client.CreateIndex(e.indexName)

	var r *elastic.IndicesCreateResult
	var err error

	if str, ok := body.(string); ok {
		r, err = c.BodyString(str).Do(todo())
	} else {
		r, err = c.BodyJson(body).Do(todo())
	}

	if err != nil {
		return false, err
	}

	return r.Acknowledged, err
}

// 删除索引
// bool true 创建成功 false 创建失败
func (e *es) IndexDelete() (bool, error) {
	r, err := e.client.DeleteIndex(e.indexName).Do(todo())
	return r.Acknowledged, err
}
