package service

import (
	"context"
	"log"
	"os"

	model "github.com/mariajz/go-utils/mongodb-client-v2/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var mongoInstance *MongoDB

type MongoDB struct {
	DBClient *mongo.Client
}

func init() {
	mongoInstance = NewMongoDB()
}

func NewMongoDB() *MongoDB {
	clientOptions := options.Client().ApplyURI(os.Getenv("MONGO_URL"))
	var err error
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		panic(err)
	}
	err = client.Ping(context.Background(), readpref.Primary())
	if err != nil {
		log.Fatal("Error in db ping:", err)
	}
	return &MongoDB{
		DBClient: client,
	}
}

func (m *MongoDB) GetCollection(DbName string, collectionName string) *mongo.Collection {
	return mongoInstance.DBClient.Database(DbName).Collection(collectionName)
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

func (m *MongoDB) FindOneByID(ctx context.Context, collection *mongo.Collection, id interface{}, result interface{}, filter bson.M) error {
	if filter == nil && id != nil {
		filter = bson.M{"_id": id}
	}

	singleResult := collection.FindOne(ctx, filter)
	err := singleResult.Decode(result)
	return err
}

func (m *MongoDB) UpdateOne(ctx context.Context, collection *mongo.Collection, filter bson.M, update bson.M) (*model.UpdateResult, error) {
	updateResult, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Fatalf("Error updating document: %v", err)
		return nil, err
	}

	return &model.UpdateResult{
		ModifiedCount: updateResult.ModifiedCount,
	}, nil
}
