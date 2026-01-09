# Printline

A CLI tool to generate and print separator lines or centered text in the terminal. Useful for scripts and logs formatting.

## Build

```bash
# Using Make
make build

# Using Go directly
go build -o printline main.go
```

## Usage

```bash
./printline [command] [flags]
```

### Commands

#### 1. `line`
Print a line of specified characters.

- `-l, --lines int`: Number of lines to print (default 1).
- `-b, --blank row string`: Print blank rows around the line (`y` or `n`, default `n`).

**Example:**
```bash
# Print a line of '='
./printline line

# Print a line of '*'
./printline line "*"
```

#### 2. `center`
Print text centered in the terminal, padded with symbols.

- `-s, --symbol string`: Specify the padding symbol (default `=`).
- `-p, --print info string`: Print debug length info (`y` to enable).
- `-b, --blank row string`: Print blank rows around (`y` or `n`, default `n`).

**Example:**
```bash
# Print "HELLO" centered with '='
./printline center "HELLO"

# Print "WARNING" centered with '!'
./printline center "WARNING" -s "!"
```

#### 3. `version`
Print version information.

#### 4. `completely-center`
*(Not Implemented)* Reserved for future completely centered printing features.
