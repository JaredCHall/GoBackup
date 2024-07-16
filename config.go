package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"
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

func LoadConfig(findConfigFilePathFunc ...func(func(string) bool) (string, error)) (Config, error) {
	var findFunc func(func(string) bool) (string, error)
	if len(findConfigFilePathFunc) > 0 {
		findFunc = findConfigFilePathFunc[0]
	} else {
		findFunc = findConfigFilePath // default function
	}

	configPath, err := findFunc(fileExists)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error locating config file: %v\n", err)
		return Config{}, err
	}

	return readConfig(configPath)
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// finds config file according to XDG specification: https://specifications.freedesktop.org/basedir-spec/latest
func findConfigFilePath(fileExistsFunc func(string) bool) (string, error) {
	path := ""
	if configHome := os.Getenv("XDG_CONFIG_HOME"); configHome != "" {
		path = filepath.Join(configHome, ".config/GoBackup.json5")
	} else {
		path = filepath.Join(os.Getenv("HOME"), ".config/GoBackup.json5")
	}

	if fileExistsFunc(path) {
		return path, nil
	}

	return "", errors.New("no config file found at " + path)
}

func readConfig(filename string) (Config, error) {
	file, err := os.Open(filename)
	if err != nil {
		return Config{}, err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error closing file: %v\n", err)
		}
	}(file)

	var config Config
	decoder := json.NewDecoder(bufio.NewReader(file))
	if err := decoder.Decode(&config); err != nil {
		return Config{}, err
	}

	return config, nil
}
