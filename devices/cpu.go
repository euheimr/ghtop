package devices

import (
	"github.com/euheimr/ghtop/util"
	"github.com/shirou/gopsutil/v3/cpu"
	"log"
	"strings"
)

var cpuInfo []cpu.InfoStat

func init() {
	cpuInfo = getCpuInfo()
}

func getCpuInfo() []cpu.InfoStat {
	info, err := cpu.Info()
	if err != nil {
		log.Fatal(util.GetFuncName(), "Could not get cpu.Info() - ", err.Error())
	}
	return info
}

func CpuModelName() string {
	// return the model name but remove `CPU ` from the name because it's redundant
	return strings.Replace(cpuInfo[0].ModelName, "CPU ", "", -1)
}

func CpuSockets() int {
	return len(cpuInfo)
}

func CpuCores() int {
	cores, err := cpu.Counts(false)
	if err != nil {
		log.Fatal(util.GetFuncName(), err.Error())
	}
	return cores
}

func CpuThreads() int {
	threads, err := cpu.Counts(true)
	if err != nil {
		log.Fatal(util.GetFuncName(), err.Error())
	}
	return threads
}
