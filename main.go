package main

import (
	"log"
)

func Check(err error) {
	if err != nil {
		panic(err)
	}
}

type Command interface {
	GetName() string
	// Flags() *flag.FlagSet
	Run()
}

type BaseCommand struct {
	Name string

	Cmds map[string]Command
}

func (me *BaseCommand) GetName() string {
	return me.Name
}

func (me *BaseCommand) Run() {
	log.Println("aaaa")
}

type HelpCommand struct {
	BaseCommand
}

func (me *HelpCommand) Run() {
	log.Println("bbb")
}

func parseCommands() map[string]Command {
	ret := make(map[string]Command)
	(func(cmds []Command) {
		for _, cmd := range cmds {
			ret[cmd.GetName()] = cmd
		}
	})([]Command{
		&HelpCommand{BaseCommand{Name: "help"}},
	})
	return ret
}

func main() {
	parseCommands()
}
