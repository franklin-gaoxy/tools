package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"test/src"
)

/*
run cpu
*/

func main() {
	BondCobra()
}

/*
Cobra
*/

// BondCobra cobra command func
func BondCobra() {
	rootCmd := &cobra.Command{
		Use:   "run cpu",
		Short: "run cpu",
	}

	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "print version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("v0.1.0")
		},
		Args: cobra.MaximumNArgs(0),
	}

	// 沾满CPU
	runCPU := &cobra.Command{
		Use:   "runcpu",
		Short: "run cpu",
		Run: func(cmd *cobra.Command, args []string) {
			src.RunCPULoad(cmd, args)
		},
		Args: cobra.MaximumNArgs(0),
	}

	// 模拟使用内存
	SUMemory := &cobra.Command{
		Use:   "memory",
		Short: "simulate use memory",
		Run: func(cmd *cobra.Command, args []string) {
			src.SimulateUseMemory(cmd, args)
		},
		Args: cobra.MaximumNArgs(0),
	}

	// simulate disk
	SUDisk := &cobra.Command{
		Use:   "disk",
		Short: "simulate use disk",
		RunE:  src.SimulateUseDisk,
		Args:  cobra.MaximumNArgs(0),
	}

	// bond command
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(runCPU)
	rootCmd.AddCommand(SUMemory)
	rootCmd.AddCommand(SUDisk)

	// bond arge
	var cpuLoad string
	runCPU.Flags().StringVarP(&cpuLoad, "load", "l", "0.8", "cpu load, demo: 0.8(80%)")
	var RunTime int
	runCPU.Flags().IntVarP(&RunTime, "time", "t", 0, "run time, Unit:minutes, How long do you want it to run for.")
	var MemorySize string
	SUMemory.Flags().StringVarP(&MemorySize, "size", "s", "0.8", "memory size, demo: 0.8(80%)")
	SUMemory.Flags().IntVarP(&RunTime, "time", "t", 0, "run time, Unit:minutes, How long do you want it to run for.")
	// disk
	SUDisk.Flags().IntP("size", "s", 2, "Size of each file in GB")
	SUDisk.Flags().StringP("path", "p", "", "Directory path to write files")
	SUDisk.Flags().IntP("time", "t", 0, "Duration in minutes (0 for infinite)")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

/*
Command execution function code area
*/
