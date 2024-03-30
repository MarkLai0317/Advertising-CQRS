package ad_sync

import (
	"context"
	"testing"
	"time"

	"github.com/MarkLai0317/Advertising-CQRS/db_sync/ad_sync/db"
	"github.com/MarkLai0317/Advertising-CQRS/domain"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type AdSynchronizerIntegrationTestSuite struct {
	testMongoClient *mongo.Client
	suite.Suite
}

func TestAdSynchronizerIntegrationTestSuite(t *testing.T) {
	suite.Run(t, &AdSynchronizerIntegrationTestSuite{})
}

func (its *AdSynchronizerIntegrationTestSuite) SetupSuite() {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://mark:markpwd@localhost:27017"))
	its.Assert().NoError(err)
	its.testMongoClient = client
}

func (its *AdSynchronizerIntegrationTestSuite) TearDownSuite() {
	its.testMongoClient.Disconnect(context.Background())
}

func (its *AdSynchronizerIntegrationTestSuite) SetupTest() {
	its.Require().NoError(its.testMongoClient.Database("advertising").Collection("active_advertisement").Drop(context.Background()))
	its.Require().NoError(its.testMongoClient.Database("advertising").Collection("all_advertisement").Drop(context.Background()))
}

func (its *AdSynchronizerIntegrationTestSuite) TestAdSynchronizer_SyncDB_active_db_empty() {
	// setup
	_, err := its.testMongoClient.Database("advertising").Collection("all_advertisement").InsertOne(context.Background(), bson.D{{"_id", "1"}, {"startAt", time.Now().Add(-time.Hour)}, {"endAt", time.Now().Add(time.Hour)}})
	its.Require().NoError(err)
	_, err = its.testMongoClient.Database("advertising").Collection("all_advertisement").InsertOne(context.Background(), bson.D{{"_id", "2"}, {"startAt", time.Now().Add(-time.Hour)}, {"endAt", time.Now().Add(time.Hour)}})
	its.Require().NoError(err)
	_, err = its.testMongoClient.Database("advertising").Collection("all_advertisement").InsertOne(context.Background(), bson.D{{"_id", "3"}, {"startAt", time.Now().Add(-time.Hour)}, {"endAt", time.Now().Add(time.Hour)}})
	its.Require().NoError(err)
	_, err = its.testMongoClient.Database("advertising").Collection("all_advertisement").InsertOne(context.Background(), bson.D{{"_id", "4"}, {"startAt", time.Now().Add(time.Hour)}, {"endAt", time.Now().Add(2 * time.Hour)}})
	its.Require().NoError(err)

	commandDB, mongoClient, err := db.NewMongoCommandDB("mongodb://mark:markpwd@localhost:27017", "all_advertisement")
	its.Assert().NoError(err)
	defer mongoClient.Disconnect(context.Background())
	queryDB, mongoClient, err := db.NewMongoQueryDB("mongodb://mark:markpwd@localhost:27017", "active_advertisement")
	its.Assert().NoError(err)
	defer mongoClient.Disconnect(context.Background())

	// test
	adSynchronizer := NewAdSynchronizor(commandDB, queryDB)
	adSynchronizer.SyncDB()

	//verify
	cursor, err := its.testMongoClient.Database("advertising").Collection("active_advertisement").Find(context.Background(), bson.D{})
	its.Assert().NoError(err)
	var adSlice []*domain.Advertisement
	err = cursor.All(context.Background(), &adSlice)
	its.Assert().NoError(err)
	its.Assert().Len(adSlice, 3)
	its.Assert().Equal("1", adSlice[0].Id)
	its.Assert().Equal("2", adSlice[1].Id)
	its.Assert().Equal("3", adSlice[2].Id)
}

func (its *AdSynchronizerIntegrationTestSuite) TestAdSynchronizer_SyncDB_active_db_exists_data() {
	// setup all_advertisement
	all_collection := "all_advertisement"
	_, err := its.testMongoClient.Database("advertising").Collection(all_collection).InsertOne(context.Background(), bson.D{{"_id", "1"}, {"startAt", time.Now().Add(-time.Hour)}, {"endAt", time.Now().Add(time.Hour)}})
	its.Require().NoError(err)
	_, err = its.testMongoClient.Database("advertising").Collection(all_collection).InsertOne(context.Background(), bson.D{{"_id", "2"}, {"startAt", time.Now().Add(-time.Hour)}, {"endAt", time.Now().Add(time.Hour)}})
	its.Require().NoError(err)
	_, err = its.testMongoClient.Database("advertising").Collection(all_collection).InsertOne(context.Background(), bson.D{{"_id", "3"}, {"startAt", time.Now().Add(-time.Hour)}, {"endAt", time.Now().Add(time.Hour)}})
	its.Require().NoError(err)
	_, err = its.testMongoClient.Database("advertising").Collection(all_collection).InsertOne(context.Background(), bson.D{{"_id", "4"}, {"startAt", time.Now().Add(time.Hour)}, {"endAt", time.Now().Add(2 * time.Hour)}})
	its.Require().NoError(err)

	// setup active_advertisement
	active_collection := "active_advertisement"
	_, err = its.testMongoClient.Database("advertising").Collection(active_collection).InsertOne(context.Background(), bson.D{{"_id", "1"}, {"startAt", time.Now().Add(-time.Hour)}, {"endAt", time.Now().Add(time.Hour)}})
	its.Require().NoError(err)

	commandDB, mongoClient, err := db.NewMongoCommandDB("mongodb://mark:markpwd@localhost:27017", "all_advertisement")
	its.Assert().NoError(err)
	defer mongoClient.Disconnect(context.Background())
	queryDB, mongoClient, err := db.NewMongoQueryDB("mongodb://mark:markpwd@localhost:27017", "active_advertisement")
	its.Assert().NoError(err)
	defer mongoClient.Disconnect(context.Background())

	// test
	adSynchronizer := NewAdSynchronizor(commandDB, queryDB)
	adSynchronizer.SyncDB()

	//verify
	cursor, err := its.testMongoClient.Database("advertising").Collection("active_advertisement").Find(context.Background(), bson.D{})
	its.Assert().NoError(err)
	var adSlice []*domain.Advertisement
	err = cursor.All(context.Background(), &adSlice)
	its.Assert().NoError(err)
	its.Assert().Len(adSlice, 3)
	its.Assert().Equal("1", adSlice[0].Id)
	its.Assert().Equal("2", adSlice[1].Id)
	its.Assert().Equal("3", adSlice[2].Id)
}
