package main

import "fmt"

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
