package drive

import (
	"context"
	"github.com/erDong01/micro-kit/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

func Init(Uri ...string) *mongo.Client {
	mongoConf := config.GetMongo()
	if len(Uri) == 0 {
		if err := config.New().Get(&mongoConf, "mongo"); err != nil {
			log.Fatal(err)
		}
	} else {
		mongoConf.Uri = Uri[0]
	}
	return Connect(options.Client().ApplyURI(mongoConf.Uri).
		SetMinPoolSize(10).
		SetMaxPoolSize(1000))
}

func Connect(opts ...*options.ClientOptions) *mongo.Client {
	ctx, _ := context.WithTimeout(context.Background(), 60*time.Second)
	client, err := mongo.Connect(ctx, opts...)
	if err != nil {
		log.Fatal(err)
	}
	return client
}
