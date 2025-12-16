package main

import (
	"fmt"
)

/*
Add a follow command. It takes a single url argument and creates a new feed follow record
for the current user. It should print the name of the feed and the current user once the
record is created (which the query we just made should support). You'll need a query to
look up feeds by URL.
*/
func handlerFollow(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("usage: %s", cmd.name)
	}

	//
	return nil
}
