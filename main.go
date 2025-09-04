package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/ercorn/gator/internal/config"
	"github.com/ercorn/gator/internal/database"
	_ "github.com/lib/pq"
)

type state struct {
	db  *database.Queries
	cfg *config.Config
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
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerReset)
	cmds.register("users", handlerUsers)
	cmds.register("agg", handlerAgg)

	//parse arguments and run the requested command
	if len(os.Args) < 2 {
		log.Fatal(fmt.Errorf("error: please provide a command name"))
	}
	cmd := command{
		name: os.Args[1],
		args: os.Args[2:],
	}

	err = cmds.run(progam_state, cmd)
	if err != nil {
		log.Fatal(err)
	}
}
