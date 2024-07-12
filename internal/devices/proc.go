package devices

import (
	"github.com/shirou/gopsutil/v4/host"
	"github.com/shirou/gopsutil/v4/process"
	"log/slog"
	"sort"
	"strings"
	"time"
)

type Process struct {
	Pid  int
	User string
	Name string
	Cpu  float64
	Mem  float64
	//Gpu float64
}

const (
	Pid = iota
	User
	Exec
	Cpu
	Mem
	//Gpu
)

var Processes []Process
var ProcsGrouped []Process
var HostInfo host.InfoStat

func sortProcesses(processes []Process, sortColumn int, sortDescending bool) ([]Process, error) {
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

func GetProcesses(sortColumn int, sortDescending bool, groupProcesses bool, update time.Duration) ([]Process, error) {
	procs, _ := process.Processes()

	Processes = make([]Process, len(procs))
	for i, proc := range procs {
		usr, _ := proc.Username()
		if strings.Contains(usr, "\\") {
			// windows only - cut out user Group names from `usr`
			_, usr, _ = strings.Cut(usr, "\\")
		}
		exec, _ := proc.Name()
		if strings.Contains(exec, ".exe") {
			// windows only... formatting the process name by removing `.exe` ...
			exec, _, _ = strings.Cut(exec, ".exe")
		}

		cpu, _ := proc.CPUPercent()
		mem, _ := proc.MemoryPercent()

		Processes[i] = Process{
			Pid:  int(proc.Pid),
			User: usr,
			Name: exec,
			Cpu:  cpu,
			Mem:  float64(mem),
		}
	}

	if !groupProcesses {
		Processes, _ = sortProcesses(Processes, sortColumn, sortDescending)
		return Processes, nil

	} else {
		// If we wanted the processes grouped, then create a map of unique processes and use
		//	Pid as a counter instead of an identifier.
		var uniqueProcesses = make(map[string]Process)

		for _, proc := range Processes {
			if val, ok := uniqueProcesses[proc.Name]; ok {
				// CPU and Memory are added together for each process with the same name
				uniqueProcesses[proc.Name] = Process{
					Pid:  val.Pid + 1,
					User: val.User,
					Name: proc.Name,
					Cpu:  val.Cpu + proc.Cpu,
					Mem:  val.Mem + proc.Mem,
				}
			} else {
				uniqueProcesses[proc.Name] = Process{
					Pid:  1,
					User: proc.User,
					Name: proc.Name,
					Cpu:  proc.Cpu,
					Mem:  proc.Mem,
				}
			}
		}

		ProcsGrouped = make([]Process, len(uniqueProcesses))
		i := 0
		for _, val := range uniqueProcesses {
			ProcsGrouped[i] = val
			i++
		}
		ProcsGrouped, _ = sortProcesses(ProcsGrouped, sortColumn, sortDescending)
		return ProcsGrouped, nil

	}
}

func GetProcessesCount() (uint64, error) {
	info, err := host.Info()
	if err != nil {
		slog.Error("Could not get host.Info() !!" + err.Error())
		return 0, err
	}
	return info.Procs, nil
}
