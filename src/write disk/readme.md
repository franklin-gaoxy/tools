# Write disk

测试硬盘写入量，可以指定写入位置、单个文件大小、写入文件数量等。

## build

```shell
go build -o writedisk main.go
```

## demo
```shell
MP990K24LJ:write disk root# ./writedisk -n 10 -f testwrite -p . -s 1
开始创建 10 个文件，每个大小 1G，保存到 .
当前时间: 2025-09-16 16:17:02 本次耗时: 0 s已创建: testwrite_1
当前时间: 2025-09-16 16:17:02 本次耗时: 1 s已创建: testwrite_2
当前时间: 2025-09-16 16:17:03 本次耗时: 0 s已创建: testwrite_3
当前时间: 2025-09-16 16:17:03 本次耗时: 0 s已创建: testwrite_4
当前时间: 2025-09-16 16:17:03 本次耗时: 0 s已创建: testwrite_5
当前时间: 2025-09-16 16:17:03 本次耗时: 1 s已创建: testwrite_6
当前时间: 2025-09-16 16:17:04 本次耗时: 0 s已创建: testwrite_7
当前时间: 2025-09-16 16:17:04 本次耗时: 0 s已创建: testwrite_8
当前时间: 2025-09-16 16:17:04 本次耗时: 1 s已创建: testwrite_9
当前时间: 2025-09-16 16:17:05 本次耗时: 0 s已创建: testwrite_10
全部文件创建完成！
```

## help
```shell
MP990K24LJ:write disk root# ./writedisk -h
File generator with Cobra

Usage:
  filegen [flags]

Flags:
      --fast            启用快速预分配模式 (使用 Truncate)
  -h, --help            help for filegen
  -n, --number int      创建文件数量
  -p, --path string     文件保存路径
  -f, --prefix string   文件名前缀 (默认: file)
  -s, --size int        单个文件大小 (单位: G)
```