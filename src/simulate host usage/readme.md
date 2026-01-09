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

---

### Command: `runcpu`

Simulate CPU load.

#### Function
Generates synthetic load on CPU cores to reach a target usage percentage.

#### Flags

##### `-l, --load string`
Target CPU load.
- **Default**: `0.8` (80%)
- **Range**: 0.0 to 1.0
- **Function**: Sets the target CPU utilization.

##### `-t, --time int`
Duration.
- **Default**: 0 (Infinite)
- **Unit**: Minutes
- **Function**: How long the simulation should run. 0 means run until manually stopped (Ctrl+C).

#### Examples
```bash
# 50% CPU load for 10 minutes
./shub runcpu -l 0.5 -t 10
```

---

### Command: `memory`

Simulate Memory usage.

#### Function
Allocates and holds a specified amount of RAM.

#### Flags

##### `-s, --size string`
Target Memory usage.
- **Default**: `0.8` (80%)
- **Range**: 0.0 to 1.0
- **Function**: Sets the target percentage of total system memory to occupy.

##### `-t, --time int`
Duration.
- **Unit**: Minutes
- **Function**: How long to hold the memory before releasing it.

#### Examples
```bash
# Consume 80% memory
./shub memory -s 0.8
```

---

### Command: `disk`

Simulate Disk usage.

#### Function
Creates files on the disk to consume storage space.

#### Flags

##### `-s, --size int`
File size.
- **Unit**: GB
- **Function**: The size of each file to create.

##### `-p, --path string`
Output path.
- **Function**: The directory where the files will be created.

##### `-t, --time int`
Duration.
- **Unit**: Minutes
- **Function**: (Note: For disk usage, this might refer to the duration to keep the files or a loop duration depending on implementation detail. Usually used to clean up after test).

#### Examples
```bash
# Create 1GB files in /tmp/test-disk
./shub disk -s 1 -p /tmp/test-disk
```
