# Write Disk (File Generator)

A high-performance file generator tool designed to create large files or fill up disk space for testing purposes.

## Build

```bash
# Using Make
make build

# Using Go directly
go build -o write_disk main.go
```

## Usage

```bash
./write_disk [flags]
```

### Flags

- `-s, --size int`: **(Required)** Size of each file in GB.
- `-n, --number int`: **(Required)** Number of files to create.
- `-p, --path string`: **(Required)** Output directory path.
- `-f, --prefix string`: Filename prefix (default "file").
- `--fast`: Enable fast pre-allocation mode using `ftruncate`.
  - **Enabled**: Creates sparse files instantly (checking file existence/metadata).
  - **Disabled (Default)**: Actually writes bytes (zeros) to disk, consuming real I/O and time.

### Examples

**1. Fast Mode (Metadata only)**
Create 10 files of 10GB each instantly.
```bash
./write_disk -s 10 -n 10 -p /data/test --fast
```

**2. Real Write Mode (I/O Stress)**
Actually write 1GB of data 5 times to disk.
```bash
./write_disk -s 1 -n 5 -p /data/test
```
