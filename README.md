# GoBackup

![logo.png](logo.png)

A small Go application that reads a YAML configuration of important directories and uses `rsync` to keep them up to date in a specified backup location.

Designed to be boring, reliable, and easy to audit — no daemons, no magic sync protocols, just plain old `rsync`.

## Features

- **Simple YAML config** for sources, destinations, and `rsync` options
- **Multiple directories** backed up under a single root
- **Concurrent jobs** for faster operation
- **Dry-run mode** for safe testing
- **Verbose logging**
- **Safety checks** to prevent destructive mistakes (e.g., syncing into `/`)

## Usage

```shell
Usage: gobackup [flags]

Flags:
  --config string     Path to YAML config (default "config.yaml")
  --dry-run           Print `rsync` actions without making changes
  --verbose           Enable verbose logging
  --timeout duration  Per-task timeout (default 12h)
```

Example:
```shell
./gobackup --config config.example.yaml --dry-run --verbose
```

## Configuration
Example `config.example.yaml`:
```yaml
backup_root: /mnt/backup

rsync:
  options: ["-a", "--delete", "--numeric-ids"]
  excludes: [".cache", "node_modules", ".git"]

jobs: 2

items:
  - name: home
    path: /home/my-user
    dest: home
  - name: photos
    path: /home/my-user/Pictures
  - name: dotfiles
    path: /home/my-user/.config
    dest: config
```

### Safety Notes
- Always run with --dry-run after changing your config.
- Never point backup_root at / or a mount that might disappear.
- Consider adding --one-file-system to avoid crossing mount points.
- Use --numeric-ids when backing up system files to preserve ownership.

### License: BSD