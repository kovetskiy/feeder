package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"

	"github.com/docopt/docopt-go"
	rss "github.com/jteeuwen/go-pkg-rss"
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
)

func main() {
	args, _ := getArgs()

	url := args["<url>"].(string)

	feed := rss.New(0, true, channelHandler, itemHandler)
	if err := feed.Fetch(url, nil); err != nil {
		fmt.Fprintf(os.Stderr, "%s", err)
		return
	}
}

func itemHandler(feed *rss.Feed, channel *rss.Channel, items []*rss.Item) {
	entries := Entries{}
	for _, item := range items {
		entry := Entry{
			Title:   item.Title,
			Preview: item.Description,
		}

		for _, link := range item.Links {
			entry.Url = link.Href
			break
		}

		pubDate, _ := item.ParsedPubDate()
		entry.CreateDate = int64(pubDate.Unix())

		hash := md5.Sum([]byte(entry.Title + entry.Url))
		entry.Uid = hex.EncodeToString(hash[:])

		entries = append(entries, &entry)
	}

	encoded, _ := json.Marshal(entries)
	fmt.Printf("%s", encoded)
}

func channelHandler(feed *rss.Feed, channels []*rss.Channel) {
}

func getArgs() (map[string]interface{}, error) {
	usage := `Sleipnir RSS 1.0
That's an application for configure your feeds.conf with rss feeds.
for example:
[kovetskiy]
command="./feeds/rss/rss https://github.com/kovetskiy.atom"
schedule="0 */5 * * *"

Usage:
	rss <url>
	rss --version | -v
	rss --help | -h

Options:
	-h --help     Show this screen.
	-v --version  Show version.
`
	return docopt.Parse(usage, nil, true, "Sleipnir RSS 1.0", false)
}
