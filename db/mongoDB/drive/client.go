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
