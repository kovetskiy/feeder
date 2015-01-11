package main

import (
	"log"

	"github.com/kovetskiy/go-crontab"
	"github.com/theairkit/runcmd"
)

const (
	feedsConfig = "feeds.conf"
)

func main() {
	feeds := Feeds{}
	if err := feeds.Load(feedsConfig); err != nil {
		log.Fatal(err)
	}

	runner, err := runcmd.NewLocalRunner()
	if err != nil {
		log.Fatal(err)
	}

	storage, err := NewMongoStorage("localhost", "feeder")
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
