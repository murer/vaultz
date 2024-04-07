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
	Name() string
	// Flags() *flag.FlagSet
	Run()
}

type BaseCommand struct {
	name string
}

func (me *BaseCommand) Name() string {
	return me.name
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

		}
	})([]Command{
		&BaseCommand{},
		&HelpCommand{},
	})
	return ret
}

func main() {
	a := BaseCommand{}
	a.Run()

	b := HelpCommand{}
	b.Run()

	cmds := []Command{
		&BaseCommand{"base"},
		&HelpCommand{BaseCommand: BaseCommand{"help"}},
	}
	log.Printf("kkkkk: %#v\n", cmds)
	for _, cmd := range cmds {
		cmd.Run()
	}

}
