package main

import (
	"encoding/json"
	"log"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/theairkit/runcmd"
)

type (
	Feeds map[string]*Feed

	Feed struct {
		Command  string `toml:"command"`
		Schedule string `toml:"schedule"`
		Name     string
	}
)

func (feeds *Feeds) Load(filepath string) error {
	log.Printf("toml decoding file '%s'", filepath)

	_, err := toml.DecodeFile(filepath, &feeds)
	if err != nil {
		return err
	}

	for name, feed := range *feeds {
		feed.Name = name
	}

	log.Printf("Decoded feeds: %+v", *feeds)

	return nil
}

func (feed *Feed) Run(runner runcmd.Runner, storage EntryStorage) {
	log.Printf("[feed:%s] Execute '%s'", feed.Name, feed.Command)

	cmd, _ := runner.Command(feed.Command)
	response, err := cmd.Run()
	if err != nil {
		log.Printf("[feed:%s] Got error during execution '%s': %s",
			feed.Name, feed.Command, err.Error())
		return
	}

	entries := Entries{}
	err = json.Unmarshal([]byte(strings.Join(response, "\n")), &entries)
	if err != nil {
		log.Printf(
			"[feed:%s] Got error during decoding json response '%s': %s",
			feed.Name, response, err.Error())
		return
	}

	log.Printf("[feed:%s] Saving entries")

	saved, err := entries.Save(storage, feed.Name)
	if err != nil {
		log.Printf("[feed:%s] Got error during saving entries: %s",
			feed.Name, err.Error())
		return
	}

	log.Printf("[feed:%s] Saved %d entries", feed.Name, saved)
}
