# Simulate host usage

## CPU

Simulate host CPU usage and run the CPU to the desired percentage.

```shell
root@zijie:~/workspace/test# ./simulatehostusage runcpu -h
run cpu

Usage:
  run runcpu [flags]

Flags:
  -h, --help          help for runcpu
  -l, --load string   cpu load, demo: 0.8(80%) (default "0.8")
  -t, --time int      run time, Unit:minutes, How long do you want it to run for.
# -l Specify a float,demo: 0.8,This will increase CPU usage to 80%
# -t specify a int,demo 10,The program will run for about 10 minutes.
```

## memory

```shell
root@zijie:~/workspace/test# ./simulatehostusage memory -h
simulate use memory

Usage:
  run memory [flags]

Flags:
  -h, --help          help for memory
  -s, --size string   memory size, demo: 0.8(80%) (default "0.8")
  -t, --time int      run time, Unit:minutes, How long do you want it to run for.
# -t specify a int,demo 10,The program will run for about 10 minutes.
# -s specify the size.Supports floating point number methods such as: 0.7 or specifying a specific value such as: 500M 1G 5G, which will occupy a specified size of memory.
```

## disk and io

```shell
root@zijie:~/workspace/test# ./main disk 
Error: failed to create directory: mkdir : no such file or directory
Usage:
  run disk [flags]

Flags:
  -h, --help          help for disk
  -p, --path string   Directory path to write files
  -s, --size int      Size of each file in GB (default 2)
  -t, --time int      Duration in minutes (0 for infinite)

failed to create directory: mkdir : no such file or directory
root@zijie:~/workspace/test# 
```

> Can increase the number of disk writes
>
> At least the - p path parameter needs to be specified.
>
> Temporary files will be created in this directory and manually deleted after execution.