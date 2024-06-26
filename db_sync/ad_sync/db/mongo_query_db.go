package db

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/MarkLai0317/Advertising-CQRS/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoQueryDB struct {
	mongoClient *mongo.Client
	collection  string
	dbName      string
}

func NewMongoQueryDB(uri, dbName, collection string) (*MongoQueryDB, *mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, nil, fmt.Errorf("error connecting to mongo: %w", err)

	}
	return &MongoQueryDB{mongoClient: client, dbName: dbName, collection: collection}, client, nil
}

func (db *MongoQueryDB) Write(adSlice []*domain.Advertisement) error {
	if len(adSlice) == 0 {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	collection := db.mongoClient.Database(db.dbName).Collection(db.collection)
	adsInterface := make([]interface{}, len(adSlice))
	for i, v := range adSlice {
		adsInterface[i] = v
	}
	result, err := collection.InsertMany(ctx, adsInterface)
	if err != nil {
		return fmt.Errorf("error inserting advertisement: %w", err)
	}
	log.Printf("Inserted document with _id: %v\n", result.InsertedIDs...)
	return nil
}

func (db *MongoQueryDB) Read(parentCtx context.Context) ([]*domain.Advertisement, error) {
	ctx, cancel := context.WithTimeout(parentCtx, 10*time.Second)
	defer cancel()
	collection := db.mongoClient.Database(db.dbName).Collection(db.collection)

	cursor, err := collection.Find(ctx, bson.D{})
	if err != nil {
		return nil, fmt.Errorf("error finding advertisement: %w", err)
	}
	defer cursor.Close(ctx)
	var adSlice []*domain.Advertisement
	if err = cursor.All(ctx, &adSlice); err != nil {
		return nil, fmt.Errorf("error decoding advertisement: %w", err)
	}
	return adSlice, nil
}
