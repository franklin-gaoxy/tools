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

---

### Global Flags

#### `-s, --size int`
File size.
- **Required**: Yes
- **Unit**: GB
- **Function**: Specifies the size of a single file to be generated.

#### `-n, --number int`
Number of files.
- **Required**: Yes
- **Function**: Specifies how many files of the given size to create.

#### `-p, --path string`
Output directory.
- **Required**: Yes
- **Function**: The absolute or relative path where files will be written.

#### `-f, --prefix string`
Filename prefix.
- **Default**: "file"
- **Function**: The prefix string for generated filenames (e.g., `prefix_1`, `prefix_2`).

#### `--fast`
Enable Fast Mode.
- **Type**: Boolean flag
- **Function**:
    - **If present**: Uses `ftruncate` (file truncation) to "pre-allocate" the file size instantly. It updates file metadata to show the size but does not write actual zeros to every block on the disk (creates sparse files). **Use this for testing file count limits or metadata operations.**
    - **If absent**: Writes actual data (zeros) to the disk block by block. **Use this for testing actual Disk I/O bandwidth and physical storage capacity.**

---

### Examples

#### 1. Fast Mode (Metadata only)
Create 10 files of 10GB each instantly.
```bash
./write_disk -s 10 -n 10 -p /data/test --fast
```

#### 2. Real Write Mode (I/O Stress)
Actually write 1GB of data 5 times to disk.
```bash
./write_disk -s 1 -n 5 -p /data/test
```
