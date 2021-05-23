// database pool
package dbo

import (
	"context"
	"github.com/erDong01/micro-kit/core"
	"github.com/erDong01/micro-kit/db/mongoDB"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"reflect"
	"time"
)

var MgoClient *mongodb

type mongodb struct {
	*core.Core
	client   *mongo.Client
	MasterDB *Database
	TestDB   *Database
}

type Database struct {
	*mongo.Database
}

func Init(dbUrl string) {
	MgoClient = &mongodb{}
	MgoClient = NewMgo(dbUrl)
	MgoClient.MongoRegister()
	mongoDB.SetDatabase(MgoClient, &Database{})

}

func NewMgo(urlStr string) *mongodb {
	var mgo mongodb
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(urlStr)) //mongodb://127.0.0.1/test
	if err != nil {
		panic(err)
	}
	mgo.client = client
	return &mgo
}

func (this *Database) FindOne(table string, q interface{}, result interface{}, orther ...interface{}) (err error) {
	var opts []*options.FindOneOptions
	for _, opt := range orther {
		one := options.FindOne()
		for k, v := range opt.(bson.M) {
			if k == "Sort" {
				one.Sort = v
			}
		}
		opts = append(opts, one)
	}

	var data primitive.D
	if err = this.Collection(table).FindOne(context.Background(), q, opts...).Decode(&data); err != nil {
		return
	}

	*result.(*map[string]interface{}) = this.Map(data)
	return
}

func (this *Database) Count(table string, q interface{}) (int64, error) {
	count, err := this.Collection(table).CountDocuments(context.Background(), q)
	return count, err
}
func (this *Database) Find(table string, q interface{}, result interface{}, orther ...interface{}) (err error) {
	var opts []*options.FindOptions
	var collOpts []*options.CollectionOptions
	for _, opt := range orther {
		one := options.Find()
		for k, v := range opt.(bson.M) {
			if k == "Sort" {
				one.Sort = v
			} else if k == "Limit" {
				//i := t.Int64(v)
				i := v.(int64)
				one.Limit = &i
			} else if k == "Skip" {
				//i := t.Int64(v)
				i := v.(int64)
				one.Skip = &i
			} else if k == "ReadPreference" {
				var ref *readpref.ReadPref
				if v == "SecondaryPreferred" {
					ref = readpref.SecondaryPreferred()
				} else if v == "Secondary" {
					ref = readpref.Secondary()
				}
				collOpts = []*options.CollectionOptions{&options.CollectionOptions{ReadPreference: ref}}
			}
		}
		opts = append(opts, one)
	}

	cursor, err := this.Collection(table, collOpts...).Find(context.Background(), q, opts...)
	if err != nil {
		return
	}
	if err = cursor.Err(); err != nil {
		return
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		//data := bson.D{}
		var data interface{}
		if err = cursor.Decode(&data); err != nil {
			return
		}
		*result.(*[]map[string]interface{}) = append(*result.(*[]map[string]interface{}), this.Map(data))
	}
	return
}

func (this *Database) Insert(table string, data interface{}) (string, error) {

	if result, err := this.Collection(table).InsertOne(context.Background(), data); err != nil {
		return "", nil
	} else {
		return result.InsertedID.(primitive.ObjectID).Hex(), err
	}
}

func (this *Database) Insert2(session *mongo.SessionContext, table string, data interface{}) (string, error) {
	if result, err := this.Collection(table).InsertOne(*session, data); err != nil {
		return "", nil
	} else {
		return result.InsertedID.(primitive.ObjectID).Hex(), err
	}
}

func (this *Database) Update(table string, q interface{}, update interface{}, orther ...interface{}) (*mongo.UpdateResult, error) {
	var opts []*options.UpdateOptions
	for _, opt := range orther {
		one := options.Update()
		for k, v := range opt.(bson.M) {
			if k == "Upsert" {
				i := v.(bool)
				one.Upsert = &i
			}
		}
		opts = append(opts, one)
	}
	if result, err := this.Collection(table).UpdateOne(context.Background(), q, update, opts...); err != nil {
		return nil, err
	} else {
		return result, nil
	}

}

func (this *Database) UpdateMany(table string, q interface{}, update interface{}, orther ...interface{}) error {
	var opts []*options.UpdateOptions
	for _, opt := range orther {
		one := options.Update()
		for k, v := range opt.(bson.M) {
			if k == "Upsert" {
				i := v.(bool)
				one.Upsert = &i
			}
		}
		opts = append(opts, one)
	}
	_, err := this.Collection(table).UpdateMany(context.Background(), q, update, opts...)
	return err
}

func (this *Database) FindOneAndUpdate(table string, q interface{}, update interface{}, result interface{}, orther ...interface{}) (err error) {
	var opts []*options.FindOneAndUpdateOptions
	for _, opt := range orther {
		one := options.FindOneAndUpdate()
		for k, v := range opt.(bson.M) {
			if k == "Sort" {
				one.Sort = v
			} else if k == "Upsert" {
				i := v.(bool)
				one.Upsert = &i
			} else if k == "New" {
				i := v.(bool)
				if i == true {
					one.SetReturnDocument(options.After)
				} else if i == false {
					one.SetReturnDocument(options.Before)
				}
			}
		}
		opts = append(opts, one)
	}
	var data interface{}
	if err = this.Collection(table).FindOneAndUpdate(context.Background(), q, update, opts...).Decode(&data); err != nil {
		return
	}
	*result.(*map[string]interface{}) = this.Map(data)
	return
}

//启动事务会话
func (this *mongodb) StartSession(callback func(mongo.SessionContext) error) error {
	ctx := context.Background()
	return this.Core.Mongo.UseSession(ctx, func(session mongo.SessionContext) (err error) {
		if err = session.StartTransaction(); err != nil {
			return
		}
		return callback(session)
	})
}

func (this *Database) DeleteOne(table string, q interface{}, ctx ...context.Context) error {
	c := context.Background()
	if ctx != nil {
		c = ctx[0]
	}
	result, err := this.Collection(table).DeleteOne(c, q)
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return err
	}
	return nil
}

func (this *Database) DeleteMany(table string, q interface{}) error {
	result, err := this.Collection(table).DeleteMany(context.Background(), q)
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return err
	}
	return nil
}
func (this *Database) Array(result interface{}) []interface{} {
	resultArray := result.(primitive.A)
	for k, v := range resultArray {
		if v == nil {
			continue
		}
		switch reflect.TypeOf(v).String() {
		case "primitive.D":
			resultArray[k] = this.Map(v)
		case "primitive.ObjectID":
			resultArray[k] = v.(primitive.ObjectID).Hex()
		case "primitive.A":
			resultArray[k] = this.Array(v)
		}
	}
	return resultArray
}

func (this *Database) Map(result interface{}) map[string]interface{} {
	resultMap := result.(primitive.D).Map()
	for k, v := range resultMap {
		if v == nil {
			continue
		}
		switch reflect.TypeOf(v).String() {
		case "primitive.D":
			resultMap[k] = this.Map(v)
		case "primitive.ObjectID":
			resultMap[k] = v.(primitive.ObjectID).Hex()
		case "primitive.A":
			resultMap[k] = this.Array(v)
		}
	}
	return resultMap
}
