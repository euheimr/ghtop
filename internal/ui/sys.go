package ui

import (
	"github.com/euheimr/ghtop/internal"
	"github.com/euheimr/ghtop/internal/devices"
	"github.com/rivo/tview"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const SysInfoLabel string = "[ System Info ]"

const (
	BLACK  = "[black]"
	BLUE   = "[blue]"
	GREEN  = "[green]"
	RED    = "[red]"
	WHITE  = "[white]"
	YELLOW = "[yellow]"
	GRAY   = "[gray]"
)

const (
	HostnameLabel     string = GRAY + "Hostname" + WHITE
	UserLabel                = GRAY + "User" + WHITE
	SocketsCoresLabel        = GRAY + "Soc/Cores" + WHITE
	ThreadsLabel             = GRAY + "Threads" + WHITE
	RefreshRateLabel         = YELLOW + "Refresh" + WHITE
	ProcessesLabel           = YELLOW + "Processes" + WHITE
	DebugLabel               = RED + "DEBUG" + WHITE
	TickLabel                = YELLOW + "Tick" + WHITE
	WidthHeightLabel         = "w*h"
)

type Tick struct {
	Index   int
	Symbols []string
}

// var SystemInfo *SysInfo
var tick *Tick

func init() {
	tick = &Tick{
		Index:   0,
		Symbols: []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"},
		//Symbols: []string{"-", "-", "+", "*", "*", "/", "\\"},
	}
}

func incrementTickSymbol() int {
	for i, _ := range tick.Symbols {
		if tick.Index == 0 {
			tick.Index++
			return tick.Index
		} else if tick.Index == len(tick.Symbols)-1 {
			tick.Index = 0
			return tick.Index
		} else if tick.Index == i {
			tick.Index++
			return tick.Index
		}
	}
	return tick.Index
}

func getSpacingOffset(str string) int {
	rgx, _ := regexp.Compile(`\[\w*]`)
	colorTags := rgx.FindAll([]byte(str), -1)
	if colorTags != nil {
		var s string
		for i := range colorTags {
			s = strings.ReplaceAll(str, string(colorTags[i]), "")
		}

		colorTags = rgx.FindAll([]byte(s), -1)
		for i := range colorTags {
			s = strings.ReplaceAll(s, string(colorTags[i]), "")
		}
		return len(s)
	}
	return len(str)
}

func formatLineSpacing(lineWidth int, title string, info string) string {
	var spaces string

	titleOffset := getSpacingOffset(title)
	infoOffset := getSpacingOffset(info)
	spacing := lineWidth - (titleOffset + infoOffset)

	for i := 0; i < spacing; i++ {
		spaces += "."
	}
	return title + spaces + info + "\n"
}

func UpdateSysInfo(app *tview.Application, sysInfo *tview.TextView,
	update time.Duration) {

	sysInfo.SetBorder(true).SetTitle(SysInfoLabel)
	sysInfo.SetDynamicColors(true)

	for {
		_, _, width, height := sysInfo.GetInnerRect()

		hostnameLine := formatLineSpacing(width, HostnameLabel, devices.SystemInfo.Hostname)
		userLine := formatLineSpacing(width, UserLabel, devices.SystemInfo.User)
		socketsCoresLine := formatLineSpacing(width, SocketsCoresLabel, devices.SystemInfo.SocketsCores)
		threadsLine := formatLineSpacing(width, ThreadsLabel, devices.SystemInfo.Threads)
		refreshLine := formatLineSpacing(width, RefreshRateLabel,
			strconv.FormatInt(int64(update/time.Millisecond), 10)+"ms")
		debugLine := formatLineSpacing(width, DebugLabel, strconv.FormatBool(internal.Config.Debug))
		widthHeightLine := formatLineSpacing(width, WidthHeightLabel,
			strconv.FormatInt(int64(width), 10)+"*"+
				strconv.FormatInt(int64(height), 10))

		var dividerLine string
		middle := (width / 2) - 1
		for i := 0; i <= (middle); i++ {
			dividerLine += " "
			if i == middle {
				dividerLine += "-\n"
			}
		}

		//procsCnt, _ := devices.GetProcsCount()
		//procsCount := strconv.FormatInt(procsCnt, 10)
		// Keep updating the procsCount
		//procsCntLine := formatLineSpacing(width, ProcessesLabel, procsCount)

		time.Sleep(update)
		app.QueueUpdateDraw(func() {
			tickLine := formatLineSpacing(width, TickLabel, tick.Symbols[tick.Index])
			tick.Index = incrementTickSymbol()

			sysInfoText := hostnameLine +
				userLine +
				socketsCoresLine +
				threadsLine +
				dividerLine +
				//procsCntLine +
				refreshLine +
				tickLine +
				dividerLine +
				widthHeightLine +
				debugLine

			sysInfo.SetText(sysInfoText)
		})
	}
}
