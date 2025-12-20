package main

import (
	"context"
	"fmt"
	"time"

	"github.com/ercorn/gator/internal/database"
	"github.com/google/uuid"
)

/*
Add a follow command. It takes a single url argument and creates a new feed follow record
for the current user. It should print the name of the feed and the current user once the
record is created (which the query we just made should support). You'll need a query to
look up feeds by URL.
*/
func handlerFollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("usage: %s <url>", cmd.name)
	}

	// //get current user record
	// user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
	// if err != nil {
	// 	return fmt.Errorf("failed to get user: %w", err)
	// }

	//parse url and use it to get corresponding feed record
	url := cmd.args[0]
	feed, err := s.db.GetFeedByUrl(context.Background(), url)
	if err != nil {
		//couldn't get feed
		return fmt.Errorf("failed to get feed: %w", err)
	}

	feed_follow, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	})

	if err != nil {
		return fmt.Errorf("failed to follow the feed: %w", err)
	}

	fmt.Println("Feed followed successfully!")
	fmt.Println("Username:", feed_follow.UserName)
	fmt.Println("Feed name:", feed_follow.FeedName)
	fmt.Println("==================================================")

	return nil
}

// Add a following command. It should print all the names of the feeds the current user is following.
func handlerFollowing(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 0 {
		return fmt.Errorf("usage: %s", cmd.name)
	}

	// //get current user from db by name
	// user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
	// if err != nil {
	// 	return fmt.Errorf("failed to get current user: %w", err)
	// }

	//get array of feeds followed by the current user
	feed_follows, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return fmt.Errorf("failed to get feed follows for the current user: %w", err)
	}

	if len(feed_follows) == 0 {
		fmt.Println("No feed follows found for this user.")
		return nil
	}

	fmt.Println("Current username:", user.Name)
	fmt.Println("Feed names:")
	for _, follow := range feed_follows {
		fmt.Println("-", follow.FeedName)
	}
	return nil
}

func handlerUnfollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("usage: %s <url>", cmd.name)
	}

	feed, err := s.db.GetFeedByUrl(context.Background(), cmd.args[0])
	if err != nil {
		return fmt.Errorf("failed to get feed from url: %w", err)
	}

	err = s.db.DeleteFeedFollow(context.Background(), database.DeleteFeedFollowParams{
		UserID: user.ID,
		FeedID: feed.ID,
	})
	if err != nil {
		return fmt.Errorf("failed to delete feed follow: %w", err)
	}

	fmt.Println("DELETED FOLLOWED FEED:", feed.Name, "FROM USER:", user.Name)

	return nil
}
