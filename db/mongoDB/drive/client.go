package drive

import (
	"context"
	"github.com/erDong01/micro-kit/config"
	"github.com/erDong01/micro-kit/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

func Init() *mongo.Client {
	mongoConf := config.GetMongo()
	if err := config.New().Get(&mongoConf, "mongo"); err != nil {
		log.Fatal(err)
	}
	options.Client().ApplyURI(mongoConf.Uri)
	return Connect(options.Client().ApplyURI(mongoConf.Uri))
}

func Connect(opts ...*options.ClientOptions) *mongo.Client {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, opts...)
	if err != nil {
		log.Fatal(err)
	}
	return client
}
