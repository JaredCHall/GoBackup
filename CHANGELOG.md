## Changelog

All notable changes to this project will be documented in this file.
The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).


### [Unreleased]

#### Added

- ChangeLog - First change log entries; format: Keep a ChangeLog 

#### Changed

#### Fixed

---

### [0.1.0] — 2025-08-15

#### Added

- Initial release of **GoBackup**.
- YAML-based configuration for:
    - Backup root directory.
    - `rsync` options and excludes.
    - Multiple backup items with per-item `name`, `path`, and optional `dest`.
    - Configurable number of concurrent jobs.
- Backup execution using `rsync` with:
    - Dry-run mode for previewing changes.
    - Verbose logging for transparency.
    - Safety checks to prevent destructive syncs (e.g., syncing into `/`).
- Example `config.example.yaml` provided.
- **Podman sandbox environment** for safe development and testing:
    - Container build with `podman-compose`.
    - Mounted project source and isolated `/mnt-backup` target.
    - Step-by-step sandbox test instructions.
- BSD license.