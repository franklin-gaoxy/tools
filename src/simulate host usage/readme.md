# Simulate Host Usage (shub)

A tool to simulate system resource usage, including CPU load, Memory consumption, and Disk I/O. Useful for testing monitoring systems or stress testing.

## Build

```bash
# Using Make
make build

# Using Go directly
go build -o shub main.go
```

## Usage

```bash
./shub [command] [flags]
```

### Commands

#### 1. `runcpu`
Simulate CPU load.

- `-l, --load string`: CPU load percentage (0.0 - 1.0). Default `0.8` (80%).
- `-t, --time int`: Duration in minutes. Default `0` (run until stopped).

**Example:**
```bash
# 50% CPU load for 10 minutes
./shub runcpu -l 0.5 -t 10
```

#### 2. `memory`
Simulate Memory usage.

- `-s, --size string`: Memory usage ratio (0.0 - 1.0). Default `0.8` (80%).
- `-t, --time int`: Duration in minutes.

**Example:**
```bash
# Consume 80% memory
./shub memory -s 0.8
```

#### 3. `disk`
Simulate Disk usage (creates files).

- `-s, --size int`: Size of each file in GB.
- `-p, --path string`: Directory path to write files.
- `-t, --time int`: Duration in minutes.

**Example:**
```bash
# Create 1GB files in /tmp/test-disk
./shub disk -s 1 -p /tmp/test-disk
```
