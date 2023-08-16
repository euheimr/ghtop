package main

import (
	"github.com/euheimr/ghtop/devices"
	"github.com/euheimr/ghtop/utils"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/shirou/gopsutil/v3/host"
	"log"
	"strconv"
)

var cfg = util.ReadConfig()

var RefreshInterval = cfg.Duration("UpdateInterval") * time.Millisecond
var GroupProcesses = cfg.Bool("GroupProcesses")
var TempScale = cfg.String("TempScale")
var cpuCores = devices.CpuCores()
var cpuThreads = devices.CpuThreads()

var sysInfoData = map[string]string{
	"hostname": hostInfo.Hostname,
	"socket-cores": strconv.FormatInt(int64(cpuSockets), 10) + "/" +
		strconv.FormatInt(int64(cpuCores)*int64(cpuSockets), 10),
	"threads": strconv.FormatInt(int64(cpuThreads), 10),
}

var widgetLabels = map[string]string{
	"sys-info": "[System Info]",
	"cpu":      "[CPU: " + cpuModelName + "]",
	"procs":    "[Processes]",
	"cpu-temp": "[CPU Temp]",
	"mem":      "[Memory]",
	"net":      "[Network]",
	"gpu":      "[GPU: ]",
	"gpu-temp": "[GPU Temp]",
}

func main() {
	screen, err := tcell.NewScreen()
	if err != nil {
		log.Fatalf("[%s] Failed to initialize tcell.NewScreen(): %v",
			utils.GetFuncName(), err)
	}

	app := tview.NewApplication().SetScreen(screen).EnableMouse(true)

	sysInfoBox := tview.NewTextView()
	cpuBox := tview.NewBox()
	cpuTempBox := tview.NewBox()
	memBox := tview.NewBox()
	procsTbl := tview.NewTable().
		SetFixed(0, 6).
		SetSeparator(tview.BoxDrawingsLightVertical)
	//SetBordersColor(tcell.ColorYellow)
	netBox := tview.NewBox()
	gpuBox := tview.NewBox()
	gpuTempBox := tview.NewBox()

	flex := tview.NewFlex().SetDirection(tview.FlexRow)
	flex.SetBorderStyle(tcell.StyleDefault)

	go func() {
		app.QueueUpdateDraw(func() {
			//flex.Clear()
			//SYSINFO
			w := sysInfoBox.BatchWriter()
			defer w.Close()
			w.Clear()
			w.Write([]byte("Hostname: " + hostInfo.Hostname + "\n" +
				//"CPU Clk: " + strconv.FormatFloat(cpuInfo[0].Mhz/1000, 'f', -1, 64) + " GHz\n" +
				"Socket/Cores: " + sysInfoData["socket-cores"] + "\n" +
				"Threads: " + sysInfoData["threads"] + "\n" +
				//"Cache Size: " + strconv.FormatInt(int64(cpuInfo[0].CacheSize), 10) + "\n" +
				"Processes: " + strconv.FormatInt(int64(hostInfo.Procs), 10) + "\n"))

			// PROCESSES
			tblHeaderColor := tcell.ColorYellow
			tblHeaderAlign := tview.AlignCenter
			procsTbl.
				SetCell(0, 0, &tview.TableCell{
					Text:      "CNT",
					Color:     tblHeaderColor,
					Align:     tblHeaderAlign,
					Expansion: 3,
					MaxWidth:  4,
				}).
				SetCell(0, 1, &tview.TableCell{
					Text:      "USER",
					Color:     tblHeaderColor,
					Align:     tblHeaderAlign,
					Expansion: 6,
					MaxWidth:  12,
				}).
				SetCell(0, 2, &tview.TableCell{
					Text:      "Exec",
					Color:     tblHeaderColor,
					Align:     tblHeaderAlign,
					Expansion: 24,
					MaxWidth:  0,
				}).
				SetCell(0, 3, &tview.TableCell{
					Text:      "CPU%",
					Color:     tblHeaderColor,
					Align:     tblHeaderAlign,
					Expansion: 4,
					MaxWidth:  5,
				}).
				SetCell(0, 4, &tview.TableCell{
					Text:      "MEM%",
					Color:     tblHeaderColor,
					Align:     tblHeaderAlign,
					Expansion: 4,
					MaxWidth:  5,
				}).
				SetCell(0, 5, &tview.TableCell{
					Text:      "GPU%",
					Color:     tblHeaderColor,
					Align:     tblHeaderAlign,
					Expansion: 4,
					MaxWidth:  5,
				}) //.
			// todo: START DATA
			//SetCell(1, 0, &tview.TableCell{
			//	Text:  strconv.Itoa(colCount),
			//	Color: tcell.ColorWhite,
			//	Align: tview.AlignCenter,
			//})

			sysInfoBox.SetBorder(true).SetTitle(widgetLabels["sys-info"])
			cpuBox.SetBorder(true).SetTitle(widgetLabels["cpu"])
			cpuTempBox.SetBorder(true).SetTitle(widgetLabels["cpu-temp"])
			memBox.SetBorder(true).SetTitle(widgetLabels["mem"])
			procsTbl.SetBorder(true).SetTitle(widgetLabels["procs"])
			netBox.SetBorder(true).SetTitle(widgetLabels["net"])
			gpuBox.SetBorder(true).SetTitle(widgetLabels["gpu"])
			gpuTempBox.SetBorder(true).SetTitle(widgetLabels["gpu-temp"])

			flex.
				AddItem(tview.NewFlex().
					AddItem(sysInfoBox, 0, 3, false).
					AddItem(cpuBox, 0, 8, false).
					AddItem(memBox, 0, 4, false),
					0, 1, false).
				AddItem(tview.NewFlex().
					AddItem(procsTbl, 0, 2, false).
					AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
						AddItem(cpuTempBox, 0, 1, false).
						AddItem(netBox, 0, 1, false),
						0, 1, false).
					AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
						AddItem(gpuBox, 0, 1, false).
						AddItem(gpuTempBox, 0, 1, false),
						0, 1, false),
					0, 2, false)
		})
	}()

	if err := app.SetRoot(flex, true).SetFocus(flex).Run(); err != nil {
		panic(err)
	}
}
