package sync

import (
	"time"
)

type Synchronizer interface {
	SyncDB() error
}

func SetIntervelSyncDB(s Synchronizer, interval time.Duration) {

	for {
		s.SyncDB()
		time.Sleep(interval)
	}

}
