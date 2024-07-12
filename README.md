# GoBackup

![logo.png](logo.png)

Utility for backing up files. It moves them from one place to another.

## Features

- Backup and Restore files
- Multiple backup locations
- Support for s3 providers
- Compression and encryption
- Systemd configuration

## Usage

```shell
Usage: gobackup [command] [args]

Commands:
  engage     Run the backup
  restore    Restore local from backup
  status     Get info about backup status

Use "gobackup [command] --help" for more information about a command.
```

## Configuration

Here's an example for now.

```json5
{
  "sourceDirs": [
    "~/Documents",
    "~/Photos",
    "~/projects",
    "~/Files"
  ],
  "destinations": [
    {
      "name": "HDD 2",
      "path": "/run/media/user/HDD2",
      "encrypt": true
    },
    {
      "name": "cloud storage",
      "path": "s3providers.com/mybucket",
      "encrypt": true,
      "username": "My_user",
      "password": "topsecret"
    }
  ],
  "log" : "none" // "none" or "syslog"
}
```