package main

import (
	"fmt"
	"os"
)

func Check(err error) {
	if err != nil {
		panic(err)
	}
}

func parseCommands(args []string) {
	subcommand := "help"
	if len(args) >= 2 {
		subcommand = args[1]
		args = args[1:]
	}
	fmt.Println(subcommand)
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
