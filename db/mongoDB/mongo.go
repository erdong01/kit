package mongoDB

import (
	"fmt"
	"github.com/erDong01/micro-kit/core"
	"go.mongodb.org/mongo-driver/mongo"
	"reflect"
)

type Database struct {
	*mongo.Database
}

func SetDatabase(d interface{}, db interface{}) {
	refValue := reflect.ValueOf(d)
	rType := reflect.TypeOf(d)
	client := core.New().Mongo

	if rType.Kind() == reflect.Ptr {
		rType = rType.Elem()
		refValue = refValue.Elem()
		fmt.Println(refValue, 1111)
	}
	for i := 0; i < rType.NumField(); i++ {
		t := rType.Field(i)
		fmt.Println(t.Type.String())
		if t.Type.String() == "*mongo.Database" {
			fmt.Println(t.Name)
			fmt.Println(refValue.FieldByName(t.Name), 88888888)
			refValue.FieldByName(t.Name).
				Set(reflect.ValueOf(client.Database(t.Name)))
		}
		if t.Type.String() == "*dbo.Database" {
			dbValue := reflect.ValueOf(db)
			fmt.Println(t.Name, 99999999)
			dbValue.Elem().Field(0).Set(reflect.ValueOf(client.Database(t.Name)))
			refValue.FieldByName(t.Name).Set(dbValue)
			//refSon := reflect.TypeOf(refValue.FieldByName(t.Name).Elem())

			//for i := 0; i < refSon.NumField(); i++ {
			//	field := refSon.Field(i)
			//	fmt.Printf("name:%s index:%d type:%v json tag:%v\n", field.Name, field.Index, field.Type, field.Tag.Get("json"))
			//}
			//refSon.Field(0).Set(reflect.ValueOf(client.Database(t.Name)))
		}
	}
}
