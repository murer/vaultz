package main

import (
	"log"
	"os"
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

func createCommands() map[string]Command {
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

func handleCommands(args []string) {
	commands := createCommands()
	subcommand := "help"
	if len(args) >= 2 {
		subcommand = args[1]
		args = args[2:]
	}
	log.Printf("Command: %s, args: %s\n", subcommand, args)
	command := commands[subcommand]
	if command == nil {
		log.Panicf("Wrong command: %s, try to use help", subcommand)
	}
	command.Run()
}

func main() {
	handleCommands(os.Args)
}
