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

### Example `nodelist`

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

#### `-f, --file string`
Path to the host list file.
- **Default**: `./nodelist`
- **Function**: Specifies which configuration file to load host information from.

#### `-g, --group string`
Target host group.
- **Default**: "" (Selects all hosts)
- **Function**: Filters hosts by the group name defined in the nodelist (e.g., `[web-servers]`).

#### `-v, --verbose int`
Log verbosity level.
- **Default**: 0
- **Function**: Controls the detail level of logs. Higher values (e.g., 2) print more debug information.

---

### Command: `parallel`

Run commands on all targeted hosts simultaneously. This is the most common mode for broadcasting commands.

#### Usage
```bash
./mhc parallel [command]
```

#### Examples
```bash
# Run command on all hosts
./mhc parallel "uname -a"

# Run only on web-servers group
./mhc parallel "systemctl status nginx" -g web-servers
```

---

### Command: `serial`

Run commands one by one on targeted hosts.

#### Function
Executes the command on the first host, waits for it to finish, then proceeds to the next. Useful for rolling restarts or when you need to avoid hitting a shared resource simultaneously.

#### Usage
```bash
./mhc serial [command]
```

#### Examples
```bash
./mhc serial "uptime"
```

---

### Command: `batch`

Run commands in batches to control concurrency and load.

#### Function
Splits the host list into chunks (batches) and executes the command on one batch at a time.

#### Flags

##### `-n, --number int`
Batch size.
- **Default**: 5
- **Function**: Determines how many hosts are processed in one batch.

##### `-t, --threadpool bool`
Use thread pool mode.
- **Default**: false
- **Function**: If enabled, uses a worker pool model instead of simple batching. This is more efficient for large numbers of hosts as it keeps a constant number of concurrent connections active.

#### Usage
```bash
./mhc batch [command] [flags]
```

#### Examples
```bash
# Run on 10 hosts at a time
./mhc batch "yum update -y" -n 10

# Use thread pool for better resource management
./mhc batch "sleep 5" -n 5 -t
```

---

### Command: `scp`

Copy files or directories from the local machine to remote hosts.

#### Function
Uses SFTP to transfer files. Supports recursive directory copying.

#### Flags

##### `-n, --number int`
Concurrency limit.
- **Default**: 0 (Unlimited)
- **Function**: Limits the number of concurrent file transfers to avoid saturating network bandwidth.

#### Usage
```bash
./mhc scp [src] [dest]
```

#### Examples
```bash
# Copy a file to all hosts
./mhc scp ./app.conf /etc/app/app.conf

# Copy directory recursively
./mhc scp ./html/ /var/www/html/

# Copy with concurrency limit (e.g., 5 at a time)
./mhc scp ./large-file /tmp/ -n 5
```

---

### Command: `test`

Test SSH connectivity to hosts.

#### Function
Attempts to establish an SSH connection to each targeted host using the provided credentials. It does not execute any command but verifies if the host is reachable and credentials are correct.

#### Usage
```bash
./mhc test
```

---

### Command: `version`

Print tool version.

#### Function
Displays the current version of the binary.

#### Usage
```bash
./mhc version
```
