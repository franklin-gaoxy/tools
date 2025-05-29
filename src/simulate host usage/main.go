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
		Args: cobra.MaximumNArgs(1),
	}

	// bond command
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(runCPU)
	rootCmd.AddCommand(SUMemory)

	// bond arge
	var cpuLoad string
	runCPU.Flags().StringVarP(&cpuLoad, "load", "l", "0.8", "cpu load, demo: 0.8(80%)")
	var RunTime int
	runCPU.Flags().IntVarP(&RunTime, "time", "t", 0, "run time, Unit:minutes, How long do you want it to run for.")
	var MemorySize string
	SUMemory.Flags().StringVarP(&MemorySize, "size", "s", "0.8", "memory size, demo: 0.8(80%)")
	SUMemory.Flags().IntVarP(&RunTime, "time", "t", 0, "run time, Unit:minutes, How long do you want it to run for.")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

/*
Command execution function code area
*/
