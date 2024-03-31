package ad_sync

import (
	"context"
	"fmt"
	"log"

	"github.com/MarkLai0317/Advertising-CQRS/domain"
	"golang.org/x/sync/errgroup"
)

type CommandDB interface {
	Read(ctx context.Context) ([]*domain.Advertisement, error)
}

type QueryDB interface {
	Write([]*domain.Advertisement) error
	Read(ctx context.Context) ([]*domain.Advertisement, error)
}

type AdSynchronizer struct {
	commandDB        CommandDB
	queryDB          QueryDB
	existAdSet       map[string]*domain.Advertisement
	commandDBAdSlice []*domain.Advertisement
}

func NewAdSynchronizor(commandDB CommandDB, queryDB QueryDB) *AdSynchronizer {
	return &AdSynchronizer{
		commandDB: commandDB,
		queryDB:   queryDB,
	}
}

// func (s *AdSynchronizer) SyncDB() error {
// 	// read all data from QueryDB

// 	queryDBAdSlice, err := s.queryDB.Read()
// 	if err != nil {
// 		return fmt.Errorf("err reading QueryDB in SyncDB: %w", err)
// 	}
// 	s.ExistAdSet = make(map[string]*domain.Advertisement)
// 	for _, ad := range queryDBAdSlice {
// 		s.ExistAdSet[ad.Id] = ad
// 		log.Printf("ad id: %v\n", ad.Id)
// 	}

// 	// read new active data from CommandDBq
// 	commandDBAdSlice, err := s.commandDB.Read()
// 	if err != nil {
// 		return fmt.Errorf("err reading CommandDB in SyncDB: %w", err)
// 	}

// 	// store the new data that is not in QueryDB
// 	newAdSlice := make([]*domain.Advertisement, 0, 1000)
// 	for _, ad := range commandDBAdSlice {
// 		log.Printf("command ad id: %v\n", ad.Id)
// 		if _, ok := s.ExistAdSet[ad.Id]; !ok {
// 			newAdSlice = append(newAdSlice, ad)
// 			log.Printf("new ad id: %v\n", ad.Id)
// 		}
// 	}

// 	log.Printf("newAdSlice: %v\n", newAdSlice)
// 	// write the new data to QueryDB
// 	if err := s.queryDB.Write(newAdSlice); err != nil {
// 		return fmt.Errorf("err writing QueryDB in SyncDB: %w", err)
// 	}
// 	return nil
// }

func (s *AdSynchronizer) SyncDB() error {

	// read all data from QueryDB
	eg, ctx := errgroup.WithContext(context.Background())
	eg.Go(func() error {
		queryDBAdSlice, err := s.queryDB.Read(ctx)
		if err != nil {
			return fmt.Errorf("err reading QueryDB in SyncDB: %w", err)
		}
		s.existAdSet = make(map[string]*domain.Advertisement)
		for _, ad := range queryDBAdSlice {
			s.existAdSet[ad.Id] = ad
			//log.Printf("ad id: %v\n", ad.Id)
		}
		return nil
	})

	eg.Go(func() error {
		// read new active data from CommandDBq
		var err error
		s.commandDBAdSlice, err = s.commandDB.Read(ctx)
		if err != nil {
			return fmt.Errorf("err reading CommandDB in SyncDB: %w", err)
		}
		return nil
	})

	if err := eg.Wait(); err != nil {
		return fmt.Errorf("err getting data in SyncDB: %w", err)
	}

	// store the new data that is not in QueryDB
	newAdSlice := make([]*domain.Advertisement, 0, 1000)
	for _, ad := range s.commandDBAdSlice {
		//log.Printf("command ad id: %v\n", ad.Id)
		if _, ok := s.existAdSet[ad.Id]; !ok {
			newAdSlice = append(newAdSlice, ad)
			//log.Printf("new ad id: %v\n", ad.Id)
		}
	}

	log.Printf("newAdSlice: %v\n", len(newAdSlice))
	// write the new data to QueryDB
	if err := s.queryDB.Write(newAdSlice); err != nil {
		return fmt.Errorf("err writing QueryDB in SyncDB: %w", err)
	}

	return nil
}
