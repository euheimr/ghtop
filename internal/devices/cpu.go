package devices

import (
	"github.com/euheimr/ghtop/internal/app/common"
	"github.com/shirou/gopsutil/v3/cpu"
	"log"
	"strings"
)

type Info struct {
	Cores     int
	CoreId    string
	ModelName string
	Threads   int
}

type Data struct {
	Timestamp string
	CoreId    string
	Percent   float64
}

const MAX_DATA = 5

var CpuInfo map[int]*Info
var CpuData map[int]*Data

func init() {
	// We get CpuInfo only once because the hardware doesn't change from the
	//	start of ghtop's execution
	if CpuInfo == nil {
		info, err := cpu.Info()
		if err != nil {
			log.Fatal(common.GetFuncName(), "Could not get cpu.Info() - ", err.Error())
		}
		cores, _ := cpu.Counts(false)

		CpuInfo = make(map[int]*Info, len(info))

		for socket := range info {
			CpuInfo[socket] = &Info{
				Cores:     cores,
				CoreId:    info[socket].CoreID,
				ModelName: strings.Replace(info[socket].ModelName, "CPU ", "", -1),
				Threads:   int(info[socket].Cores),
			}
		}
	}

	// setup initial CpuData that gets populated in the main Cpu box
	// CpuData on the other hand does get updated, unlike CpuInfo
	CpuData = make(map[int]*Data, MAX_DATA)
	for i := 0; i < MAX_DATA; i++ {
		CpuData[i] = &Data{
			Timestamp: "",
			CoreId:    CpuInfo[0].CoreId,
			Percent:   0.34,
		}
	}

}

func getCpuInfo() {
	if CpuInfo == nil {
		info, err := cpu.Info()
		if err != nil {
			log.Fatal(common.GetFuncName(), "Could not get cpu.Info() - ", err.Error())
		}
		cores, _ := cpu.Counts(false)

		CpuInfo = make(map[int]*Info, len(info))

		for socket := range info {
			CpuInfo[socket] = &Info{
				Cores:     cores,
				CoreId:    info[socket].CoreID,
				ModelName: strings.Replace(info[socket].ModelName, "CPU ", "", -1),
				Threads:   int(info[socket].Cores),
			}
		}
	}

}

func GetCpuData() {
	for {
		//time.Sleep(update)

	}
}
