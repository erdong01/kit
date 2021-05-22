package mongoDB

import (
	"github.com/erDong01/micro-kit/core"
	"go.mongodb.org/mongo-driver/mongo"
	"reflect"
)

type Database struct {
	*mongo.Database
}

func SetDatabase(cs interface{}, db interface{}) {
	csValue := reflect.ValueOf(cs)
	csType := reflect.TypeOf(cs)
	client := core.New().Mongo
	if csType.Kind() == reflect.Ptr {
		csType = csType.Elem()
		csValue = csValue.Elem()
	}
	for i := 0; i < csType.NumField(); i++ {
		t := csType.Field(i)
		if t.Type.String() == "*mongo.Database" {
			csValue.FieldByName(t.Name).
				Set(reflect.ValueOf(client.Database(t.Name)))
		}
		if t.Type.String() == "*dbo.Database" {
			dbValue := reflect.ValueOf(db)
			dbValue.Elem().Field(0).Set(reflect.ValueOf(client.Database(t.Name)))
			csValue.FieldByName(t.Name).Set(dbValue)
		}
	}
}
