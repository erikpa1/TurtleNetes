package dbclient

import (
	"context"
	"time"
	"turtle/core/lgr"
	"turtle/core/serverKit"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Client wraps the MongoDB client with generic operations
type Client struct {
	client   *mongo.Client
	database *mongo.Database
	timeout  time.Duration
}

var MongoClient *Client

// NewClient creates a new MongoDB wrapper client
func InitMongoDb() {

	clientOptions := options.Client().ApplyURI(serverKit.SERVER_CONFIG.Mongo)

	client, err := mongo.Connect(context.Background(), clientOptions)

	if err != nil {
		lgr.ErrorStack("failed to connect to MongoDB: %w", err)
		return
	}

	// Ping the database to verify connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := client.Ping(ctx, nil); err != nil {
		lgr.ErrorStack("failed to ping MongoDB: %w", err)
		return
	}

	tmp := &Client{
		client:   client,
		database: client.Database(serverKit.SERVER_CONFIG.MongoDbName),
		timeout:  10 * time.Second,
	}

	MongoClient = tmp

}

func Insert(ctx context.Context, collection string, document interface{}) (*mongo.InsertOneResult, error) {

	tmp, err := MongoClient.database.Collection(collection).InsertOne(ctx, document)

	if err != nil {
		lgr.ErrorStack("failed to insert documents: %w", err)
	}

	return tmp, err

}

func FindOne[T any](ctx context.Context, collection string, filter interface{}) (*T, error) {
	var result T
	err := MongoClient.database.Collection(collection).FindOne(ctx, filter).Decode(&result)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			//lgr.Error("no document found in collection %s", collection)
			return nil, err
		}
		//lgr.Error("failed to find document: %w", err)
		return nil, err
	}

	return &result, nil
}

func FindMany[T any](ctx context.Context, collection string, filter interface{}, opts ...*options.FindOptions) ([]T, error) {
	results := []T{}

	cursor, err := MongoClient.database.Collection(collection).Find(ctx, filter, opts...)
	if err != nil {
		lgr.Error("failed to find documents: %w", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	err = cursor.All(ctx, &results)
	if err != nil {
		lgr.Error("failed to find documents: %w", err)
		return results, err
	}

	return results, nil
}

func SetById(ctx context.Context, collection string, id primitive.ObjectID, dataToSet interface{}) (*mongo.UpdateResult, error) {

	tmp, err := MongoClient.database.Collection(collection).UpdateByID(ctx, id, bson.M{"$set": dataToSet})

	if err != nil {
		lgr.Error("failed to update documents: %w", err)
		return tmp, err
	} else {
		return tmp, nil
	}

}

// IncrementBy increments fields in a MongoDB document by specified values
func IncrementBy(ctx context.Context, collection string, id primitive.ObjectID, dataToIncrement interface{}) (*mongo.UpdateResult, error) {

	result, err := MongoClient.database.Collection(collection).UpdateByID(ctx, id, bson.M{"$inc": dataToIncrement})

	if err != nil {
		lgr.Error("failed to increment document fields: %w", err)
		return result, err
	}

	return result, nil
}

// DeleteById deletes a document from MongoDB by its ID
func DeleteById(ctx context.Context, collection string, id primitive.ObjectID) (*mongo.DeleteResult, error) {

	result, err := MongoClient.database.Collection(collection).DeleteOne(ctx, bson.M{"_id": id})

	if err != nil {
		lgr.Error("failed to delete document: %w", err)
		return result, err
	}

	return result, nil
}
