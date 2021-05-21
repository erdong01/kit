package mongoDB

import (
	"github.com/erDong01/micro-kit/core"
	"reflect"
)

func SetDatabase(d interface{}) {
	refValue := reflect.ValueOf(d)
	rType := reflect.TypeOf(d)
	client := core.New().Mongo
	if rType.Kind() == reflect.Ptr {
		rType = rType.Elem()
		refValue = refValue.Elem()
	}
	for i := 0; i < rType.NumField(); i++ {
		t := rType.Field(i)

		if t.Type.String() == "*mongo.Database" {
			refValue.FieldByName(t.Name).
				Set(reflect.ValueOf(client.Database(t.Name)))
		}
	}
}
