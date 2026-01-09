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

- `-p, --port int`: Port to listen on (default `8080`).
- `--v=int`: Log level verbosity for klog (e.g., `--v=2`).

### Commands

- `version`: Print the version information.
- `help`: Print help information.

### API Endpoints

Once started (e.g., on port 8080), you can access:

- `GET /version`: Check service status and version.
- `GET /all`: Get all system information (CPU, Mem, Disk, Net, Node).
- `GET /cpu`: Get CPU information.
- `GET /memory`: Get Memory information.
- `GET /disk`: Get Disk information.
- `GET /network`: Get Network information.
- `GET /node`: Get Node/OS information.

### Example

```bash
# Start server on port 9090 with debug logging
./hostinfo --port 9090 --v=2
```
