package main

import (
	"log"

	"github.com/BurntSushi/toml"
	"github.com/kovetskiy/go-crontab"
	"github.com/theairkit/runcmd"
	"github.com/zazab/zhash"
)

const (
	feedsFile  = "feeds.conf"
	configFile = "sleipnir.conf"
)

func main() {

	feeds := Feeds{}
	if err := feeds.Load(feedsFile); err != nil {
		log.Fatal(err)
	}

	runner, err := runcmd.NewLocalRunner()
	if err != nil {
		log.Fatal(err)
	}

	config, err := getConfig(configFile)
	if err != nil {
		log.Fatal(err)
	}

	dbUrl, err := config.GetString("database", "url")
	if err != nil {
		log.Fatal(err)
	}

	dbName, err := config.GetString("database", "name")
	if err != nil {
		log.Fatal(err)
	}

	storage, err := NewMongoStorage(dbUrl, dbName)
	if err != nil {
		log.Fatal(err)
	}

	jobs := cron.Jobs{}
	for name, _ := range feeds {
		feed := feeds[name]

		log.Printf("[feed:%s] creating schedule '%s'", name, feed.Schedule)

		job, err := cron.NewJob(feed.Schedule, func() {
			feed.Run(runner, storage)
		})

		if err != nil {
			log.Fatal(err)
		}

		jobs = append(jobs, *job)
	}

	jobs.Process()
}

func getConfig(filepath string) (*zhash.Hash, error) {
	configData := map[string]interface{}{}
	if _, err := toml.DecodeFile(filepath, &configData); err != nil {
		return nil, err
	}

	hash := zhash.HashFromMap(configData)
	return &hash, nil
}
