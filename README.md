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

## Development Sandbox with Podman

For safe testing without touching your host filesystem, you can run GoBackup inside an isolated Alpine Linux container using podman-compose.

```shell
# build
podman-compose up --build -d

# drop into the commandline
podman exec -it go-backup /bin/bash

# rebuild
podman-compose down; podman-compose up --build -d

# copy config example
podman cp ./config.example.yaml go-backup:/app
```

__Notes:__
- `/app` inside the container is your mounted project source directory.
- `/mnt-backup` inside the container is the safe target for backup runs.
- Always rebuild after changing Go source files to refresh the binary.

## Quick Sandbox Test
Once you’re inside the container:
```shell
# create safe target directory
mkdir /mnt-backup

# create some dummy source data
mkdir photos home
touch photos/photo1.jpg photos/photo2.png
mkdir -p home/Music
touch home/Favorite-Recipes.doc home/Music/song1.mp3 home/Music/song2.mp3

# write an example config that matches this sandbox layout
cat > /app/config.example.yaml <<'EOF'
backup_root: /mnt-backup
rsync:
  options: ["--numeric-ids", "--human-readable", "--itemize-changes"]
  excludes: [".cache", "node_modules", "Thumbs.db"]
jobs: 2
items:
  - name: home
    path: /app/home
    dest: home
  - name: photos
    path: /app/photos
EOF

# 1. Dry-run: see what would be backed up
gobackup --config /app/config.example.yaml --dry-run --verbose

# 2. Real run: copy data into /mnt-backup
gobackup --config /app/config.example.yaml --verbose

# 3. Inspect backup contents
tree /mnt-backup

# --- CHANGES DEMO ---
# Remove a file from source
rm /app/home/Music/song2.mp3

# Dry-run again to preview deletion
gobackup --config /app/config.example.yaml --dry-run --verbose

# Real run to apply deletion in backup
gobackup --config /app/config.example.yaml --verbose

# Verify file is gone from backup
tree /mnt-backup
```


### License: BSD