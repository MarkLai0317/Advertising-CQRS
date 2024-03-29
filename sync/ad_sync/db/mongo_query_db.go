package db

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/MarkLai0317/Advertising-CQRS/domain"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoQueryDB struct {
	mongoClient *mongo.Client
	collection  string
}

func NewMongoQueryDB(uri, collection string) (*MongoQueryDB, *mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	return &MongoQueryDB{mongoClient: client}, client, err
}

func (db *MongoQueryDB) Write(adSlice []*domain.Advertisement) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	collection := db.mongoClient.Database("advertising").Collection(db.collection)
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

func (db *MongoQueryDB) Read() ([]*domain.Advertisement, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	collection := db.mongoClient.Database("advertising").Collection(db.collection)

	cursor, err := collection.Find(ctx, nil)
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
