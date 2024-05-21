package db

import (
	"context"
	"reflect"

	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"

	"github.com/rickywei/sparrow/project/logger"
)

type MongoModel interface {
	GetId() string
}

func getCol[T MongoModel]() *mongo.Collection {
	var zero [0]T
	t := reflect.TypeOf(zero)
	for t.Kind() == reflect.Pointer || t.Kind() == reflect.Slice || t.Kind() == reflect.Array {
		t = t.Elem()
	}

	return mongoClient.Database(dbName).Collection(t.Name())
}

func toAnyList[T MongoModel](data []T) []any {
	return lo.Map(data, func(d T, _ int) any { return d })
}

func InsertMany[T MongoModel](ctx context.Context, data []T) (err error) {
	col := getCol[T]()
	if _, err = col.InsertMany(ctx, toAnyList(data)); err != nil {
		logger.L().Error("InsertMany failed", zap.Any("data", data), zap.Error(err))
	}

	return
}

func DeleteMany[T MongoModel](ctx context.Context, filter bson.D) {
	col := getCol[T]()
	col.DeleteOne(ctx, filter)
}

func UpdateMany[T MongoModel](ctx context.Context, data []T) {
	col := getCol[T]()
	sess, err := mongoClient.StartSession()
	if err != nil {
		logger.L().Error("start session failed", zap.Error(err))
	}
	defer sess.EndSession(ctx)
	_, err = sess.WithTransaction(ctx, func(ctx mongo.SessionContext) (_ any, err error) {
		for _, dt := range data {
			if _, err = col.UpdateByID(ctx, dt.GetId(), dt); err != nil {
				return
			}
		}
		return
	})
	if err != nil {
		logger.L().Error("UpdateMany failed", zap.Any("data", data), zap.Error(err))
	}

	return
}

func GetMany[T MongoModel](ctx context.Context, pageIndex, pageSize int64, filter, sorter, selector bson.D) (count int64, data []T, err error) {
	opt := &options.FindOptions{}
	opt.SetSkip((pageIndex - 1) * pageSize)
	opt.SetLimit(pageSize)
	if len(sorter) > 0 {
		opt.SetSort(sorter)
	}
	if len(selector) > 0 {
		opt.SetProjection(selector)
	}

	col := getCol[T]()
	count, err = col.CountDocuments(ctx, filter)
	if err != nil {
		logger.L().Error("GetMany count failed", zap.Error(err))
		return
	}
	cur, err := col.Find(ctx, filter, opt)
	if err != nil {
		logger.L().Error("GetMany find failed", zap.Error(err))
		return
	}
	err = cur.All(ctx, &data)
	if err != nil {
		logger.L().Error("GetMany cur failed", zap.Error(err))
	}

	return
}

func GetOne[T MongoModel](ctx context.Context, pageIndex, pageSize int64, filter, sorter, selector bson.D) (data T, err error) {
	opt := &options.FindOneOptions{}
	opt.SetSkip((pageIndex - 1) * pageSize)
	if len(sorter) > 0 {
		opt.SetSort(sorter)
	}
	if len(selector) > 0 {
		opt.SetProjection(selector)
	}

	col := getCol[T]()
	err = col.FindOne(ctx, filter, opt).Decode(&data)
	if err != nil {
		logger.L().Error("GetOne find failed", zap.Error(err))
		return
	}

	return
}
