package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/goccy/go-yaml"
)

type Config struct {
	BackupRoot string   `yaml:"backup_root"`
	Rsync      RsyncCfg `yaml:"rsync"`
	Jobs       int      `yaml:"jobs"`
	Items      []Item   `yaml:"items"`
}

type RsyncCfg struct {
	Options  []string `yaml:"options"`
	Excludes []string `yaml:"excludes"`
}

type Item struct {
	Name string `yaml:"name"`
	Path string `yaml:"path"`
	Dest string `yaml:"dest"`
}

func readConfig(path string) (*Config, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var c Config
	if err := yaml.Unmarshal(b, &c); err != nil {
		return nil, err
	}
	if c.BackupRoot == "" {
		return nil, errors.New("backup_root is required")
	}
	if len(c.Items) == 0 {
		return nil, errors.New("items is empty")
	}
	if c.Jobs <= 0 {
		c.Jobs = max(1, runtime.NumCPU()/2)
	}
	// Normalize rsync excludes once: trim whitespace, drop empties, deduplicate.
	c.Rsync.Excludes = uniqTrimmed(c.Rsync.Excludes)

	return &c, nil
}

// uniqTrimmed returns a copy of ss with whitespace-trimmed entries,
// all empty strings removed, and duplicates eliminated while preserving order.
func uniqTrimmed(ss []string) []string {
	seen := make(map[string]struct{}, len(ss))
	out := make([]string, 0, len(ss))
	for _, s := range ss {
		s = strings.TrimSpace(s)
		if s == "" {
			continue
		}
		if _, ok := seen[s]; ok {
			continue
		}
		seen[s] = struct{}{}
		out = append(out, s)
	}
	return out
}

func safePath(p string) (string, error) {
	ap, err := filepath.Abs(p)
	if err != nil {
		return "", err
	}
	if ap == "/" {
		return "", fmt.Errorf("refusing to use '/' as a path: %s", p)
	}
	return ap, nil
}

func main() {
	var (
		cfgPath = flag.String("config", "config.yaml", "path to YAML config")
		dryRun  = flag.Bool("dry-run", false, "print rsync actions without changing anything")
		verbose = flag.Bool("verbose", false, "verbose logging")
		timeout = flag.Duration("timeout", 12*time.Hour, "per-task timeout")
	)
	flag.Parse()

	logger := log.New(os.Stdout, "gobackup ", log.LstdFlags)

	cfg, err := readConfig(*cfgPath)
	if err != nil {
		log.Fatal(err)
	}

	backupRoot, err := safePath(cfg.BackupRoot)
	if err != nil {
		log.Fatal(err)
	}

	if _, err := exec.LookPath("rsync"); err != nil {
		log.Fatalf("rsync not found in PATH: %v", err)
	}

	tasks := make(chan Item)
	var wg sync.WaitGroup
	for i := 0; i < cfg.Jobs; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for it := range tasks {
				if err := syncItem(it, backupRoot, cfg.Rsync, *dryRun, *verbose, *timeout, logger); err != nil {
					logger.Printf("[worker %d] %s: ERROR: %v", id, label(it), err)
				} else {
					logger.Printf("[worker %d] %s: OK", id, label(it))
				}
			}
		}(i + 1)
	}

	for _, it := range cfg.Items {
		tasks <- it
	}
	close(tasks)
	wg.Wait()
}

func label(it Item) string {
	name := it.Name
	if name == "" {
		name = filepath.Base(it.Path)
	}
	return name
}

func syncItem(it Item, backupRoot string, rc RsyncCfg, dry, verbose bool, timeout time.Duration, logger *log.Logger) error {
	src, err := safePath(it.Path)
	if err != nil {
		return err
	}
	if _, err := os.Stat(src); err != nil {
		return fmt.Errorf("source missing: %w", err)
	}

	destName := it.Dest
	if destName == "" {
		destName = filepath.Base(src)
	}
	dest := filepath.Join(backupRoot, destName)
	if err := os.MkdirAll(dest, 0o755); err != nil {
		return err
	}

	args := []string{"-a", "--delete"}
	if verbose {
		args = append(args, "-v")
	}
	args = append(args, rc.Options...)
	for _, ex := range rc.Excludes {
		// already trimmed & deduplicated in readConfig
		args = append(args, "--exclude", ex)
	}
	if dry {
		args = append(args, "--dry-run")
	}

	srcWithSlash := src
	if !strings.HasSuffix(srcWithSlash, string(os.PathSeparator)) {
		srcWithSlash += string(os.PathSeparator)
	}
	args = append(args, srcWithSlash, dest+string(os.PathSeparator))

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, "rsync", args...)
	cmd.Stdout = logger.Writer()
	cmd.Stderr = logger.Writer()

	if verbose || dry {
		logger.Printf("rsync %s", shellEscape(args))
	}
	return cmd.Run()
}

func shellEscape(args []string) string {
	q := make([]string, len(args))
	for i, a := range args {
		if strings.ContainsAny(a, " \"'\t\n$") {
			q[i] = "'" + strings.ReplaceAll(a, "'", "'\\''") + "'"
		} else {
			q[i] = a
		}
	}
	return strings.Join(q, " ")
}
