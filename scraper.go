package main

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/soramon0/rssagg/internal/database"
)

func scrape(db *database.Queries, concurrency int, timeBetweenRequest time.Duration) {
	log.Printf("Scraping on %v goroutines every %s duration", concurrency, timeBetweenRequest)

	ticker := time.NewTicker(timeBetweenRequest)
	for ; ; <-ticker.C {
		feeds, err := db.GetNextFeedsToFetch(context.Background(), int32(concurrency))
		if err != nil {
			log.Println("error fetching feeds:", err)
			continue
		}

		wg := &sync.WaitGroup{}
		for _, f := range feeds {
			wg.Add(1)

			go func(wg *sync.WaitGroup, db *database.Queries, feed database.Feed) {
				defer wg.Done()

				_, err := db.MakrFeedAsFetched(context.Background(), feed.ID)
				if err != nil {
					log.Println("error marking feed as fetched:", err)
					return
				}

				rss, err := urlToFeed(feed.Url)
				if err != nil {
					log.Println("error fetching rss feed:", err)
					return
				}

				for _, item := range rss.Channel.Item {
					log.Println("Found post", item.Title, "on feed", feed.Name)
				}
				log.Printf("Feed %s collected, %v posts found\n", feed.Name, len(rss.Channel.Item))
			}(wg, db, f)
		}
		wg.Wait()
	}
}
