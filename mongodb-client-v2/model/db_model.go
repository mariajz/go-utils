package model

type InsertOneResult struct {
	InsertedID interface{}
}

type UpdateResult struct {
	ModifiedCount int64 `bson:"n"`
}

type DeleteResult struct {
	DeletedCount int64 `bson:"n"`
}
