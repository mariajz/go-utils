package service

import (
	"context"
	"log"
	"os"

	model "github.com/mariajz/go-utils/mongodb-client/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var DBClient *mongo.Client

func InitDB() {
	clientOptions := options.Client().ApplyURI(os.Getenv("MONGO_URL"))
	var err error
	DBClient, err = mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		panic(err)
	}
	err = DBClient.Ping(context.Background(), readpref.Primary())
	if err != nil {
		log.Fatal("Error in db ping:", err)
	}
}

func GetCollection(DbName string, collectionName string) *mongo.Collection {
	return DBClient.Database(DbName).Collection(collectionName)
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
