package main

import (
	"context"
	"fmt"
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
	}
	return nil
}
