# Many Host Command (SSH Tool)

A powerful SSH automation tool to execute commands or copy files across multiple hosts in parallel, serially, or in batches.

## Build

```bash
# Using Make
make build

# Using Go directly
go build -o mhc main.go
```

## Configuration (Node List)

Create a host list file (default: `./nodelist`). The format supports INI-style groups.

**Example `nodelist`:**

```ini
[web-servers]
address=192.168.1.101 user=admin password="your_password"
address=192.168.1.102 user=root ssh_key="/path/to/private_key"

[db-servers]
address=10.0.0.5 user=dbadmin password="db_password"
```

> **Authentication Order**: SSH Key > Password. If `ssh_key` is provided, it is used first. If `password` is provided, it is used as a fallback or if no key is specified.

## Usage

```bash
./mhc [command] [flags]
```

### Global Flags

- `-f, --file string`: Path to the host list file (default `./nodelist`).
- `-g, --group string`: Target a specific group defined in the nodelist (default selects all hosts).
- `-v, --verbose int`: Log verbosity level (e.g., `-v=2` for debug logs).

### Commands

#### 1. Parallel Execution (`parallel`)
Run commands on all targeted hosts simultaneously.

```bash
# Run command on all hosts
./mhc parallel "uname -a"

# Run only on web-servers group
./mhc parallel "systemctl status nginx" -g web-servers
```

#### 2. Serial Execution (`serial`)
Run commands one by one (useful for rolling updates or debugging).

```bash
./mhc serial "uptime"
```

#### 3. Batch Execution (`batch`)
Run commands in batches to control concurrency and load.

- `-n, --number int`: Batch size (default 5).
- `-t, --threadpool bool`: Use thread pool mode (default false).

```bash
# Run on 10 hosts at a time
./mhc batch "yum update -y" -n 10

# Use thread pool for better resource management
./mhc batch "sleep 5" -n 5 -t
```

#### 4. SCP (File Copy) (`scp`)
Copy files or directories to remote hosts.

- `-n, --number int`: Concurrency limit for SCP (0 for unlimited).

```bash
# Syntax: ./mhc scp [src] [dest]

# Copy a file to all hosts
./mhc scp ./app.conf /etc/app/app.conf

# Copy directory recursively
./mhc scp ./html/ /var/www/html/

# Copy with concurrency limit (e.g., 5 at a time)
./mhc scp ./large-file /tmp/ -n 5
```

#### 5. Connectivity Test (`test`)
Test SSH connectivity to all targeted hosts without running commands.

```bash
./mhc test
./mhc test -g db-servers
```

#### 6. Version (`version`)
Print tool version.

```bash
./mhc version
```
