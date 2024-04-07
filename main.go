package main

import (
	"flag"
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
	PrepareFlags(flags *flag.FlagSet)
	Run()
}

type BaseCommand struct {
	Name string
	Cmds map[string]Command
}

func (me *BaseCommand) GetName() string {
	return me.Name
}

func (me *BaseCommand) PrepareFlags(flags *flag.FlagSet) {
}

func (me *BaseCommand) Run() {
	log.Println("aaaa")
}

type HelpCommand struct {
	BaseCommand
}

func (me *HelpCommand) Run() {
	// me.Flags().Output().Write([]byte{10})
	// for _, cmd := range me.cmds {
	// 	cmd.Flags().Usage()
	// 	os.Stdout.Write([]byte{10, 10})
	// }
}

type KeygenCommand struct {
	BaseCommand
	FlagName *string
}

func (me *KeygenCommand) Run() {
	log.Printf("Key gen name: %s\n", *me.FlagName)
}

func (me *KeygenCommand) PrepareFlags(flags *flag.FlagSet) {
	me.FlagName = flags.String("name", "", "Key name")
}

func createCommands() map[string]Command {
	ret := make(map[string]Command)
	(func(cmds []Command) {
		for _, cmd := range cmds {
			cmd.PrepareFlags(flag.NewFlagSet(cmd.GetName(), flag.ExitOnError))
			ret[cmd.GetName()] = cmd
		}
	})([]Command{
		&HelpCommand{BaseCommand: BaseCommand{"help", ret}},
		&KeygenCommand{BaseCommand: BaseCommand{"keygen", ret}},
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
