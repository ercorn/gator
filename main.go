package main

import (
	"fmt"
	"log"
	"os"

	"github.com/ercorn/gator/internal/config"
	_ "github.com/lib/pq"
)

type state struct {
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
		log.Fatal(fmt.Errorf("the login handler expects a single argument, the username"))
	}
	s.cfg.SetUser(cmd.args[0])
	fmt.Println("user has been set")
	return nil
}

func main() {
	//cfg := config.Read()
	//cfg.SetUser("ronyo")
	//cfg = config.Read()
	//fmt.Println(cfg)

	curr_state := &state{
		cfg: config.Read(),
	}

	cmds := commands{
		cmd_map: make(map[string]func(*state, command) error),
	}

	cmds.register("login", handlerLogin)
	if len(os.Args) < 2 {
		log.Fatal(fmt.Errorf("error: please provide a command name"))
	}
	cmd := command{
		name: os.Args[1],
		args: os.Args[2:],
	}
	fmt.Println("input", os.Args[1])
	err := cmds.run(curr_state, cmd)
	if err != nil {
		fmt.Println(err)
	}
}
