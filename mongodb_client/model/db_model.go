package model

type InsertOneResult struct {
	InsertedID interface{}
}

type DeleteResult struct {
	DeletedCount int64 `bson:"n"`
}
