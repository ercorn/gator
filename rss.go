package main

import (
	"context"
	"database/sql"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
	"time"

	"github.com/ercorn/gator/internal/database"
	"github.com/google/uuid"
)

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

func fetchFeed(ctx context.Context, feedUrl string) (*RSSFeed, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", feedUrl, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}
	req.Header.Set("User-Agent", "gator")

	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %w", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading request: %w", err)
	}

	feed := RSSFeed{}
	err = xml.Unmarshal(body, &feed)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling xml: %w", err)
	}

	//decode escaped HTML entities (like &ldquo;) with html.UnescapeString(...)
	//specifically decode Title and Description fields of the channel and the
	//individual items
	feed.Channel.Title = html.UnescapeString(feed.Channel.Title)
	feed.Channel.Description = html.UnescapeString(feed.Channel.Description)
	for i := range feed.Channel.Item {
		feed.Channel.Item[i].Title = html.UnescapeString(feed.Channel.Item[i].Title)
		feed.Channel.Item[i].Description = html.UnescapeString(feed.Channel.Item[i].Description)
	}

	return &feed, nil
}

func handlerAgg(s *state, cmd command) error {
	_ = s
	if len(cmd.args) != 0 {
		return fmt.Errorf("usage: %s", cmd.name)
	}

	url := "https://www.wagslane.dev/index.xml"
	ctx := context.Background()
	feed, err := fetchFeed(ctx, url)
	if err != nil {
		return fmt.Errorf("failed to feed rss feed: %w", err)
	}

	fmt.Println(feed)
	return nil
}

func handlerAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 2 {
		return fmt.Errorf("usage: %s <feed_name> <url>", cmd.name)
	}

	ctx := context.Background()
	// user, err := s.db.GetUser(ctx, s.cfg.CurrentUserName)
	// if err != nil {
	// 	return fmt.Errorf("failed to get user: %w", err)
	// }

	feed, err := s.db.CreateFeed(ctx, database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      cmd.args[0],
		Url:       cmd.args[1],
		UserID:    user.ID,
	})
	if err != nil {
		return fmt.Errorf("failed to create feed: %w", err)
	}

	fmt.Println("Feed created successfully:")
	printFeed(feed, user)
	fmt.Println("\n======================================")

	//Automatically create a feed follow record for the current user when they add a feed.
	f_f_cmd := command{
		name: "follow",
		args: []string{cmd.args[1]},
	}
	err = handlerFollow(s, f_f_cmd, user)
	if err != nil {
		return fmt.Errorf("failed call to <follow> command handler: %w", err)
	}

	return nil

}

func handlerFeeds(s *state, cmd command) error {
	if len(cmd.args) != 0 {
		return fmt.Errorf("usage: %s", cmd.name)
	}

	ctx := context.Background()
	feeds, err := s.db.GetFeeds(ctx)
	if err != nil {
		return fmt.Errorf("failed to get list of feeds: %w", err)
	}

	if len(feeds) == 0 {
		fmt.Println("No feeds found.")
		return nil
	}

	fmt.Printf("Found %d feeds:\n", len(feeds))
	for _, feed := range feeds {
		fmt.Println("Feeds:")
		user, err := s.db.GetUserByID(ctx, feed.UserID)
		if err != nil {
			return fmt.Errorf("failed to get user of feed creator: %w", err)
		}
		printFeed(feed, user)
		fmt.Println("======================================================")
	}

	return nil
}

func printFeed(feed database.Feed, user database.User) {
	fmt.Printf("* ID:            %s\n", feed.ID)
	fmt.Printf("* Created:       %v\n", feed.CreatedAt)
	fmt.Printf("* Updated:       %v\n", feed.UpdatedAt)
	fmt.Printf("* Name:          %s\n", feed.Name)
	fmt.Printf("* URL:           %s\n", feed.Url)
	fmt.Printf("* User:          %s\n", user.Name)
}

func scrapeFeeds(s *state) {
	next_feed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		fmt.Println("failed to get next feed:", err)
		return
	}

	now := time.Now().UTC()
	err = s.db.MarkFeedFetched(context.Background(), database.MarkFeedFetchedParams{
		ID: next_feed.ID,
		LastFetchedAt: sql.NullTime{
			Time:  now,
			Valid: true,
		},
		UpdatedAt: now,
	})
	if err != nil {
		fmt.Println("failed to mark feed as fetched:", err)
		return
	}

	fetched_feed, err := fetchFeed(context.Background(), next_feed.Url)
	if err != nil {
		fmt.Println("failed to fetch feed:", err)
		return
	}

	fmt.Println("Feed Item Titles:")
	for _, item := range fetched_feed.Channel.Item {
		fmt.Println("-", item.Title)
	}
}
