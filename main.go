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
	// Name() string
	// Flags() *flag.FlagSet
	Run()
}

type BaseCommand struct {
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

func main() {
	a := BaseCommand{}
	a.Run()

	b := HelpCommand{}
	b.Run()

	cmds := []Command{
		&BaseCommand{},
		&HelpCommand{},
	}
	log.Printf("kkkkk: %#v\n", cmds)
	for _, cmd := range cmds {
		cmd.Run()
	}

}
