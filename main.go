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
	Name() string
	Flags() *flag.FlagSet
	Run(args []string)
}

type HelpCommand struct {
	cmds map[string]Command
}

func (me *HelpCommand) Name() string { return "help" }

func (me *HelpCommand) Flags() *flag.FlagSet {
	return flag.NewFlagSet(me.Name(), flag.ExitOnError)
}

func (me *HelpCommand) Run(args []string) {
	me.Flags().Output().Write([]byte{10})
	for _, cmd := range me.cmds {
		cmd.Flags().Usage()
		os.Stdout.Write([]byte{10, 10})
	}
}

type KeygenCommand struct {
	flags    *flag.FlagSet
	flagName *string
}

func (me *KeygenCommand) Name() string { return "keygen" }

func (me *KeygenCommand) Flags() *flag.FlagSet {
	if me.flags == nil {
		me.flags = flag.NewFlagSet(me.Name(), flag.ExitOnError)
		me.flagName = me.flags.String("name", "", "Key name")
	}
	return me.flags
}

func (me *KeygenCommand) Run(args []string) {
	err := me.Flags().Parse(args)
	Check(err)
	log.Printf("Keygen name: %s\n", *me.flagName)
}

func createCommands() map[string]Command {
	ret := make(map[string]Command)
	func(cmds []Command) {
		for _, element := range cmds {
			ret[element.Name()] = element
		}
	}([]Command{
		&HelpCommand{ret},
		&KeygenCommand{},
	})
	return ret
}

func parseCommands(args []string) {
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
	command.Run(args)
}

func main() {
	parseCommands(os.Args)
	// if len(os.Args) < 2 {
	// 	fmt.Println("Expected a subcommand")
	// 	os.Exit(1)
	// }

	// // Subcommands
	// serveCmd := flag.NewFlagSet("serve", flag.ExitOnError)
	// servePort := serveCmd.Int("port", 8080, "Port to run the server on")

	// migrateCmd := flag.NewFlagSet("migrate", flag.ExitOnError)
	// migrateDir := migrateCmd.String("dir", "./migrations", "Directory with migration files")

	// // Check which subcommand is invoked
	// switch os.Args[1] {
	// case "serve":
	// 	// Parse flags for the 'serve' subcommand
	// 	serveCmd.Parse(os.Args[2:])
	// 	fmt.Printf("Serving on port %d...\n", *servePort)

	// case "migrate":
	// 	// Parse flags for the 'migrate' subcommand
	// 	migrateCmd.Parse(os.Args[2:])
	// 	fmt.Printf("Running migrations from directory '%s'...\n", *migrateDir)

	// default:
	// 	fmt.Println("Expected 'serve' or 'migrate' subcommands")
	// 	os.Exit(1)
	// }
}
