package db

import (
	"context"
	"testing"
	"time"

	"github.com/MarkLai0317/Advertising-CQRS/domain"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoIntegrationTestSuite struct {
	testMongoClient *mongo.Client
	suite.Suite
}

func TestMongoIntegrationTestSuite(t *testing.T) {
	suite.Run(t, &MongoIntegrationTestSuite{})
}

const dbName = "advertising"
const allAdCollection = "all_advertisement"
const activeAdCollection = "active_advertisement"

func (its *MongoIntegrationTestSuite) SetupSuite() {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://mark:markpwd@localhost:27017"))
	its.Assert().NoError(err)
	its.testMongoClient = client
}

func (its *MongoIntegrationTestSuite) TearDownSuite() {
	its.testMongoClient.Disconnect(context.Background())
}

func (its *MongoIntegrationTestSuite) SetupTest() {
	its.Assert().NoError(its.testMongoClient.Database(dbName).Collection(activeAdCollection).Drop(context.Background()))
	its.Assert().NoError(its.testMongoClient.Database(dbName).Collection(allAdCollection).Drop(context.Background()))
}

func (its *MongoIntegrationTestSuite) TestMongoCommandDB_Read() {
	// setup
	_, err := its.testMongoClient.Database(dbName).Collection(allAdCollection).InsertOne(context.Background(), bson.D{{"_id", "1"}, {"startAt", time.Now().Add(-time.Hour)}, {"endAt", time.Now().Add(time.Hour)}})
	its.Assert().NoError(err)
	_, err = its.testMongoClient.Database(dbName).Collection(allAdCollection).InsertOne(context.Background(), bson.D{{"_id", "2"}, {"startAt", time.Now().Add(-time.Hour)}, {"endAt", time.Now().Add(time.Hour)}})
	its.Assert().NoError(err)
	_, err = its.testMongoClient.Database(dbName).Collection(allAdCollection).InsertOne(context.Background(), bson.D{{"_id", "3"}, {"startAt", time.Now().Add(-time.Hour)}, {"endAt", time.Now().Add(time.Hour)}})
	its.Assert().NoError(err)
	_, err = its.testMongoClient.Database(dbName).Collection(allAdCollection).InsertOne(context.Background(), bson.D{{"_id", "4"}, {"startAt", time.Now().Add(time.Hour)}, {"endAt", time.Now().Add(2 * time.Hour)}})
	its.Assert().NoError(err)

	db, mongoClient, err := NewMongoCommandDB("mongodb://mark:markpwd@localhost:27017", dbName, allAdCollection)
	its.Assert().NoError(err)
	defer mongoClient.Disconnect(context.Background())
	// test
	adSlice, err := db.Read(context.Background())
	its.Assert().NoError(err)
	its.Assert().Len(adSlice, 3)
	its.Assert().Equal("1", adSlice[0].Id)
	its.Assert().Equal("2", adSlice[1].Id)
	its.Assert().Equal("3", adSlice[2].Id)
}

func (its *MongoIntegrationTestSuite) TestMongoQueryDB_Write() {
	// setup
	db, mongoClient, err := NewMongoQueryDB("mongodb://mark:markpwd@localhost:27017", dbName, activeAdCollection)
	its.Assert().NoError(err)
	defer mongoClient.Disconnect(context.Background())
	// test
	err = db.Write([]*domain.Advertisement{
		{Id: "1", StartAt: time.Now().Add(-time.Hour), EndAt: time.Now().Add(time.Hour)},
		{Id: "2", StartAt: time.Now().Add(-time.Hour), EndAt: time.Now().Add(time.Hour)},
		{Id: "3", StartAt: time.Now().Add(-time.Hour), EndAt: time.Now().Add(time.Hour)},
	})
	its.Assert().NoError(err)
	// verify
	cursor, err := its.testMongoClient.Database(dbName).Collection(activeAdCollection).Find(context.Background(), bson.D{})
	its.Assert().NoError(err)
	var adSlice []*domain.Advertisement
	err = cursor.All(context.Background(), &adSlice)
	its.Assert().NoError(err)
	its.Assert().Len(adSlice, 3)
	its.Assert().Equal("1", adSlice[0].Id)
	its.Assert().Equal("2", adSlice[1].Id)
	its.Assert().Equal("3", adSlice[2].Id)
}

func (its *MongoIntegrationTestSuite) TestMongoQueryDB_Read() {
	// setup
	_, err := its.testMongoClient.Database(dbName).Collection(activeAdCollection).InsertOne(context.Background(), bson.D{{"_id", "1"}, {"startAt", time.Now().Add(-time.Hour)}, {"endAt", time.Now().Add(time.Hour)}})
	its.Assert().NoError(err)
	_, err = its.testMongoClient.Database(dbName).Collection(activeAdCollection).InsertOne(context.Background(), bson.D{{"_id", "2"}, {"startAt", time.Now().Add(-time.Hour)}, {"endAt", time.Now().Add(time.Hour)}})
	its.Assert().NoError(err)
	_, err = its.testMongoClient.Database(dbName).Collection(activeAdCollection).InsertOne(context.Background(), bson.D{{"_id", "3"}, {"startAt", time.Now().Add(-time.Hour)}, {"endAt", time.Now().Add(time.Hour)}})
	its.Assert().NoError(err)

	db, mongoClient, err := NewMongoQueryDB("mongodb://mark:markpwd@localhost:27017", dbName, activeAdCollection)
	its.Assert().NoError(err)
	defer mongoClient.Disconnect(context.Background())
	// test
	adSlice, err := db.Read(context.Background())
	its.Assert().NoError(err)
	its.Assert().Len(adSlice, 3)
	its.Assert().Equal("1", adSlice[0].Id)
	its.Assert().Equal("2", adSlice[1].Id)
	its.Assert().Equal("3", adSlice[2].Id)
}
