package db

import "go.mongodb.org/mongo-driver/bson"

const (
	dbName     = "sparrow"
	colSparrow = "sparrow"
)

var (
	SortByCreatedAt     = bson.D{{Key: "created_at", Value: 1}}
	SortByCreatedAtDesc = bson.D{{Key: "created_at", Value: -1}}
	SortByUpdatedAt     = bson.D{{Key: "updated_at", Value: 1}}
	SortByUpdatedAtDesc = bson.D{{Key: "updated_at", Value: -1}}
)
