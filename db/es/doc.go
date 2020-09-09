package es

import (
	"errors"
	"github.com/olivere/elastic"
)

// 文档结构
type DocParam struct {
	Id      string      // 文档id
	Body    interface{} // 文档内容
	Refresh string      // 强制刷新 默认为不强制 "true" 为单节点刷新 "wait_for"全节点刷新
}

// 根据id获取doc信息
func (e *es) Id(id string) (*elastic.GetResult, error) {
	return e.client.Get().Index(e.indexName).Type(e.typeName).Id(id).Do(todo())
}

// 根据id切片获取多个doc信息
func (e *es) Ids(ids []string) ([]*elastic.GetResult, error) {

	mGet := e.client.Mget()

	for _, id := range ids {
		mGet.Add(elastic.NewMultiGetItem().Index(e.indexName).Type(e.typeName).Id(id))
	}

	r, err := mGet.Do(todo())

	if err != nil {
		return nil, err
	}

	return r.Docs, nil
}

// 创建单个文档
func (e *es) Create(p *DocParam) (*elastic.IndexResponse, error) {
	create := e.client.Index().Index(e.indexName).Type(e.typeName).Id(p.Id)
	if str, ok := p.Body.(string); ok {
		create.BodyString(str)
	} else {
		create.BodyJson(p.Body)
	}

	switch p.Refresh {
	case "true", "wait_for":
		create.Refresh(p.Refresh)
	default:
		if "" != p.Refresh {
			return nil, errors.New("refresh 错误")
		}
	}

	return create.Do(todo())
}

// 修改单个文档
func (e *es) Update(p *DocParam) (*elastic.UpdateResponse, error) {
	update := e.client.Update().Index(e.indexName).Type(e.typeName).Id(p.Id)

	update.Doc(p.Body).DocAsUpsert(true)

	switch p.Refresh {
	case "true", "wait_for":
		update.Refresh(p.Refresh)
	default:
		if "" != p.Refresh {
			return nil, errors.New("refresh 错误")
		}
	}

	return update.Do(todo())
}

// 删除单个文档
func (e *es) Delete(p *DocParam) (*elastic.DeleteResponse, error) {
	del := e.client.Delete().Index(e.indexName).Type(e.typeName).Id(p.Id)
	switch p.Refresh {
	case "true", "wait_for":
		del.Refresh(p.Refresh)
	default:
		if "" != p.Refresh {
			return nil, errors.New("refresh 错误")
		}
	}

	return del.Do(todo())
}
