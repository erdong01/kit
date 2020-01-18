package es

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/olivere/elastic"
)

// es实例
type es struct {
	indexName string
	typeName  string
	client    *elastic.Client
	mapping   string
}

// 配置文件
type Config struct {
	IndexName string
	TypeName  string
	Client    *elastic.Client
	Mapping   string
}

type JSON map[string]interface{}

func New(c *Config) (*es, error) {
	if c.IndexName == "" {
		return nil, errors.New("index 不能为空")
	}

	if c.TypeName == "" {
		return nil, errors.New("type 不能为空")
	}

	if c.Client == nil {
		return nil, errors.New("client 不能为空")
	}

	return &es{
		indexName: c.IndexName,
		typeName:  c.TypeName,
		client:    c.Client,
		mapping:   c.Mapping,
	}, nil
}

//获取最底层 elastic.Client
func (e *es) Client() *elastic.Client {
	return e.client
}

func (e *es) IndexName() string {
	return e.indexName
}

func (e *es) TypeName() string {
	return e.typeName
}

// 打印
func Dump(src interface{}) {
	data, err := json.MarshalIndent(src, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(data))
}

// 格式化打印sql代码
func DumpSql(query elastic.Query) {
	src, err := query.Source()
	if err != nil {
		panic(err)
	}
	data, err := json.MarshalIndent(src, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(data))
}

func todo() context.Context {
	return context.TODO()
}
