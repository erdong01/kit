package drive

import (
	"context"
	"github.com/erdong01/kit/config"
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
		SetMinPoolSize(32).
		SetMaxPoolSize(8192).
		SetConnectTimeout(time.Second * 65).
		SetMaxConnIdleTime(time.Minute * 5).
		SetSocketTimeout(time.Second * 60))
}

func Connect(opts ...*options.ClientOptions) *mongo.Client {
	ctx, _ := context.WithTimeout(context.Background(), 60*time.Second)
	client, err := mongo.Connect(ctx, opts...)
	if err != nil {
		log.Fatal(err)
	}
	return client
}
