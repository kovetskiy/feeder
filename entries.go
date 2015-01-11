package main

import (
	"errors"
	"log"
)

type (
	Entries []*Entry

	Entry struct {
		Uid        string `json:"uid"`
		Url        string `json:"url"`
		Title      string `json:"title"`
		Image      string `json:"image,omitempty"`
		Preview    string `json:"preview"`
		CreateDate int64  `json:"create_date"`
	}

	EntryStorage interface {
		GetByUid(uid string) (*Entry, error)
		Add(entry *Entry, feedName string) error
	}
)

func (entry *Entry) Validate() error {
	switch {
	case entry.Uid == "":
		return errors.New("uid is a required field")
	case entry.Url == "":
		return errors.New("url is a required field")
	case entry.Title == "":
		return errors.New("title is a required field")
	case entry.Preview == "":
		return errors.New("preview is a required field")
	case entry.CreateDate == 0:
		return errors.New("create_date is a required field")
	}

	return nil
}

func (entries *Entries) Save(
	storage EntryStorage, feedName string,
) (saved int64, err error) {
	for _, entry := range *entries {
		log.Printf("[feed:%s] Validating entry %+v", feedName, entry)
		err = entry.Validate()
		if err != nil {
			return saved, err
		}

		log.Printf("[feed:%s] Looking up for entry with uid='%s'", feedName, entry.Uid)
		foundEntry, _ := storage.GetByUid(entry.Uid)
		if foundEntry != nil {
			// that is not error, that is normal situation
			// because storage may contains only unique entries
			continue
		}

		log.Printf("[feed:%s] Writing entry with uid='%s'", feedName, entry.Uid)
		err = storage.Add(entry, feedName)
		if err != nil {
			return saved, err
		}

		saved++
	}

	return saved, err
}
