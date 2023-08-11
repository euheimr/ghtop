package devices

import (
	"github.com/shirou/gopsutil/v3/cpu"
	"log"
	"strings"
)

func getCpuInfo() []cpu.InfoStat {
	info, err := cpu.Info()
	if err != nil {
		log.Fatalf("Could not get CPU Info!")
	}
	return info
}

func CpuModelName() string {
	info := getCpuInfo()
	return strings.Replace(info[0].ModelName, "CPU ", "", -1)
}

func CpuSockets() int {
	info := getCpuInfo()
	return len(info)
}

func CpuCores() int {
	cores, err := cpu.Counts(false)

	if err != nil {
		log.Fatalf("Could not get CPU Cores count!")
	}
	return cores
}

func CpuThreads() int {
	threads, err := cpu.Counts(true)
	if err != nil {
		log.Fatalf("Could not get CPU Threads (logical) count!")
	}
	return threads
}
