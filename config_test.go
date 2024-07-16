package main

import (
	"errors"
	"os"
	"reflect"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	tests := []struct {
		name               string
		findConfigFilePath func(func(string) bool) (string, error)
		mockConfigFilePath string
		mockConfigContent  string
		expectedConfig     Config
		expectingError     bool
	}{
		{
			name: "ValidConfig",
			findConfigFilePath: func(fileExists func(string) bool) (string, error) {
				return "valid_config.json", nil
			},
			mockConfigContent: `{
				"sourceDirs": ["/path/to/source1", "/path/to/source2"],
				"destinations": [
					{
						"name": "dest1",
						"path": "/path/to/dest1",
						"encrypt": true,
						"username": "user1",
						"password": "pass1"
					},
					{
						"name": "dest2",
						"path": "/path/to/dest2",
						"encrypt": false
					}
				]
			}`,
			expectedConfig: Config{
				SourceDirs: []string{"/path/to/source1", "/path/to/source2"},
				Destinations: []Destination{
					{
						Name:     "dest1",
						Path:     "/path/to/dest1",
						Encrypt:  true,
						Username: "user1",
						Password: "pass1",
					},
					{
						Name:    "dest2",
						Path:    "/path/to/dest2",
						Encrypt: false,
					},
				},
			},
			expectingError: false,
		},
		{
			name: "ConfigFileNotFound",
			findConfigFilePath: func(fileExists func(string) bool) (string, error) {
				return "", errors.New("config file not found")
			},
			expectedConfig: Config{},
			expectingError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temporary file with mock content if needed
			var tmpfile *os.File
			var err error

			if tt.mockConfigContent != "" {
				tmpfile, err = os.CreateTemp("", "test_config_*.json")
				if err != nil {
					t.Fatal(err)
				}
				defer os.Remove(tmpfile.Name())

				if _, err := tmpfile.Write([]byte(tt.mockConfigContent)); err != nil {
					t.Fatal(err)
				}
				if err := tmpfile.Close(); err != nil {
					t.Fatal(err)
				}

				tt.findConfigFilePath = func(fileExists func(string) bool) (string, error) {
					return tmpfile.Name(), nil
				}
			}

			got, err := LoadConfig(tt.findConfigFilePath)
			if (err != nil) != tt.expectingError {
				t.Errorf("LoadConfig() error = %v, expectingError %v", err, tt.expectingError)
				return
			}
			if !reflect.DeepEqual(got, tt.expectedConfig) {
				t.Errorf("LoadConfig() got = %v, want %v", got, tt.expectedConfig)
			}
		})
	}
}

func Test_fileExists(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "executable exists",
			args: args{
				path: func() string {
					ex, err := os.Executable()
					if err != nil {
						t.Fatalf("Failed to get executable path: %v", err)
					}
					return ex
				}(),
			},
			want: true,
		},
		{
			name: "non-existent file",
			args: args{
				path: "/path/to/nonexistent/file",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := fileExists(tt.args.path); got != tt.want {
				t.Errorf("fileExists() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_findConfigFilePath(t *testing.T) {
	tests := []struct {
		name           string
		envVars        map[string]string
		mockFileExists func(path string) bool
		want           string
		wantErr        bool
	}{
		{
			name: "Config file in XDG_CONFIG_HOME",
			envVars: map[string]string{
				"XDG_CONFIG_HOME": "/tmp/config_home",
				"HOME":            "/tmp/home",
			},
			mockFileExists: func(path string) bool {
				return path == "/tmp/config_home/.config/GoBackup.json5"
			},
			want:    "/tmp/config_home/.config/GoBackup.json5",
			wantErr: false,
		},
		{
			name: "Config file in HOME",
			envVars: map[string]string{
				"XDG_CONFIG_HOME": "",
				"HOME":            "/tmp/home",
			},
			mockFileExists: func(path string) bool {
				return path == "/tmp/home/.config/GoBackup.json5"
			},
			want:    "/tmp/home/.config/GoBackup.json5",
			wantErr: false,
		},
		{
			name: "No config file found",
			envVars: map[string]string{
				"XDG_CONFIG_HOME": "",
				"HOME":            "/tmp/home",
			},
			mockFileExists: func(path string) bool {
				return false
			},
			want:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set environment variables
			for key, value := range tt.envVars {
				os.Setenv(key, value)
			}

			got, err := findConfigFilePath(tt.mockFileExists)
			if (err != nil) != tt.wantErr {
				t.Errorf("findConfigFilePath() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("findConfigFilePath() got = %v, want %v", got, tt.want)
			}

			// Unset environment variables
			for key := range tt.envVars {
				os.Unsetenv(key)
			}
		})
	}
}

func Test_readConfig(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		content  string
		want     Config
		wantErr  bool
	}{
		{
			name:     "ValidConfig",
			filename: "test_valid_config.json",
			content: `{
				"sourceDirs": ["/path/to/source1", "/path/to/source2"],
				"destinations": [
					{
						"name": "dest1",
						"path": "/path/to/dest1",
						"encrypt": true,
						"username": "user1",
						"password": "pass1"
					},
					{
						"name": "dest2",
						"path": "/path/to/dest2",
						"encrypt": false
					}
				]
			}`,
			want: Config{
				SourceDirs: []string{"/path/to/source1", "/path/to/source2"},
				Destinations: []Destination{
					{
						Name:     "dest1",
						Path:     "/path/to/dest1",
						Encrypt:  true,
						Username: "user1",
						Password: "pass1",
					},
					{
						Name:    "dest2",
						Path:    "/path/to/dest2",
						Encrypt: false,
					},
				},
			},
			wantErr: false,
		},
		{
			name:     "InvalidConfig",
			filename: "test_invalid_config.json",
			content:  `{ invalid json }`,
			want:     Config{},
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temporary file with test content
			tmpfile, err := os.CreateTemp("", tt.filename)
			if err != nil {
				t.Fatal(err)
			}
			defer os.Remove(tmpfile.Name())

			if _, err := tmpfile.Write([]byte(tt.content)); err != nil {
				t.Fatal(err)
			}
			if err := tmpfile.Close(); err != nil {
				t.Fatal(err)
			}

			got, err := readConfig(tmpfile.Name())
			if (err != nil) != tt.wantErr {
				t.Errorf("readConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("readConfig() got = %v, want %v", got, tt.want)
			}
		})
	}
}
