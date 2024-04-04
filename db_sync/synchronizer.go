package db_sync

import (
	"log"
	"time"
)

type Synchronizer interface {
	SyncDB() error
}

func SetIntervalSyncDB(s Synchronizer, interval time.Duration) {

	for {
		err := s.SyncDB()
		if err != nil {
			// send message to monitoring system
			// hear for simplicity just print the error
			log.Printf("error syncing db: %v\n", err)
		}
		time.Sleep(interval)
	}

}
