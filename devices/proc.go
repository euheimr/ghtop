package devices

import (
	"github.com/euheimr/ghtop/util"
	"github.com/shirou/gopsutil/v3/process"
	"log"
	"strings"
)

type Process struct {
	Pid  int
	User string
	Name string
	Cpu  float64
	Mem  float64
}

var _procs []Process

func groupProcs(procs []Process) []Process {
	// Create a map of unique processes and use Pid as a counter of the same
	//	processes. CPU and Memory are added for each process with the same name
	uniqueProcsMap := make(map[string]Process)
	for _, proc := range procs {

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
	_groupProcs := make([]Process, len(uniqueProcsMap))
	i := 0
	for _, val := range uniqueProcsMap {
		_groupProcs[i] = val
		i++
	}
	// then return the new []Process struct data
	return _groupProcs

}

func GetProcs(group bool) ([]Process, error) {
	// GET PROCESSES
	processes, err := process.Processes()
	if err != nil {
		log.Fatal(util.GetFuncName(), err.Error())
	}

	_procs = make([]Process, len(processes))
	for i, proc := range processes {
		pid := proc.Pid
		user, err := proc.Username()
		if err != nil {
			//log.Println(util.GetFuncName(), err.Error())
			continue
		}
		// Windows only - Cut out the user Group names
		_, user, _ = strings.Cut(user, "\\")
		if err != nil {
			//log.Println(util.GetFuncName(), err.Error())
			continue
		}
		name, err := proc.Name()
		name, _, _ = strings.Cut(name, ".exe")
		if err != nil {
			//log.Println(util.GetFuncName(), err.Error())
			continue
		}
		cpu, err := proc.CPUPercent()
		if err != nil {
			//log.Println(util.GetFuncName(), err.Error())
			continue
		}
		mem, err := proc.MemoryPercent()
		if err != nil {
			//log.Println(util.GetFuncName(), err.Error())
			continue
		}
		if int(pid) > 0 {
			_procs[i] = Process{
				Pid:  int(pid),
				User: user,
				Name: name,
				Cpu:  cpu,
				Mem:  float64(mem),
			}
		}
	}
	if group {
		_procs = groupProcs(_procs)
	}

	return _procs, nil
}
