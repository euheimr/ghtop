package ui

import (
	"github.com/euheimr/ghtop/devices"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"math"
	"sort"
	"strconv"
	"time"
)

type ProcSortMethod string
type ProcSortDirection string

type TblHeaderNames struct {
	Pid  string
	Cnt  string
	User string
	Exec string
	Cpu  string
	Mem  string
	Gpu  string
}

type TblColMinWidth struct {
	Pid  int
	Cnt  int
	User int
	Exec int
	Cpu  int
	Mem  int
	Gpu  int
}

type TblColMaxWidth struct {
	Pid  int
	Cnt  int
	User int
	Exec int
	Cpu  int
	Mem  int
	Gpu  int
}

type TblHeaderRow struct {
	Text     string
	Expand   int
	MaxWidth int
}

const (
	SortPid  ProcSortMethod = "p"
	SortUser                = "u"
	SortExec                = "e"
	SortCpu                 = "c"
	SortMem                 = "m"
)

const (
	DownArrow string = "▼"
	UpArrow          = "▲"
)

const (
	SortAsc  ProcSortDirection = "a"
	SortDesc                   = "d"
)

var procBoxLabel = "[ Processes ]"
var SortMethod ProcSortMethod
var SortDirection ProcSortDirection

var tblHeaderNames = []TblHeaderNames{
	{
		Pid:  "PID",
		Cnt:  "CNT",
		User: "USER",
		Exec: "EXEC",
		Cpu:  "CPU%",
		Mem:  "MEM%",
		Gpu:  "GPU%",
	},
}

var tblColMinWidth = []TblColMinWidth{
	{
		Pid:  3,
		Cnt:  3,
		User: 6,
		Exec: 12,
		Cpu:  4,
		Mem:  4,
		Gpu:  4,
	},
}
var tblColMaxWidth = []TblColMaxWidth{
	{
		Pid:  4,
		Cnt:  4,
		User: 8,
		Exec: 16,
		Cpu:  4,
		Mem:  4,
		Gpu:  4,
	},
}

func roundFloat(float float64, precision int) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(float*ratio) / ratio
}

func formatFloat(float float64, precision int) string {
	roundedFloat := roundFloat(float, precision)
	return strconv.FormatFloat(roundedFloat, 'g', precision, 64)
}

func getHeaderNames(groupProcesses bool) []TblHeaderRow {
	header := []TblHeaderRow{
		{Text: tblHeaderNames[0].Pid, Expand: tblColMinWidth[0].Pid,
			MaxWidth: tblColMaxWidth[0].Pid},
		{Text: tblHeaderNames[0].User, Expand: tblColMinWidth[0].User,
			MaxWidth: tblColMaxWidth[0].User},
		{Text: tblHeaderNames[0].Exec, Expand: tblColMinWidth[0].Exec,
			MaxWidth: tblColMaxWidth[0].Exec},
		{Text: tblHeaderNames[0].Cpu, Expand: tblColMinWidth[0].Cpu,
			MaxWidth: tblColMaxWidth[0].Cpu},
		{Text: tblHeaderNames[0].Mem, Expand: tblColMinWidth[0].Mem,
			MaxWidth: tblColMaxWidth[0].Mem},
		//{Text: tblHeaderNames[0].Gpu, Expand: tblColMinWidth[0].Gpu,
		//	MaxWidth: tblColMaxWidth[0].Gpu},
	}
	if groupProcesses {
		header[0].Text = tblHeaderNames[0].Cnt
	}
	return header
}

