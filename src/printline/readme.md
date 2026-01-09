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

---

### Command: `line`

Print a line of specified characters.

#### Function
Repeats a character (or string) across the entire width of the terminal.

#### Flags

##### `-l, --lines int`
Number of lines.
- **Default**: 1
- **Function**: Specifies how many times to repeat the line vertically.

##### `-b, --blank row string`
Print blank rows.
- **Values**: `y` (yes), `n` (no)
- **Default**: `n`
- **Function**: If set to `y`, adds an empty line before and after the printed line for spacing.

#### Examples
```bash
# Print a line of '='
./printline line

# Print a line of '*'
./printline line "*"
```

---

### Command: `center`

Print text centered in the terminal.

#### Function
Calculates the terminal width and prints the provided text in the middle, padding the left and right with a specified symbol.

#### Flags

##### `-s, --symbol string`
Padding symbol.
- **Default**: `=`
- **Function**: The character used to fill the empty space around the centered text.

##### `-p, --print info string`
Print debug info.
- **Values**: `y`, `n`
- **Function**: If set to `y`, prints details about terminal width and string length calculations.

##### `-b, --blank row string`
Print blank rows.
- **Values**: `y`, `n`
- **Default**: `n`
- **Function**: Adds vertical padding (empty lines) around the output.

#### Examples
```bash
# Print "HELLO" centered with '='
./printline center "HELLO"

# Print "WARNING" centered with '!'
./printline center "WARNING" -s "!"
```

---

### Command: `version`

Print version information.

#### Function
Displays the current version of the tool.

---

### Command: `completely-center`

*(Reserved)*

#### Function
Currently not implemented. Reserved for future features to center text both vertically and horizontally.
