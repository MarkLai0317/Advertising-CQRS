package main

import (
	"log"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/MarkLai0317/Advertising-CQRS/db_sync"
	"github.com/MarkLai0317/Advertising-CQRS/db_sync/ad_sync"
	"github.com/MarkLai0317/Advertising-CQRS/db_sync/ad_sync/db"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}

	// get env for mongoRepo
	dbName := os.Getenv("DB_NAME")
	commandDBUrl := os.Getenv("COMMAND_DB_URL")
	commandDBColName := os.Getenv("COMMAND_DB_COL_NAME")
	commandDB, _, err := db.NewMongoCommandDB(commandDBUrl, dbName, commandDBColName)
	if err != nil {
		log.Fatalf("Error connecting to command db: %s", err)
	}
	queryDBUrl := os.Getenv("QUERY_DB_URL")
	queryDBColName := os.Getenv("QUERY_DB_COL_NAME")
	queryDB, _, err := db.NewMongoQueryDB(queryDBUrl, dbName, queryDBColName)
	if err != nil {
		log.Fatalf("Error connecting to query db: %s", err)
	}
	synchrorizer := ad_sync.NewAdSynchronizor(commandDB, queryDB)

	// sync db every X seconds
	secNum, err := strconv.Atoi(os.Getenv("INTERVAL_SYNC_DB"))
	if err != nil {
		log.Fatalf("Error decode INTERVAL_SYNC_DB: %s", err)
	}

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go db_sync.SetIntervalSyncDB(synchrorizer, time.Duration(secNum)*time.Second)
	wg.Wait()

	// commandclientDB.Disconnect(context.Background())
	// queryclientDB.Disconnect(context.Background())

}
