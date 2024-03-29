package domain

import "time"

type Advertisement struct {
	Id         string     `bson:"_id"`
	Title      string     `bson:"title"`
	StartAt    time.Time  `bson:"startAt"`
	EndAt      time.Time  `bson:"endAt"`
	Conditions Conditions `bson:"conditions"`
}

type Conditions struct {
	AgeStart  int      `bson:"ageStart"`
	AgeEnd    int      `bson:"ageEnd"`
	Genders   []string `bson:"genders"`
	Countries []string `bson:"countries"`
	Platforms []string `bson:"platforms"`
}
