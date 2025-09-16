package main

/*
用于测试硬盘写入 测试命令 ./writedisk -s 10 -n 10 -p /xx/xx
*/

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/spf13/cobra"
)

func main() {
	var sizeGB int
	var numFiles int
	var path string
	var prefix string
	var fast bool

	var rootCmd = &cobra.Command{
		Use:   "filegen",
		Short: "File generator with Cobra",
		RunE: func(cmd *cobra.Command, args []string) error {
			if sizeGB <= 0 || numFiles <= 0 || path == "" {
				return fmt.Errorf("参数错误: -s (文件大小, G), -n (文件数量), -p (路径) 必须都指定且有效")
			}

			// 确保路径存在
			if err := os.MkdirAll(path, 0755); err != nil {
				return fmt.Errorf("创建路径失败: %v", err)
			}

			sizeBytes := int64(sizeGB) * 1024 * 1024 * 1024
			if prefix == "" {
				prefix = "file"
			}

			fmt.Printf("开始创建 %d 个文件，每个大小 %dG，保存到 %s\n", numFiles, sizeGB, path)
			if fast {
				fmt.Println("启用快速预分配模式 (Truncate)")
			}

			if fast {
				// 使用 Truncate 快速创建文件
				for i := 0; i < numFiles; i++ {
					filename := filepath.Join(path, prefix+"_"+strconv.Itoa(i+1))
					f, err := os.Create(filename)
					if err != nil {
						return fmt.Errorf("创建文件失败 %s: %v", filename, err)
					}
					if err := f.Truncate(sizeBytes); err != nil {
						func() {
							if err := f.Close(); err != nil {
								panic("文件关闭失败！")
							}
						}()
						return fmt.Errorf("预分配文件失败 %s: %v", filename, err)
					}

					func() {
						if err := f.Close(); err != nil {
							panic("文件关闭失败！")
						}
					}()

					fmt.Printf("已创建完成: %s\n", filename)
				}
			} else {
				// 循环写满磁盘
				buf := make([]byte, 1024*1024) // 1MB 缓冲
				for i := 0; i < numFiles; i++ {
					// 获取时间戳
					timeStamp := time.Now().Unix()

					filename := filepath.Join(path, prefix+"_"+strconv.Itoa(i+1))
					f, err := os.Create(filename)
					if err != nil {
						return fmt.Errorf("创建文件失败 %s: %v", filename, err)
					}

					var written int64
					for written < sizeBytes {
						toWrite := sizeBytes - written
						if toWrite > int64(len(buf)) {
							toWrite = int64(len(buf))
						}
						n, err := f.Write(buf[:toWrite])
						if err != nil {
							func() {
								if err := f.Close(); err != nil {
									panic("文件关闭失败！")
								}
							}()
							return fmt.Errorf("写入文件失败 %s: %v", filename, err)
						}
						written += int64(n)
					}
					func() {
						if err := f.Close(); err != nil {
							panic("文件关闭失败！")
						}
					}()

					// 计算总耗时
					stopTimeStamp := time.Now().Unix()
					var timeConsuming string
					if (stopTimeStamp-timeStamp)/60 == 0 {
						timeConsuming = fmt.Sprintf("本次耗时: %d s", stopTimeStamp-timeStamp)
					} else {
						timeConsuming = fmt.Sprintf("本次耗时: %d m", (stopTimeStamp-timeStamp)/60)
					}

					fmt.Printf("当前时间: %s %s", time.Unix(timeStamp, 0).Format("2006-01-02 15:04:05"), timeConsuming)
					fmt.Printf("已创建: %s\n", filename)
				}
			}

			fmt.Println("全部文件创建完成！")
			return nil
		},
	}

	rootCmd.Flags().IntVarP(&sizeGB, "size", "s", 0, "单个文件大小 (单位: G)")
	rootCmd.Flags().IntVarP(&numFiles, "number", "n", 0, "创建文件数量")
	rootCmd.Flags().StringVarP(&path, "path", "p", "", "文件保存路径")
	rootCmd.Flags().StringVarP(&prefix, "prefix", "f", "", "文件名前缀 (默认: file)")
	rootCmd.Flags().BoolVar(&fast, "fast", false, "启用快速预分配模式 (使用 Truncate)")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println("错误:", err)
		os.Exit(1)
	}
}
