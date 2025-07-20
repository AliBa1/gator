package main

import (
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
)

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func (f *RSSFeed) Print() {
	fmt.Printf("%s: %s\n", f.Channel.Title, f.Channel.Link)
	fmt.Printf("Description: %s\n", f.Channel.Description)
	fmt.Printf("Items:\n")
	for _, item := range f.Channel.Item {
		fmt.Printf("%s (%s)\n", item.Title, item.Link)
		fmt.Printf("%s\n", item.PubDate)
		fmt.Printf("%s\n", item.Description)
	}
	fmt.Println()
}

func fetchFeed(context context.Context, feedURL string) (*RSSFeed, error) {
	request, err := http.NewRequestWithContext(context, "GET", feedURL, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("User-Agent", "gator")

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}

	rawData, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var feed *RSSFeed
	err = xml.Unmarshal(rawData, &feed)
	if err != nil {
		return nil, err
	}

	feed.Channel.Title = html.UnescapeString(feed.Channel.Title)
	feed.Channel.Description = html.UnescapeString(feed.Channel.Description)

	for _, item := range feed.Channel.Item {
		item.Title = html.UnescapeString(item.Title)
		item.Description = html.UnescapeString(item.Description)
	}
	return feed, nil
}

func scrapeFeeds(s *state) error {
	feedToFetch, err := s.database.GetNextFeedToFetch(context.Background())
	if err != nil {
		return err
	}

	feed, err := fetchFeed(context.Background(), feedToFetch.Url)
	if err != nil {
		return err
	}
	feed.Print()

	err = s.database.MarkFeedFetched(context.Background(), feedToFetch.ID)
	if err != nil {
		return err
	}

	return nil
}