func sortProcs(processes []devices.Process, sortDirection ProcSortDirection,
	sortMethod ProcSortMethod) []devices.Process {
	sort.Slice(processes, func(i, j int) bool {
		switch sortDirection {
		case SortAsc:
			switch sortMethod {
			case SortPid:
				return processes[i].Pid < processes[j].Pid
			case SortUser:
				return processes[i].User < processes[j].User
			case SortExec:
				return processes[i].Name < processes[j].Name
			case SortCpu:
				return processes[i].Cpu < processes[j].Cpu
			case SortMem:
				return processes[i].Mem < processes[j].Mem
			}
		case SortDesc:
			switch sortMethod {
			case SortPid:
				return processes[i].Pid > processes[j].Pid
			case SortUser:
				return processes[i].User > processes[j].User
			case SortExec:
				return processes[i].Name > processes[j].Name
			case SortCpu:
				return processes[i].Cpu > processes[j].Cpu
			case SortMem:
				return processes[i].Mem > processes[j].Mem
			}
		}
		// Default sort Descending by CPU
		return processes[i].Cpu > processes[j].Cpu
	})
	SortDirection = sortDirection
	SortMethod = sortMethod
	return processes
}

func UpdateProcBox(app *tview.Application, procsTbl *tview.Table,
	refresh time.Duration, groupProcesses bool) {

	tblHeaderTextColor := tcell.ColorYellow
	tblHeaderBkgdColor := tcell.ColorRed
	tblHeaderAlign := tview.AlignCenter

	headerText := getHeaderNames(groupProcesses)

	procsTbl.SetFixed(1, 0).
		SetSeparator(tview.BoxDrawingsLightVertical).
		SetSelectable(false, true).
		SetEvaluateAllRows(true).
		SetBorder(true).
		SetTitle(procBoxLabel)

	//// Construct the header row by column
	//for col := 0; col < len(headerText); col++ {
	//	procsTbl.SetCell(
	//		0, col, &tview.TableCell{
	//			Text:            headerText[col].Text,
	//			Expansion:       headerText[col].Expand,
	//			MaxWidth:        headerText[col].MaxWidth,
	//			Color:           tblHeaderTextColor,
	//			BackgroundColor: tblHeaderBkgdColor,
	//			Align:           tblHeaderAlign,
	//			Transparent:     true,
	//		})
	//}
	SortDirection = SortDesc
	procsTbl.Select(0, 3)

	//.
	//	SetSelectedStyle(
	//		tcell.Style.Foreground(tcell.StyleDefault, tcell.ColorRed).
	//		Background(tcell.ColorYellow))

	for {
		// First get the processes BEFORE we sleep for the refresh. This
		//	helps us load the data within the window/time period of sleeping.
		processes, _ := devices.GetProcs(groupProcesses)

		// Construct the header row by column
		for col := 0; col < len(headerText); col++ {
			procsTbl.SetCell(
				0, col, &tview.TableCell{
					Text:            headerText[col].Text,
					Expansion:       headerText[col].Expand,
					MaxWidth:        headerText[col].MaxWidth,
					Color:           tblHeaderTextColor,
					BackgroundColor: tblHeaderBkgdColor,
					Align:           tblHeaderAlign,
					Transparent:     true,
				})
		}

		procs := sortProcs(processes, SortDirection, SortMethod)

		//switch SortDirection {
		//case SortAsc:
		//	SortDirection = SortDesc
		//case SortDesc:
		//	SortDirection = SortAsc
		//}

		//procsTbl.GetMouseCapture()

		// Get the initial sorting method (by default sorts by CPU)
		procsTbl.SetMouseCapture(
			func(action tview.MouseAction, event *tcell.EventMouse) (tview.MouseAction, *tcell.EventMouse) {
				// TODO: first clear all up or down arrows in the table header names
				headerText = getHeaderNames(groupProcesses)

				// get the selected column and set the sort direction
				_, selectedCol := procsTbl.GetSelection()
				if selectedCol == 0 {
					switch SortDirection {
					case SortAsc:
						headerText[selectedCol].Text = tblHeaderNames[0].Pid + UpArrow
						if groupProcesses {
							headerText[selectedCol].Text = tblHeaderNames[0].Cnt + UpArrow
						}
						procs = sortProcs(procs, SortAsc, SortPid)
					case SortDesc:
						headerText[selectedCol].Text = tblHeaderNames[0].Pid + DownArrow
						if groupProcesses {
							headerText[selectedCol].Text = tblHeaderNames[0].Cnt + DownArrow
						}
						procs = sortProcs(procs, SortDesc, SortPid)
					}
				}
				if selectedCol == 1 {
					switch SortDirection {
					case SortAsc:
						headerText[selectedCol].Text = tblHeaderNames[0].User + UpArrow
						procs = sortProcs(procs, SortAsc, SortUser)
					case SortDesc:
						headerText[selectedCol].Text = tblHeaderNames[0].User + DownArrow
						procs = sortProcs(procs, SortDesc, SortUser)
					}
				}
				if selectedCol == 2 {
					switch SortDirection {
					case SortAsc:
						headerText[selectedCol].Text = tblHeaderNames[0].Exec + UpArrow
						procs = sortProcs(procs, SortAsc, SortExec)
					case SortDesc:
						headerText[selectedCol].Text = tblHeaderNames[0].Exec + DownArrow
						procs = sortProcs(procs, SortDesc, SortExec)
					}
				}
				if selectedCol == 3 {
					switch SortDirection {
					case SortAsc:
						headerText[selectedCol].Text = tblHeaderNames[0].Cpu + UpArrow
						procs = sortProcs(procs, SortAsc, SortCpu)
					case SortDesc:
						headerText[selectedCol].Text = tblHeaderNames[0].Cpu + DownArrow
						procs = sortProcs(procs, SortDesc, SortCpu)
					}
				}
				if selectedCol == 4 {
					switch SortDirection {
					case SortAsc:
						headerText[selectedCol].Text = tblHeaderNames[0].Mem + UpArrow
						procs = sortProcs(procs, SortAsc, SortMem)
					case SortDesc:
						headerText[selectedCol].Text = tblHeaderNames[0].Mem + DownArrow
						procs = sortProcs(procs, SortDesc, SortMem)
					}
				}
				return action, event
			})

		time.Sleep(refresh)
		// We want to keep the amount of rows scrolled down then set it last
		//	after everything is drawn.
		rowOffset, _ := procsTbl.GetOffset()

		app.QueueUpdateDraw(func() {
			for i := range procs {
				// Start at _row == 1 and not zero because we don't want to
				//	overwrite the header row!
				_row := i + 1
				procsTbl.
					SetCell(_row, 0, &tview.TableCell{
						Text:        strconv.Itoa(procs[i].Pid),
						Align:       tview.AlignCenter,
						Expansion:   tblColMinWidth[0].Pid,
						MaxWidth:    tblColMaxWidth[0].Pid,
						Transparent: true,
					}).
					SetCell(_row, 1, &tview.TableCell{
						Text:        procs[i].User,
						Align:       tview.AlignLeft,
						Expansion:   tblColMinWidth[0].User,
						MaxWidth:    tblColMaxWidth[0].User,
						Transparent: true,
					}).
					SetCell(_row, 2, &tview.TableCell{
						Text:        procs[i].Name,
						Align:       tview.AlignLeft,
						Expansion:   tblColMinWidth[0].Exec,
						MaxWidth:    tblColMaxWidth[0].Exec,
						Transparent: true,
					}).
					SetCell(_row, 3, &tview.TableCell{
						Text:        formatFloat(procs[i].Cpu, 2),
						Align:       tview.AlignCenter,
						Expansion:   tblColMinWidth[0].Cpu,
						MaxWidth:    tblColMaxWidth[0].Cpu,
						Transparent: true,
					}).
					SetCell(_row, 4, &tview.TableCell{
						Text:        formatFloat(procs[i].Mem, 2),
						Align:       tview.AlignCenter,
						Expansion:   tblColMinWidth[0].Mem,
						MaxWidth:    tblColMaxWidth[0].Mem,
						Transparent: true,
					}) //.
				//	SetCell(i, 5, &tview.TableCell{
				//		Text:        "", // formatFloat(procs[i].Gpu),
				//		Align:       tview.AlignCenter,
				//		Expansion:   tblColMinWidth[0].Gpu,
				//		MaxWidth:    tblColMinWidth[0].Gpu,
				//		Transparent: true,
				//	})
				//}
			}
			// This helps retain the amount of rows scrolled by the user
			procsTbl.SetOffset(rowOffset, 0)
		})
	}
}
