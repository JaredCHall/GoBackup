package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
)

type Destination struct {
	Name     string `json:"name"`
	Path     string `json:"path"`
	Encrypt  bool   `json:"encrypt"`
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

type Config struct {
	SourceDirs   []string      `json:"sourceDirs"`
	Destinations []Destination `json:"destinations"`
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func FindConfigFilePath() (string, error) {
	usr, err := user.Current()
	if err != nil {
		fmt.Println("Error fetching current user:", err)
		return "", err
	}
	homeDirPath := filepath.Join(usr.HomeDir, ".gobackup", "config.json")
	etcPath := "/etc/gobackup/config.json"

	if fileExists(homeDirPath) {
		return homeDirPath, nil
	}
	if fileExists(etcPath) {
		return etcPath, nil
	}
	return "", errors.New("no file at " + homeDirPath + " or " + etcPath)
}

func ReadConfig(filename string) (Config, error) {
	file, err := os.Open(filename)
	if err != nil {
		return Config{}, err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Printf("Error closing file: %v\n", err)
		}
	}(file)

	var config Config
	decoder := json.NewDecoder(bufio.NewReader(file))
	if err := decoder.Decode(&config); err != nil {
		return Config{}, err
	}

	return config, nil
}

func LoadConfig() (Config, error) {
	configPath, err := FindConfigFilePath()
	if err != nil {
		fmt.Printf("Error locating config file: %v\n", err)
		return Config{}, err
	}

	return ReadConfig(configPath)
}

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
