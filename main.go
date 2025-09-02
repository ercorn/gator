package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ercorn/gator/internal/config"
	"github.com/ercorn/gator/internal/database"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type state struct {
	db  *database.Queries
	cfg *config.Config
}

type command struct {
	name string
	args []string
}

type commands struct {
	cmd_map map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmd command) error {
	if cmd_func, exists := c.cmd_map[cmd.name]; exists {
		fmt.Println("running command")
		err := cmd_func(s, cmd)
		if err != nil {
			return err
		}
		return nil
	}
	return fmt.Errorf("command <%s> does not exist", cmd.name)
}

func (c *commands) register(name string, f func(s *state, cmd command) error) error {
	c.cmd_map[name] = f
	return nil
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("usage: %s <name>", cmd.name)
	}

	//check if current username exists in the database, error if not
	ctx := context.Background()
	_, err := s.db.GetUser(ctx, cmd.args[0])
	if err != nil {
		return fmt.Errorf("failed, username doesn't exist in the database: %w", err)
	}

	err = s.cfg.SetUser(cmd.args[0])
	if err != nil {
		return fmt.Errorf("couldn't set the user: %w", err)
	}
	fmt.Println("user has been set!")
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("usage: %s <name>", cmd.name)
	}

	ctx := context.Background()
	_, err := s.db.CreateUser(ctx, database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      cmd.args[0],
	})
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	err = s.cfg.SetUser(cmd.args[0])
	if err != nil {
		return fmt.Errorf("couldn't set the user: %w", err)
	}

	fmt.Println("user was created:", s.cfg.CurrentUserName)

	return nil
}

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}

	db, err := sql.Open("postgres", cfg.DBUrl) //open connection to db
	if err != nil {
		log.Fatalf("error opening db connection: %v", err)
	}
	defer db.Close()
	dbQueries := database.New(db)

	progam_state := &state{
		db:  dbQueries,
		cfg: &cfg,
	}

	cmds := commands{
		cmd_map: make(map[string]func(*state, command) error),
	}

	cmds.register("login", handlerLogin)
	//register "register" command here, usage: "go run . register lane"
	cmds.register("register", handlerRegister)

	//parse arguments and run the requested command
	if len(os.Args) < 2 {
		log.Fatal(fmt.Errorf("error: please provide a command name"))
	}
	cmd := command{
		name: os.Args[1],
		args: os.Args[2:],
	}
	//fmt.Println("input", os.Args[1])
	err = cmds.run(progam_state, cmd)
	if err != nil {
		log.Fatal(err)
	}
}
