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

const (
	SortPid  ProcSortMethod = "p"
	SortUser                = "u"
	SortExec                = "e"
	SortCpu                 = "c"
	SortMem                 = "m"
)

const (
	SortAsc  ProcSortDirection = "a"
	SortDesc                   = "d"
)

const (
	DownArrow string = "▼"
	UpArrow          = "▲"
)

type ProcTbl struct {
	Text     string
	MinWidth int
	MaxWidth int
}

type ProcTblCol struct {
	Pid  ProcTbl
	Cnt  ProcTbl
	User ProcTbl
	Exec ProcTbl
	Cpu  ProcTbl
	Mem  ProcTbl
	Gpu  ProcTbl
}

var procBoxLabel = "[ Processes ]"
var SortMethod ProcSortMethod
var SortDirection ProcSortDirection

var procTblCol = ProcTblCol{
	Pid:  ProcTbl{"PID", 3, 4},
	Cnt:  ProcTbl{"CNT", 3, 4},
	User: ProcTbl{"USER", 6, 8},
	Exec: ProcTbl{"EXEC", 12, 16},
	Cpu:  ProcTbl{"CPU%", 4, 4},
	Mem:  ProcTbl{"MEM%", 4, 4},
	Gpu:  ProcTbl{"GPU%", 4, 4},
}

func roundFloat(float float64, precision int) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(float*ratio) / ratio
}

func formatFloat(float float64, precision int) string {
	roundedFloat := roundFloat(float, precision)
	return strconv.FormatFloat(roundedFloat, 'g', precision, 64)
}

func getHeaderNames(groupProcesses bool) []ProcTbl {
	header := []ProcTbl{
		{Text: procTblCol.Pid.Text, MinWidth: procTblCol.Pid.MinWidth,
			MaxWidth: procTblCol.Pid.MaxWidth},
		{Text: procTblCol.User.Text, MinWidth: procTblCol.User.MinWidth,
			MaxWidth: procTblCol.User.MaxWidth},
		{Text: procTblCol.Exec.Text, MinWidth: procTblCol.Exec.MinWidth,
			MaxWidth: procTblCol.Exec.MaxWidth},
		{Text: procTblCol.Cpu.Text, MinWidth: procTblCol.Cpu.MinWidth,
			MaxWidth: procTblCol.Cpu.MaxWidth},
		{Text: procTblCol.Mem.Text, MinWidth: procTblCol.Mem.MinWidth,
			MaxWidth: procTblCol.Mem.MaxWidth},
		//{Text: procTblCol.Gpu.Text, MinWidth: procTblCol.Gpu.MinWidth,
		//	MaxWidth: procTblCol.Gpu.MaxWidth},
	}
	if groupProcesses {
		header[0].Text = procTblCol.Cnt.Text
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

	// Construct the header row by column
	headerText := getHeaderNames(groupProcesses)

	procsTbl.SetFixed(1, 0).
		SetSeparator(tview.BoxDrawingsLightVertical).
		SetSelectable(false, true).
		SetEvaluateAllRows(true).
		SetBorder(true).
		SetTitle(procBoxLabel)

	// Set the default sort direction and selected cell (CPU%, descending)
	SortDirection = SortDesc
	procsTbl.Select(0, 3)

	for {
		// First get the processes BEFORE we sleep for the refresh. This
		//	helps us load the data within the window/time period of sleeping.
		processes, _ := devices.GetProcs(groupProcesses)

		// Construct the header row by column
		for col := 0; col < len(headerText); col++ {
			procsTbl.SetCell(
				0, col, &tview.TableCell{
					Text:            headerText[col].Text,
					Expansion:       headerText[col].MinWidth,
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
						headerText[selectedCol].Text = procTblCol.Pid.Text + UpArrow
						if groupProcesses {
							headerText[selectedCol].Text = procTblCol.Cnt.Text + UpArrow
						}
						procs = sortProcs(procs, SortAsc, SortPid)
					case SortDesc:
						headerText[selectedCol].Text = procTblCol.Pid.Text + DownArrow
						if groupProcesses {
							headerText[selectedCol].Text = procTblCol.Cnt.Text + DownArrow
						}
						procs = sortProcs(procs, SortDesc, SortPid)
					}
				}
				if selectedCol == 1 {
					switch SortDirection {
					case SortAsc:
						headerText[selectedCol].Text = procTblCol.User.Text + UpArrow
						procs = sortProcs(procs, SortAsc, SortUser)
					case SortDesc:
						headerText[selectedCol].Text = procTblCol.User.Text + DownArrow
						procs = sortProcs(procs, SortDesc, SortUser)
					}
				}
				if selectedCol == 2 {
					switch SortDirection {
					case SortAsc:
						headerText[selectedCol].Text = procTblCol.Exec.Text + UpArrow
						procs = sortProcs(procs, SortAsc, SortExec)
					case SortDesc:
						headerText[selectedCol].Text = procTblCol.Exec.Text + DownArrow
						procs = sortProcs(procs, SortDesc, SortExec)
					}
				}
				if selectedCol == 3 {
					switch SortDirection {
					case SortAsc:
						headerText[selectedCol].Text = procTblCol.Cpu.Text + UpArrow
						procs = sortProcs(procs, SortAsc, SortCpu)
					case SortDesc:
						headerText[selectedCol].Text = procTblCol.Cpu.Text + DownArrow
						procs = sortProcs(procs, SortDesc, SortCpu)
					}
				}
				if selectedCol == 4 {
					switch SortDirection {
					case SortAsc:
						headerText[selectedCol].Text = procTblCol.Mem.Text + UpArrow
						procs = sortProcs(procs, SortAsc, SortMem)
					case SortDesc:
						headerText[selectedCol].Text = procTblCol.Mem.Text + DownArrow
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
						Expansion:   procTblCol.Pid.MinWidth,
						MaxWidth:    procTblCol.Pid.MaxWidth,
						Transparent: true,
					}).
					SetCell(_row, 1, &tview.TableCell{
						Text:        procs[i].User,
						Align:       tview.AlignLeft,
						Expansion:   procTblCol.User.MinWidth,
						MaxWidth:    procTblCol.User.MaxWidth,
						Transparent: true,
					}).
					SetCell(_row, 2, &tview.TableCell{
						Text:        procs[i].Name,
						Align:       tview.AlignLeft,
						Expansion:   procTblCol.Exec.MinWidth,
						MaxWidth:    procTblCol.Exec.MaxWidth,
						Transparent: true,
					}).
					SetCell(_row, 3, &tview.TableCell{
						Text:        formatFloat(procs[i].Cpu, 2),
						Align:       tview.AlignCenter,
						Expansion:   procTblCol.Cpu.MinWidth,
						MaxWidth:    procTblCol.Cpu.MaxWidth,
						Transparent: true,
					}).
					SetCell(_row, 4, &tview.TableCell{
						Text:        formatFloat(procs[i].Mem, 2),
						Align:       tview.AlignCenter,
						Expansion:   procTblCol.Mem.MinWidth,
						MaxWidth:    procTblCol.Mem.MaxWidth,
						Transparent: true,
					}) //.
				//	SetCell(i, 5, &tview.TableCell{
				//		Text:        formatFloat(procs[i].Gpu, 2),
				//		Align:       tview.AlignCenter,
				//		Expansion:   procTblCol.Gpu.MinWidth,
				//		MaxWidth:    procTblCol.Gpu.MaxWidth,
				//		Transparent: true,
				//	})
				//}
			}
			// This helps retain the amount of rows scrolled by the user
			procsTbl.SetOffset(rowOffset, 0)
		})
	}
}
