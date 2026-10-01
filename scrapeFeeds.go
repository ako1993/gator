package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/ako1993/internal/database"
	"github.com/google/uuid"
)

func scrapeFeeds(s *state) error {
	feed, err := s.dbQueries.GetNextFeedToFetch(context.Background())
	if err != nil {
		fmt.Printf("ERROR GETTING NEXT FEED TO FETCH:%v", err)
		return err
	}
	err = s.dbQueries.MarkFeedFetched(context.Background(), feed.ID)
	if err != nil {
		fmt.Printf("ERROR MARKING FEED FETCHED:%v", err)
		return err
	}
	rss_feed, err := fetchFeed(context.Background(), feed.Url)
	if err != nil {
		fmt.Printf("ERROR FETCHING FEED:%v", err)
		return err
	}
	for _, item := range rss_feed.Channel.Item {
		fmt.Println(item.Title)
		parsedTime, err := time.Parse(time.RFC1123, item.PubDate)
		if err != nil {
			fmt.Printf("ERROR PARSING PUBLISHED TIME:%v", err)
			return err
		}
		save_post_params := database.CreatePostParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Title: sql.NullString{String: item.Title,
				Valid: item.Title != "",
			},
			Url: sql.NullString{String: feed.Url,
				Valid: feed.Url != "",
			},
			Description: sql.NullString{String: item.Description,
				Valid: item.Description != "",
			},
			PublishedAt: sql.NullTime{Time: parsedTime,
				Valid: item.PubDate != "",
			},
			FeedID: uuid.NullUUID{UUID: feed.ID,
				Valid: true,
			},
		}
		_, err = s.dbQueries.CreatePost(context.Background(), save_post_params)
		if err != nil {
			fmt.Printf("ERROR CREATING POST:%v", err)
		}
	}
	return nil
}
