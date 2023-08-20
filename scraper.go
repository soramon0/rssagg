package main

import (
	"context"
	"database/sql"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
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
					description := sql.NullString{}
					if item.Description != "" {
						description.String = item.Description
						description.Valid = true
					}

					pubDate, err := time.Parse(time.RFC1123Z, item.PubDate)
					if err != nil {
						log.Printf("couldn't parse date %v with err %v\n", item.PubDate, err)
						continue
					}

					_, err = db.CreatePost(context.Background(), database.CreatePostParams{
						ID:          uuid.New(),
						Url:         item.Link,
						Title:       item.Title,
						Description: description,
						PublishedAt: pubDate,
						FeedID:      feed.ID,
						CreatedAt:   time.Now().UTC(),
						UpdatedAt:   time.Now().UTC(),
					})
					if err != nil {
						if strings.Contains(err.Error(), `violates unique constraint "posts_url_key"`) {
							continue
						}
						log.Println("failed to create post for:", item.Title, item.Link)
						log.Println(err)
					}
				}

				log.Printf("Feed %s collected, %v posts found\n", feed.Name, len(rss.Channel.Item))
			}(wg, db, f)
		}
		wg.Wait()
	}
}
