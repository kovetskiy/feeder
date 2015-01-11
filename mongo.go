package main

import (
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

type (
	MongoStorage struct {
		entries *mgo.Collection
	}

	MongoDocument struct {
		Uid        string `bson:"uid"`
		Feed       string `bson:"feed"`
		Url        string `bson:"url"`
		Title      string `bson:"title"`
		Image      string `bson:"image,omitempty"`
		Preview    string `bson:"preview"`
		CreateDate int64  `bson:"create_date"`
	}
)

func NewMongoStorage(url string, database string) (*MongoStorage, error) {
	session, err := mgo.Dial(url)
	if err != nil {
		return nil, err
	}

	entries := session.DB(database).C("entries")

	storage := &MongoStorage{
		entries,
	}

	return storage, nil
}

func (storage *MongoStorage) GetByUid(uid string) (*Entry, error) {
	query := storage.entries.Find(bson.M{
		"uid": uid,
	})

	entity := MongoDocument{}
	err := query.One(&entity)
	if err != nil {
		return nil, err
	}

	if entity.Uid == "" {
		return nil, nil
	}

	//without Feed (name of feed)
	entry := &Entry{
		Uid:        entity.Uid,
		Url:        entity.Url,
		Title:      entity.Title,
		Image:      entity.Image,
		Preview:    entity.Preview,
		CreateDate: entity.CreateDate,
	}

	return entry, nil
}

func (storage *MongoStorage) Add(entry *Entry, feedName string) error {
	entity := &MongoDocument{
		Uid:        entry.Uid,
		Feed:       feedName,
		Url:        entry.Url,
		Title:      entry.Title,
		Image:      entry.Image,
		Preview:    entry.Preview,
		CreateDate: entry.CreateDate,
	}

	err := storage.entries.Insert(entity)

	return err
}
