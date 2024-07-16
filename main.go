package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	// Create a flag set for the main command to capture global flags like -help
	mainFlagSet := flag.NewFlagSet("main", flag.ExitOnError)
	helpPtr := mainFlagSet.Bool("help", false, "Display this help message")
	shortHelpPtr := mainFlagSet.Bool("h", false, "Display this help message")

	// Parse the main flag set
	mainFlagSet.Parse(os.Args[1:])

	// Check if the help flag is set
	if *helpPtr || *shortHelpPtr || len(os.Args) < 2 {
		printMainHelp()
		os.Exit(0)
		return
	}

	// Handle subcommands
	switch os.Args[1] {
	case "engage":
		handleEngage(os.Args[2:])
	case "status":
		handleStatus(os.Args[2:])
	default:
		fmt.Printf("Unknown command: %s\n", os.Args[1])
		printMainHelp()
		os.Exit(1)
	}
}

func handleEngage(args []string) {
	engageFlagSet := flag.NewFlagSet("engage", flag.ExitOnError)
	destPtr := engageFlagSet.String("dest", "", "Name a specific destination for backup")
	allPtr := engageFlagSet.Bool("all", false, "Backup to all destinations")

	// Parse the engage flag set
	engageFlagSet.Parse(args)

	config, err := LoadConfig()
	if err != nil {
		fmt.Printf("Error reading config: %v\n", err)
		return
	}

	fmt.Println("Running backup...")

	fmt.Println("Dest:", *destPtr)
	fmt.Println("All:", *allPtr)
	printConfigInfo(config)

}

func handleStatus(args []string) {
	statusFlagSet := flag.NewFlagSet("status", flag.ExitOnError)

	// Parse the status flag set (no specific flags for status in this example)
	statusFlagSet.Parse(args)

	fmt.Println("Getting backup status...")
}

func printConfigInfo(config Config) {
	fmt.Println("Printing config information...")
	fmt.Println("Source Directories:")
	for _, dir := range config.SourceDirs {
		fmt.Println(" -", dir)
	}
	fmt.Println("Destinations:")
	for _, dest := range config.Destinations {
		fmt.Printf(" - %s (%s), Encrypt: %v\n", dest.Name, dest.Path, dest.Encrypt)
		if dest.Username != "" {
			fmt.Printf("   Username: %s, Password: %s\n", dest.Username, dest.Password)
		}
	}
}

func printMainHelp() {
	fmt.Printf(`Usage: gobackup [command] [args]

Commands:
  engage     Run the backup
  status     Get info about backup status

Use "gobackup [command] --help" for more information about a command.
`)
}
