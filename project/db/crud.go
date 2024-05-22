package db

import (
	"context"
	"reflect"
	"time"

	"github.com/pkg/errors"
	"github.com/rickywei/sparrow/project/graph/model"
	"github.com/samber/lo"
	"github.com/spf13/cast"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoModel interface {
	GetMM() *model.MongoModel
	SetMM(*model.MongoModel)
}

func getCol[T MongoModel]() *mongo.Collection {
	var zero [0]T
	t := reflect.TypeOf(zero).Elem()
	for t.Kind() == reflect.Pointer || t.Kind() == reflect.Slice || t.Kind() == reflect.Array {
		t = t.Elem()
	}

	return mongoClient.Database(dbName).Collection(t.Name())
}

func toAnyList[T MongoModel](data []T) []any {
	return lo.Map(data, func(d T, _ int) any { return d })
}

func getIdFromMM(mm *model.MongoModel) string {
	if mm == nil || mm.ID == nil {
		return ""
	}
	return *mm.ID
}

func getUpdates[T MongoModel](data T) (updates bson.M, err error) {
	bs, err := bson.Marshal(data)
	if err != nil {
		err = errors.Wrap(err, "")
		return
	}
	err = bson.Unmarshal(bs, &updates)
	if err != nil {
		err = errors.Wrap(err, "")
		return
	}
	delete(updates, "model")
	return
}

func InsertMany[T MongoModel](ctx context.Context, data []T) (err error) {
	if len(data) <= 0 {
		return
	}
	col := getCol[T]()
	now := time.Now()
	mm := &model.MongoModel{
		CreatedAt: lo.ToPtr(now.Unix()),
		UpdatedAt: lo.ToPtr(now.Unix()),
	}
	for _, dt := range data {
		dt.SetMM(mm)
	}
	res, err := col.InsertMany(ctx, toAnyList(data))
	if err != nil {
		err = errors.Wrap(err, "")
		return
	}
	for i, id := range res.InsertedIDs {
		data[i].GetMM().ID = lo.ToPtr(cast.ToString(id))
	}

	return
}

func InsertOne[T MongoModel](ctx context.Context, data T) (err error) {
	if reflect.ValueOf(data).IsNil() {
		return
	}

	col := getCol[T]()
	now := time.Now()
	data.SetMM(&model.MongoModel{
		CreatedAt: lo.ToPtr(now.Unix()),
		UpdatedAt: lo.ToPtr(now.Unix()),
	})
	res, err := col.InsertOne(ctx, data)
	if err != nil {
		err = errors.Wrap(err, "")
		return
	}
	data.GetMM().ID = lo.ToPtr(cast.ToString(res.InsertedID))

	return
}

func DeleteMany[T MongoModel](ctx context.Context, filter bson.D) (err error) {
	col := getCol[T]()
	if _, err = col.DeleteMany(ctx, filter); err != nil {
		err = errors.Wrap(err, "")
	}
	return
}

func DeleteOne[T MongoModel](ctx context.Context, filter bson.D) (err error) {
	col := getCol[T]()
	if _, err = col.DeleteOne(ctx, filter); err != nil {
		err = errors.Wrap(err, "")
		return
	}
	return
}

func UpdateMany[T MongoModel](ctx context.Context, data []T) (err error) {
	col := getCol[T]()
	sess, err := mongoClient.StartSession()
	if err != nil {
		err = errors.Wrap(err, "")
	}
	defer sess.EndSession(ctx)
	_, err = sess.WithTransaction(ctx, func(ctx mongo.SessionContext) (_ any, err error) {
		for _, dt := range data {
			if _, err = col.UpdateByID(ctx, getIdFromMM(dt.GetMM()), dt); err != nil {
				return
			}
		}
		return
	})
	if err != nil {
		err = errors.Wrap(err, "")
		return
	}

	return
}

func UpdateOne[T MongoModel](ctx context.Context, data T) (err error) {
	col := getCol[T]()
	if _, err = col.UpdateByID(ctx, getIdFromMM(data.GetMM()), data); err != nil {
		err = errors.Wrap(err, "")
		return
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
		err = errors.Wrap(err, "")
		return
	}
	cur, err := col.Find(ctx, filter, opt)
	if err != nil {
		err = errors.Wrap(err, "")
		return
	}
	err = cur.All(ctx, &data)
	if err != nil {
		err = errors.Wrap(err, "")
		return
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
		err = errors.Wrap(err, "")
		return
	}

	return
}
