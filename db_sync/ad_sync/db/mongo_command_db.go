package db

import (
	"context"
	"fmt"
	"time"

	"github.com/MarkLai0317/Advertising-CQRS/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoCommandDB struct {
	mongoClient *mongo.Client
	collection  string
	dbName      string
}

func NewMongoCommandDB(uri, dbName, collection string) (*MongoCommandDB, *mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, nil, fmt.Errorf("error connecting to mongo: %w", err)

	}
	return &MongoCommandDB{mongoClient: client, dbName: dbName, collection: collection}, client, nil
}

func (db *MongoCommandDB) Read(parentCtx context.Context) ([]*domain.Advertisement, error) {
	ctx, cancel := context.WithTimeout(parentCtx, 10*time.Second)
	defer cancel()
	collection := db.mongoClient.Database("advertising").Collection(db.collection)

	now := time.Now()
	filter := bson.D{
		{"endAt", bson.D{{"$gt", now}}},
		{"startAt", bson.D{{"$lt", now}}},
	}

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var adSlice []*domain.Advertisement
	if err = cursor.All(ctx, &adSlice); err != nil {
		return nil, err
	}
	return adSlice, nil
}
