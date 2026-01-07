# SSH Execution Tool

A powerful tool to execute commands on multiple remote hosts via SSH. It supports parallel, serial, and batch execution modes.

## Build

```shell
go build -o ssh-tool .
# or
make
```

## Configuration (Nodelist)

Create a `nodelist` file to define your hosts.

### Format

The `nodelist` file supports two formats.

**1. Legacy Format**

Simple space-separated values. If port is omitted, it defaults to 22.

```
# IP[:Port] User Password
192.168.1.100 root password123
192.168.1.101:2222 admin secret
```

**2. New INI-style Format**

Supports groups and key-value pairs. Field order is flexible.
- If `user` is omitted, defaults to `root`.
- If `port` is omitted, defaults to `22`.
- **Authentication**: Supports both `password` and `ssh_key`. If both are provided, `ssh_key` is prioritized.

```ini
[web-servers]
address=192.168.1.10 user=root password=pass
address=192.168.1.11 user=admin port=2022 password=secret

[db-servers]
address=10.0.0.5 port=22 ssh_key="/path/to/id_rsa"
address=10.0.0.6 port=22 password=dbpass
```

## Usage

```shell
./ssh-tool [command] [flags]
```

### Global Flags

*   `-f, --file string`: Path to the hosts file (default `"./nodelist"`).
*   `-v, --verbose int`: Log verbosity level (default `0`).
*   `-g, --group string`: Target group filter. Only hosts in this group will be targeted (default `""`).

### Commands

#### 1. Parallel Execution (Default)
Run commands on all hosts simultaneously.

```shell
./ssh-tool parallel <command> [flags]
# or simply (root command defaults to parallel)
./ssh-tool <command> [flags]
```

**Flags:**
*   `-g, --group string`: Target group (default `""`).

**Example:**
```shell
./ssh-tool "uptime" -g web-servers
```

#### 2. Serial Execution
Run commands on hosts one by one. Useful for debugging or strict ordering.

```shell
./ssh-tool serial <command> [flags]
```

**Flags:**
*   `-g, --group string`: Target group (default `""`).

**Example:**
```shell
./ssh-tool serial "date"
```

#### 3. Batch Execution
Run commands in batches to control concurrency and load.

```shell
./ssh-tool batch <command> [flags]
```

**Flags:**
*   `-n, --number int`: Batch size (default `5`).
*   `-t, --threadpool bool`: Use thread pool mode (default `false`). If true, maintains `n` concurrent connections; otherwise, executes in batches of `n`.
*   `-g, --group string`: Target group (default `""`).

**Example:**
```shell
# Run in batches of 10 (wait for all 10 to finish before next batch)
./ssh-tool batch "yum update -y" -n 10 -g db-servers

# Run with thread pool of 10 (keep 10 running constantly)
./ssh-tool batch "sleep 10" -n 10 -t -g db-servers
```

#### 4. SCP File Transfer
Copy files or directories to remote hosts.

```shell
./ssh-tool scp <src> <dest> [flags]
```

**Flags:**
*   `-n, --number int`: Concurrent limit (default `0` for unlimited).
*   `-g, --group string`: Target group (default `""`).

**Example:**
```shell
# Copy local file to remote /tmp
./ssh-tool scp ./config.json /tmp/config.json -g web-servers

# Copy directory recursively
./ssh-tool scp ./dist /var/www/html -g web-servers
```

#### 5. Test Connection
Test SSH connectivity and authentication for all (or filtered) hosts without running commands.

```shell
./ssh-tool test [flags]
```

**Flags:**
*   `-g, --group string`: Target group (default `""`).

#### 6. Version
Show version information.

```shell
./ssh-tool version
```
