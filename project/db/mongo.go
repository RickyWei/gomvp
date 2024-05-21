package db

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.uber.org/zap"

	"github.com/rickywei/sparrow/project/conf"
	"github.com/rickywei/sparrow/project/logger"
)

var (
	mongoClient *mongo.Client
)

func init() {
	var err error

	ctx := context.Background()

	mongoClient, err = mongo.Connect(ctx, (&options.ClientOptions{}).
		ApplyURI(fmt.Sprintf("mongodb://%s:%s@%s:%s/?connect=direct",
			conf.String("mongo.user"), conf.String("mongo.password"), conf.String("mongo.ip"), conf.String("mongo.port"))).
		SetMaxPoolSize(5))
	if err != nil {
		logger.L().Fatal("mongo init failed", zap.Error(err))
	}

	if err = mongoClient.Ping(ctx, readpref.Primary()); err != nil {
		logger.L().Fatal("mongo ping failed", zap.Error(err))
	}

	if err = createIndexes(); err != nil {
		logger.L().Fatal("mongo indexes failed", zap.Error(err))
	}
}

func createIndexes() (err error) {
	// if _, err = getCol[]().Indexes().CreateOne(ctx,
	// 	mongo.IndexModel{Keys: bson.D{{Key: "name", Value: 1}}, Options: options.Index().SetUnique(true)}); err != nil {
	// 	return
	// }

	return
}
