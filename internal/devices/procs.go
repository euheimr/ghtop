package devices

import (
	"fmt"
	"github.com/euheimr/ghtop/internal"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/process"
	"sort"
	"strings"
)

type Process struct {
	Pid  int
	User string
	Name string
	Cpu  float64
	Mem  float64
}

const (
	Pid  int = 0
	User     = 1
	Exec     = 2
	Cpu      = 3
	Mem      = 4
	//Gpu      = 5
)

var Procs []Process
var GroupedProcs []Process

func SortProcs(processes []Process, sortColumn int, sortDescending bool) ([]Process, error) {
	sort.Slice(processes, func(i, j int) bool {
		switch sortDescending {
		case true:
			switch sortColumn {
			case Pid:
				return processes[i].Pid > processes[j].Pid
			case User:
				return processes[i].User > processes[j].User
			case Exec:
				return processes[i].Name > processes[j].Name
			case Cpu:
				return processes[i].Cpu > processes[j].Cpu
			case Mem:
				return processes[i].Mem > processes[j].Mem
			}
		case false:
			switch sortColumn {
			case Pid:
				return processes[i].Pid < processes[j].Pid
			case User:
				return processes[i].User < processes[j].User
			case Exec:
				return processes[i].Name < processes[j].Name
			case Cpu:
				return processes[i].Cpu < processes[j].Cpu
			case Mem:
				return processes[i].Mem < processes[j].Mem
			}
		}
		// Default sort Descending by CPU
		return processes[i].Cpu > processes[j].Cpu
	})
	return processes, nil
}

func GetProcs(group bool) ([]Process, error) {
	var (
		pid  int32
		usr  string
		exec string
		cpu  float64
		mem  float32
	)

	processes, _ := process.Processes()
	Procs = make([]Process, len(processes))

	for i, proc := range processes {
		pid = proc.Pid
		usr, _ = proc.Username()
		// Windows only - Cut out the user Group names
		if strings.Contains(usr, "\\") {
			_, usr, _ = strings.Cut(usr, "\\")
		}

		exec, _ = proc.Name()
		exec, _, _ = strings.Cut(exec, ".exe")

		cpu, _ = proc.CPUPercent()

		// !!TODO: this seems terrible and like a workaround.. because it is one
		var totCpuCores int
		//var totCpuThreads int
		for socket := range CpuInfo {
			// a hacky way to get the numbers more correct? divide by total cpu cores?
			totCpuCores += CpuInfo[socket].Cores
			cpu = cpu / float64(totCpuCores)
			//totCpuThreads += CpuInfo[socket].Threads
			//cpu = cpu / float64(totCpuThreads)
		}

		mem, _ = proc.MemoryPercent()

		//u, _ := user.Current()

		if internal.Config.ShowOnlyUserProcesses {
			if usr == SystemInfo.User {
				Procs[i] = Process{
					Pid:  int(pid),
					User: usr,
					Name: exec,
					Cpu:  float64(cpu),
					Mem:  float64(mem),
				}
			}
		} else {
			Procs[i] = Process{
				Pid:  int(pid),
				User: usr,
				Name: exec,
				Cpu:  float64(cpu),
				Mem:  float64(mem),
			}
		}
	}

	if group {
		// Create a map of unique processes and use Pid as a counter of the same
		//	processes. CPU and Memory are added for each process with the same name
		var uniqueProcsMap = make(map[string]Process)

		for _, proc := range Procs {
			if val, ok := uniqueProcsMap[proc.Name]; ok {
				uniqueProcsMap[proc.Name] = Process{
					Pid:  val.Pid + 1,
					User: val.User,
					Name: proc.Name,
					Cpu:  val.Cpu + proc.Cpu,
					Mem:  val.Mem + proc.Mem,
				}
			} else {
				uniqueProcsMap[proc.Name] = Process{
					Pid:  1,
					User: proc.User,
					Name: proc.Name,
					Cpu:  proc.Cpu,
					Mem:  proc.Mem,
				}
			}
		}
		GroupedProcs = make([]Process, len(uniqueProcsMap))
		i := 0
		for _, val := range uniqueProcsMap {
			GroupedProcs[i] = val
			i++
		}
		// then return the new []Process struct data
		return GroupedProcs, nil
	}
	return Procs, nil
}
