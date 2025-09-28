package server

import (
	"bwrs/tools"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/net"
	"k8s.io/klog"
	"time"
)

/*
get all information.
all func return data in JSON format
*/

func GetInfoCPU() []tools.CPU {
	var cpuInfoList []tools.CPU
	cpuInfo, err := cpu.Info()
	if err != nil {
		klog.Errorf("get CPU info failed: %v\n", err)
		return nil
	}

	for i, info := range cpuInfo {
		var cpuinformation tools.CPU
		cpuinformation.Number = i
		cpuinformation.ModelName = info.ModelName
		cpuinformation.Cores = info.Cores
		cpuinformation.Mhz = info.Mhz
		cpuinformation.CacheSize = info.CacheSize
		cpuinformation.Flags = info.Flags
		percent, _ := cpu.Percent(time.Second*3, true)
		cpuinformation.Percent = percent[i]

		cpuInfoList = append(cpuInfoList, cpuinformation)
	}

	return cpuInfoList
}

func GetInfoMemory() tools.MemoryInfo {
	vm, err := mem.VirtualMemory()
	if err != nil {
		return tools.MemoryInfo{}
	}

	return tools.MemoryInfo{
		Total:       vm.Total / 1024 / 1024,
		Available:   vm.Available / 1024 / 1024,
		Used:        vm.Used / 1024 / 1024,
		UsedPercent: vm.UsedPercent,
		Free:        vm.Free / 1024 / 1024,
		Cached:      vm.Cached / 1024 / 1024,
	}
}

func GetInfoDisk() []tools.DiskInfo {
	partitions, err := disk.Partitions(true)
	if err != nil {
		return nil
	}

	var diskInfos []tools.DiskInfo
	for _, partition := range partitions {
		usage, err := disk.Usage(partition.Mountpoint)
		if err != nil {
			continue // 跳过无法获取使用情况的分区
		}

		diskInfos = append(diskInfos, tools.DiskInfo{
			Total:       usage.Total / 1024 / 1024,
			Available:   usage.Free / 1024 / 1024,
			Used:        usage.Used / 1024 / 1024,
			UsedPercent: usage.UsedPercent,
			Free:        usage.Free / 1024 / 1024,
			Name:        partition.Device,
			Mountpoint:  partition.Mountpoint,
			Type:        partition.Fstype,
		})
	}

	return diskInfos
}

func GetInfoNetwork() []tools.NetworkInfo {
	interfaces, _ := net.Interfaces()
	var networkInfos []tools.NetworkInfo

	for _, iface := range interfaces {
		var addresses []string
		for _, addr := range iface.Addrs {
			addresses = append(addresses, addr.Addr)
		}
		networkInfos = append(networkInfos, tools.NetworkInfo{
			Name:    iface.Name,
			Address: addresses,
		})
	}

	return networkInfos
}

func GetInfoNode() tools.NodeInfo {
	info, _ := host.Info()
	return tools.NodeInfo{
		Hostname:        info.Hostname,
		OS:              info.OS,
		Platform:        info.Platform,
		PlatformVersion: info.PlatformVersion,
		KernelVersion:   info.KernelVersion,
		Arch:            info.KernelArch,
	}
}
