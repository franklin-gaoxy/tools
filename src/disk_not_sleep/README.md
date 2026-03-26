# disk_not_sleep

通过周期性向目标磁盘写入一个极小文件，阻止外置盘/机械盘进入休眠。

## 构建

在当前目录执行：

```bash
go build -o disk_not_sleep .
```

Windows（在 Windows 上执行或交叉编译）：

```bash
GOOS=windows GOARCH=amd64 go build -o disk_not_sleep.exe .
```

## 使用

```bash
./disk_not_sleep [flags]
```

### Flags

#### `-f, --config`

指定配置文件路径。

- 不传 `-f`：只会查找当前工作目录下是否存在 `disk_not_sleep.yaml`；如果不存在，则使用内置默认配置。

#### `--tray`

是否以托盘模式运行（默认 `true`）。

- `--tray=true`：最小化到托盘运行，托盘菜单里可以点击“退出”。
- `--tray=false`：前台运行（方便在命令行里调试/观察日志）。

## 配置说明

配置文件为 YAML，顶层 key 为 `disk_not_sleep`：

```yaml
disk_not_sleep:
  tmp_file_path: "/path/to/disk"
  tmp_file_name: "disk_not_sleep"
  time_interval: "180s"
  log_level: "info"
  log_file_path: "/path/to/logs"
  log_file_prefix: "disk_not_sleep"
```

- `tmp_file_path`：要在哪个目录下创建/写入临时文件，不存在则创建。
- `tmp_file_name`：临时文件名，最终写入位置为 `tmp_file_path/tmp_file_name`。
- `time_interval`：间隔时间（Go `time.ParseDuration` 格式，如 `180s`、`5m`）。
- `log_level`：日志级别（`debug|info|warn|error`），用于设置 klog 的详细程度。
- `log_file_path`：日志输出目录；每次启动都会创建一个新日志文件。
- `log_file_prefix`：日志文件名前缀；文件名为 `<prefix>_<启动时间>.log`。

## 示例

### 1) 前台运行（推荐先这样验证配置）

```bash
./disk_not_sleep --tray=false -f ./disk_not_sleep.yaml
```

### 2) 使用默认配置文件名（不传 -f）

把配置文件放到当前目录并命名为 `disk_not_sleep.yaml`：

```bash
./disk_not_sleep --tray=false
```

### 3) 托盘运行

```bash
./disk_not_sleep -f ./disk_not_sleep.yaml
```
