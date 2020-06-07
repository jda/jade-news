package main

import (
	"context"
	"net/url"
	"time"

	"github.com/golang/glog"
	"github.com/mmcdole/gofeed"
)

// Link represents a link to an article from a feed
type Link struct {
	Title     string
	URL       string
	Host      string
	Published time.Time
	FirstSeen time.Time
}

func fetchLinks(feeds []string, timeout time.Duration) (links []Link, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// TODO: do this in parallel?
	for _, feed := range feeds {
		feedLinks, err := getFeedLinks(ctx, feed)
		if err != nil {
			glog.Errorf("failed to process %s: %s", feed, err)
		}
		links = append(links, feedLinks...)
	}

	return links, nil
}

func getFeedLinks(ctx context.Context, feedURL string) (links []Link, err error) {
	runTimestamp := time.Now()

	parser := gofeed.NewParser()
	feed, err := parser.ParseURLWithContext(feedURL, ctx)
	if err != nil {
		return links, err
	}

	for _, item := range feed.Items {
		published := item.PublishedParsed
		if published == nil {
			published = item.UpdatedParsed
		}

		parsedLink, err := url.Parse(item.Link)
		if err != nil {
			glog.Warningf("could not parse link %s on feed %s: %s", item.Link, feedURL, err)
			continue
		}

		links = append(links, Link{
			Title:     item.Title,
			URL:       item.Link,
			Published: published.In(time.UTC),
			Host:      parsedLink.Host,
			FirstSeen: runTimestamp,
		})
	}

	return links, nil
}
