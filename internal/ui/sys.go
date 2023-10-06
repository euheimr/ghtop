package ui

import (
	"github.com/euheimr/ghtop/internal/devices"
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

type SysInfo struct {
	Hostname     string
	SocketsCores string
	Threads      string
}

type Tick struct {
	SymbolIndex int
	Symbols     []string
}

var tick *Tick

func init() {
	tick = &Tick{
		SymbolIndex: 0,
		Symbols:     []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"},
		//Symbols: []string{"-", "-", "+", "*", "*", "/", "\\"},
	}
}

func incrementTickSymbol() int {
	for i, _ := range tick.Symbols {
		if tick.SymbolIndex == 0 {
			tick.SymbolIndex++
			return tick.SymbolIndex
		} else if tick.SymbolIndex == len(tick.Symbols)-1 {
			tick.SymbolIndex = 0
			return tick.SymbolIndex
		} else if tick.SymbolIndex == i {
			tick.SymbolIndex++
			return tick.SymbolIndex
		}
	}
	return tick.SymbolIndex
}

func formatLine(lineWidth int, title string, info string) string {
	spaces := ""
	spacing := lineWidth - len(title+info)
	for i := 0; i < spacing; i++ {
		spaces += " "
	}
	return title + spaces + info
}

func UpdateSysInfo(app *tview.Application, sysInfoBox *tview.TextView,
	update time.Duration) {

	// Get Sysinfo data - this isn't in the for loop because this doesn't change
	//	during the lifetime of the program, thus we only get it once
	var hostInfo, _ = host.Info()
	var cpus = devices.CpuInfo
	// these variables grab info using functions defined in devices/cpu.go

	sysInfoBox.SetBorder(true).SetTitle(SysInfoLabel)

	sysInfo := &SysInfo{
		Hostname: strings.Split(hostInfo.Hostname, ".")[0],
		SocketsCores: strconv.FormatInt(int64(len(cpus)), 10) + "/" +
			strconv.FormatInt(int64(cpus[0].Cores)*int64(len(cpus)), 10),
		Threads: strconv.FormatInt(int64(cpus[0].Threads), 10),
	}

	for {
		_, _, width, height := sysInfoBox.GetInnerRect()

		hostnameLine := formatLine(width, HostnameLabel, sysInfo.Hostname)
		socketsCoresLine := formatLine(width, SocketsCoresLabel, sysInfo.SocketsCores)
		threadsLine := formatLine(width, ThreadsLabel, sysInfo.Threads)
		refreshLine := formatLine(width, RefreshRateLabel,
			strconv.FormatInt(int64(update/time.Millisecond), 10)+"ms")

		// we want the number of processes updated, unlike the rest of the
		//	system info, so we call host.Info() again to update the number
		//	of processes with each draw
		hostInfo, _ = host.Info()

		time.Sleep(update)

		app.QueueUpdateDraw(func() {
			procsCount := strconv.FormatInt(int64(hostInfo.Procs), 10)
			procsLine := formatLine(width, ProcessesLabel, procsCount)
			tickLine := formatLine(width, TickLabel, tick.Symbols[tick.SymbolIndex])
			tick.SymbolIndex = incrementTickSymbol()

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
