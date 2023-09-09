package ui

import (
	"github.com/euheimr/ghtop/devices"
	"github.com/rivo/tview"
	"github.com/shirou/gopsutil/v3/host"
	"strconv"
	"strings"
	"time"
)

const SysInfoLabel string = "[ System Info ]"
const (
	HostnameLabel     string = "Hostname:"
	SocketsCoresLabel        = "Sockets/Cores:"
	ThreadsLabel             = "Threads:"
	RefreshRateLabel         = "Refresh rate:"
	ProcessesLabel           = "Processes:"
	TickLabel                = "Tick:"
)

var tick string
var tickSymbols = [10]string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}

type SysInfo struct {
	Hostname     string
	SocketsCores string
	Threads      string
}

func formatLine(lineWidth int, title string, info string) string {
func tickCycleSymbol(tick string) string {
	if tick == "" {
		return tickSymbols[0]
	} else if tick == tickSymbols[0] {
		return tickSymbols[1]
	} else if tick == tickSymbols[1] {
		return tickSymbols[2]
	} else if tick == tickSymbols[2] {
		return tickSymbols[3]
	} else if tick == tickSymbols[3] {
		return tickSymbols[4]
	} else if tick == tickSymbols[4] {
		return tickSymbols[5]
	} else if tick == tickSymbols[5] {
		return tickSymbols[6]
	} else if tick == tickSymbols[6] {
		return tickSymbols[7]
	} else if tick == tickSymbols[7] {
		return tickSymbols[8]
	} else if tick == tickSymbols[8] {
		return tickSymbols[9]
	} else if tick == tickSymbols[9] {
		return tickSymbols[0]
	}
	return tick
}

func formatLine(lineWidth int, title string, info string) string {
	spaces := ""
	spacing := lineWidth - len(title+info)
	for i := 0; i < spacing; i++ {
		spaces += " "
	}
	return title + spaces + info
}

func UpdateSysInfoBox(app *tview.Application, sysInfoBox *tview.TextView,
	update time.Duration) {

	// Get Sysinfo data - this isn't in the for loop because this doesn't change
	//	during the lifetime of the program, thus we only get it once
	var hostInfo, _ = host.Info()
	// these variables grab info using functions defined in devices/cpu.go
	var cpuSockets = devices.CpuSockets()
	var cpuCores = devices.CpuCores()
	var cpuThreads = devices.CpuThreads()

	sysInfoBox.SetBorder(true).SetTitle(SysInfoLabel)

	sysInfo := &SysInfo{
		Hostname: strings.Split(hostInfo.Hostname, ".")[0],
		SocketsCores: strconv.FormatInt(int64(cpuSockets), 10) + "/" +
			strconv.FormatInt(int64(cpuCores)*int64(cpuSockets), 10),
		Threads: strconv.FormatInt(int64(cpuThreads), 10),
	}

	for {
		_, _, width, height := sysInfoBox.GetInnerRect()

		hostnameLine := formatLine(width, HostnameLabel, sysInfo.Hostname)
		socketsCoresLine := formatLine(width, SocketsCoresLabel, sysInfo.SocketsCores)
		threadsLine := formatLine(width, ThreadsLabel, sysInfo.Threads)
		refreshLine := formatLine(width, RefreshRateLabel,
			strconv.FormatInt(int64(update/time.Millisecond), 10)+"ms")

		refreshLine := formatLine(width, "Refresh rate:", strconv.FormatInt(int64(update/time.Millisecond), 10)+"ms")
		tick = tickCycleSymbol(tick)

		// we want the number of processes updated, unlike the rest of the
		//	system info, so we call host.Info() again to update the number
		//	of processes with each draw
		hostInfo, _ = host.Info()

		time.Sleep(update)
		app.QueueUpdateDraw(func() {

			procsCount := strconv.FormatInt(int64(hostInfo.Procs), 10)
			procsLine := formatLine(width, ProcessesLabel, procsCount)
			tickLine := formatLine(width, TickLabel, tick)

			sysInfoText := hostnameLine + "\n" +
				socketsCoresLine + "\n" +
				threadsLine + "\n" +
				procsLine + "\n" +
				"\n" +
				refreshLine + "\n" +
				tickLine + "\n" +
				strconv.FormatInt(int64(height), 10)

			sysInfoBox.SetText(sysInfoText)
		})
	}
}
