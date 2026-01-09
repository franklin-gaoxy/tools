# Host Info Tool

A lightweight host information collection server based on Gin framework. It provides RESTful APIs to retrieve system information such as CPU, Memory, Disk, Network, and Host details.

## Build

You can build the binary using `make` or `go build`.

```bash
# Using Make
make build

# Using Go directly
cd src && go build -o ../hostinfo main.go
```

## Usage

Run the server specifying the port.

```bash
./hostinfo [flags] [command]
```

### Global Flags

#### `-p, --port int`
Port to listen on.
- **Default**: `8080`
- **Function**: Sets the HTTP port for the web server.

#### `--v=int`
Log verbosity.
- **Function**: Sets the logging level for klog. Higher values output more debug logs.

---

### Command: `version`

Print version information.

#### Function
Outputs the current version of the application to the console.

---

### Command: `help`

Print help information.

#### Function
Displays a summary of available commands and flags.

---

### API Endpoints

Once the server is started, you can access the following HTTP endpoints.

#### `GET /version`
Check service status and version.
- **Returns**: JSON object with status and version.

#### `GET /all`
Get all system information.
- **Returns**: A comprehensive JSON containing CPU, Memory, Disk, Network, and Node info.

#### `GET /cpu`
Get CPU information.
- **Returns**: CPU model, cores, frequency, and usage statistics.

#### `GET /memory`
Get Memory information.
- **Returns**: Total memory, used memory, free memory, and usage percentage.

#### `GET /disk`
Get Disk information.
- **Returns**: List of partitions, mount points, total space, used space, and usage percentage.

#### `GET /network`
Get Network information.
- **Returns**: Network interfaces, IP addresses, and I/O statistics.

#### `GET /node`
Get Node/OS information.
- **Returns**: Hostname, OS type, kernel version, platform, and uptime.
