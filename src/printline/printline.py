import shutil
import sys

def PrintString(symbol):
    print(symbol * terminal_width)

def PrintCenter(symbol):
    DefaultSymbol = "="
    if len(sys.argv) == 4:
        DefaultSymbol = sys.argv[3]
    print(f' {symbol} '.center(terminal_width, DefaultSymbol))

# check
if len(sys.argv) < 3:
    print("[ERROR]: parameter error!")
    print("The first parameter is optional: [line] or [center].")
    print("The second parameter: The content to be output.")
    print("If using string mode, the third parameter can be selected.")
    print("demo: ./[file name] line -")
    print("demo: ./[file name] line --- :this will output three lines")
    print("demo: ./[file name] line '   ' :this will output three blank lines")
    print("demo: ./[file name] center 'hello world'")
    print("demo: ./[file name] center 'hello world' -")
    sys.exit(1)

# 获取终端宽度
terminal_width = shutil.get_terminal_size().columns
command = sys.argv[1]
symbol = sys.argv[2]

if command == "line":
    PrintString(symbol)
elif command == "center":
    PrintCenter(symbol)
else:
    print("[ERROR]: parameter error: please input [line] or [center]")
