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

### Flags

- `-port string`: Service port (format `:8080`). Default is `:8080`.
- `-path string`: The local directory path to serve.
  - Default on Windows: `D:\`
  - Default on Linux/macOS: `./` (Current directory)

### Features

- **Directory Listing**: Browses subdirectories.
- **File Download**: direct download for files.
- **Access Logs**: Prints access logs to stdout.

### Example

```bash
# Serve the current directory on port 8080
./localweb

# Serve a specific path on port 9000
./localweb -port :9000 -path /var/www/html
```

> **Note**: Access files via `http://localhost:8080/file/YOUR_FILE_PATH`.
