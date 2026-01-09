# Disk Not Sleep

A utility tool designed to prevent hard drives (especially external HDDs) from entering sleep/spin-down mode by periodically performing a lightweight write-delete operation.

## Build

```bash
# Using Go directly
go build -o disk_not_sleep main.go
```

## Usage

```bash
./disk_not_sleep [flags]
```

### Global Flags

#### `-p, --path string`
Target directory path.
- **Default**: `.` (Current directory)
- **Function**: Specifies the directory where the temporary keep-alive file will be created. Select the mount point of the disk you want to keep awake.

#### `-f, --filename string`
Temporary file name.
- **Default**: `tmp_file`
- **Function**: The name of the temporary file used for the write operation.

#### `-t, --time int`
Interval time.
- **Default**: `60`
- **Unit**: Seconds
- **Function**: How often the write-delete cycle runs.

---

### Examples

#### Basic Usage
Prevent the disk mounted at `/mnt/external_drive` from sleeping, writing every 60 seconds (default).

```bash
./disk_not_sleep -p /mnt/external_drive
```

#### Custom Interval and Filename
Write to `/Volumes/MyDisk` every 5 minutes (300 seconds) using a custom filename `.keep_alive`.

```bash
./disk_not_sleep -p /Volumes/MyDisk -t 300 -f .keep_alive
```
