package service

import (
	"context"
	"mongodb_client/constants"
	"mongodb_client/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DBClient *mongo.Client

func init() {
	clientOptions := options.Client().ApplyURI(constants.ConnectionString)
	var err error
	DBClient, err = mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		panic(err)
	}
}

func GetCollection(collectionName string) *mongo.Collection {
	return DBClient.Database(constants.DbName).Collection(collectionName)
}

func InsertOne(ctx context.Context, collection *mongo.Collection, document interface{}) (*model.InsertOneResult, error) {
	result, err := collection.InsertOne(ctx, document)
	if err != nil {
		return nil, err
	}
	return &model.InsertOneResult{
		InsertedID: result.InsertedID,
	}, nil
}

func FindOneByID(ctx context.Context, collection *mongo.Collection, id interface{}, result interface{}) error {
	filter := bson.M{"_id": id}
	singleResult := collection.FindOne(ctx, filter)
	err := singleResult.Decode(result)
	return err
}

func DeleteOneByID(ctx context.Context, collection *mongo.Collection, id interface{}) (*model.DeleteResult, error) {
	filter := bson.M{"_id": id}
	result, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		return nil, err
	}
	return &model.DeleteResult{
		DeletedCount: result.DeletedCount,
	}, nil
}
