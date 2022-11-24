package mongoDB

import (
	"github.com/erdong01/kit/core"
	"go.mongodb.org/mongo-driver/mongo"
	"reflect"
	"strings"
)

type Database struct {
	*mongo.Database
}

func SetDatabase(cs interface{}, db interface{}, databaseArr ...string) {
	csValue := reflect.ValueOf(cs)
	csType := reflect.TypeOf(cs)
	if csType.Kind() == reflect.Ptr {
		csType = csType.Elem()
		csValue = csValue.Elem()
	}
	client := core.New().Mongo
	for i := 0; i < csType.NumField(); i++ {
		t := csType.Field(i)
		dbName := t.Tag.Get("db")
		if dbName == "" {
			dbName = t.Name
		}
		if len(databaseArr) > 0 {
			dbName = ""
			for _, databaseName := range databaseArr {
				if databaseName == dbName {
					dbName = databaseName
				}
			}
			if dbName == "" {
				continue
			}
		}
		if t.Type.String() == "*mongo.Database" {
			csValue.FieldByName(t.Name).
				Set(reflect.ValueOf(client.Database(dbName)))
		} else if strings.Contains(t.Type.String(), "Database") {
			dbValueJ := reflect.ValueOf(db).Elem()
			dbValue := reflect.New(dbValueJ.Type()).Elem()
			dbValue.Field(0).Set(reflect.ValueOf(client.Database(dbName)))
			csValue.FieldByName(t.Name).Set(dbValue.Addr())
		}
	}
}
