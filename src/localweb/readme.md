# Local Web Server

A simple static file server that exposes a local directory over HTTP. It allows browsing directories and downloading files.

## Build

```bash
# Using Make
make build

# Using Go directly
go build -o localweb main.go
```

## Usage

```bash
./localweb [flags]
```

---

### Global Flags

#### `-port string`
Service port.
- **Default**: `:8080`
- **Format**: `:PORT` (e.g., `:9090`)
- **Function**: Sets the TCP port the server listens on.

#### `-path string`
Root directory.
- **Default (Windows)**: `D:\`
- **Default (Linux/macOS)**: `./` (Current directory)
- **Function**: The local file system path that will be served as the root of the web server.

---

### Features

#### Directory Listing
When you access a URL that points to a directory, the server lists all files and subdirectories within it, allowing for easy navigation.

#### File Download
Clicking on a file in the list or accessing its URL directly initiates a download. The server sets `Content-Disposition: attachment` to force the browser to download the file instead of trying to display it.

#### Access Logs
The server prints a log line to stdout for every request received, including the timestamp and client IP address.

---

### Examples

#### Basic Usage
Serve the current directory on port 8080.
```bash
./localweb
```

#### Custom Port and Path
Serve the `/var/www/html` directory on port 9000.
```bash
./localweb -port :9000 -path /var/www/html
```

> **Note**: Access files via `http://localhost:8080/file/YOUR_FILE_PATH` (The code mounts the handler at `/`).
