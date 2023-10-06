package ui

import (
	"github.com/euheimr/ghtop/internal"
	"github.com/euheimr/ghtop/internal/devices"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"math"
	"strconv"
	"time"
)

type HeaderAttr struct {
	Text          string
	TextAlign     int
	TextColor     tcell.Color
	TextBkgdColor tcell.Color
	MinWidth      int
	MaxWidth      int
}
type Header struct {
	Pid  HeaderAttr
	Cnt  HeaderAttr
	User HeaderAttr
	Exec HeaderAttr
	Cpu  HeaderAttr
	Mem  HeaderAttr
	//Gpu  *HeaderAttr
}

const (
	PROCS_LABEL       string      = "[ Processes ]"
	HEADER_TEXT_COLOR tcell.Color = tcell.ColorYellow
	HEADER_BKGD_COLOR             = tcell.ColorRed
	HEADER_ALIGN      int         = tview.AlignCenter
	HEADER_PID_LABEL  string      = "PID"
	HEADER_CNT_LABEL              = "CNT"
	HEADER_USER_LABEL             = "USER"
	HEADER_EXEC_LABEL             = "EXEC"
	HEADER_CPU_LABEL              = "CPU%"
	HEADER_MEM_LABEL              = "MEM%"
	//HEADER_GPU_LABEL              = "GPU%"
)

var (
	cfg                *internal.ConfigVars
	sortColumn         int
	sortColumnPrevious int
	sortDescending     bool
	baseHeader         Header
	header             Header
	//headerNames        [5]string
	processes []devices.Process
)

func init() {
	cfg = &internal.Config
	sortColumn = devices.Cpu
	sortDescending = true

	baseHeader = Header{
		Pid: HeaderAttr{HEADER_PID_LABEL, HEADER_ALIGN,
			HEADER_TEXT_COLOR, HEADER_BKGD_COLOR,
			3, 4},
		User: HeaderAttr{HEADER_USER_LABEL, HEADER_ALIGN,
			HEADER_TEXT_COLOR, HEADER_BKGD_COLOR,
			5, 8},
		Exec: HeaderAttr{HEADER_EXEC_LABEL, HEADER_ALIGN,
			HEADER_TEXT_COLOR, HEADER_BKGD_COLOR,
			12, 16},
		Cpu: HeaderAttr{HEADER_CPU_LABEL, HEADER_ALIGN,
			HEADER_TEXT_COLOR, HEADER_BKGD_COLOR,
			4, 4},
		Mem: HeaderAttr{HEADER_MEM_LABEL, HEADER_ALIGN,
			HEADER_TEXT_COLOR, HEADER_BKGD_COLOR,
			4, 4},
		//Gpu: HeaderAttr{HEADER_GPU_LABEL, HEADER_ALIGN,
		//  HEADER_TEXT_COLOR, HEADER_BKGD_COLOR,
		//	4, 4},
	}

	processes, _ = devices.GetProcs(cfg.GroupProcesses)
	header = baseHeader
}

func updateHeaderNames() (h Header) {
	const DownArrow string = "▼" // descending
	const UpArrow = "▲"          // ascending
	// Reset the table header
	h = baseHeader

	if cfg.GroupProcesses {
		h.Pid.Text = HEADER_CNT_LABEL
	}

	switch sortColumn {
	case devices.Pid:
		switch sortDescending {
		case true:
			if cfg.GroupProcesses {
				h.Pid.Text = HEADER_CNT_LABEL + DownArrow
			} else {
				h.Pid.Text = baseHeader.Pid.Text + DownArrow
			}
		case false:
			if cfg.GroupProcesses {
				h.Pid.Text = HEADER_CNT_LABEL + UpArrow
			} else {
				h.Pid.Text = baseHeader.Pid.Text + UpArrow
			}

		}
	case devices.User:
		switch sortDescending {
		case true:
			h.User.Text = baseHeader.User.Text + DownArrow
		case false:
			h.User.Text = baseHeader.User.Text + UpArrow
		}
	case devices.Exec:
		switch sortDescending {
		case true:
			h.Exec.Text = baseHeader.Exec.Text + DownArrow
		case false:
			h.Exec.Text = baseHeader.Exec.Text + UpArrow
		}
	case devices.Cpu:
		switch sortDescending {
		case true:
			h.Cpu.Text = baseHeader.Cpu.Text + DownArrow
		case false:
			h.Cpu.Text = baseHeader.Cpu.Text + UpArrow
		}
	case devices.Mem:
		switch sortDescending {
		case true:
			h.Mem.Text = baseHeader.Mem.Text + DownArrow
		case false:
			h.Mem.Text = baseHeader.Mem.Text + UpArrow
		}
	}
	return h
}

func formatValue(float interface{}, precision int) (f string) {
	switch f := float.(type) {
	case float64:
		ratio := math.Pow(10, float64(precision))
		// round the float first
		roundedFloat := math.Round(f*ratio) / ratio
		// return the float converted to string
		return strconv.FormatFloat(roundedFloat, 'g', 3, 64)
	case uint64:
		return strconv.FormatUint(f, 10)
	}
	return f
}

func UpdateProcs(app *tview.Application, procsTbl *tview.Table, update time.Duration) {

	procsTbl.SetFixed(1, 0).
		SetSeparator(tview.BoxDrawingsLightVertical).
		SetSelectable(false, true).
		SetEvaluateAllRows(true).
		SetBorder(true).
		SetTitle(PROCS_LABEL)
	// Set the default sort direction and selected cell (CPU%, descending)
	procsTbl.Select(0, sortColumn)
	// todo: Input capture get and set sortColumn and sortDescending
	procsTbl.SetMouseCapture(
		func(action tview.MouseAction, event *tcell.EventMouse) (tview.MouseAction, *tcell.EventMouse) {
			if action == tview.MouseLeftClick {
				switch sortDescending {
				case true:
					sortDescending = false
				case false:
					sortDescending = true
				}
			}
			return action, event
		})

	for {
		processes, _ = devices.GetProcs(cfg.GroupProcesses)
		sortColumnPrevious = sortColumn         // We remember previous column sorting
		_, sortColumn = procsTbl.GetSelection() // Overwrite and get current selection
		// If the newly selected column isn't the same, always start off by sortDescending
		if sortColumnPrevious != sortColumn {
			sortDescending = true
		}
		procs, _ := devices.SortProcs(processes, sortColumn, sortDescending)
		header = updateHeaderNames()

		time.Sleep(update)
		// We want to keep the amount of rows scrolled down then set it last
		//	after everything is drawn.
		rowOffset, _ := procsTbl.GetOffset()

		app.QueueUpdateDraw(func() {

			// ** HEADER SETUP **
			procsTbl.SetCell(
				0, devices.Pid, &tview.TableCell{
					//Reference:       nil,
					Text:            header.Pid.Text,
					Align:           header.Pid.TextAlign,
					MaxWidth:        header.Pid.MaxWidth,
					Expansion:       header.Pid.MinWidth,
					Color:           HEADER_TEXT_COLOR,
					BackgroundColor: HEADER_BKGD_COLOR,
					Transparent:     false,
					//Attributes:      0,
					//NotSelectable:   false,
					//Clicked:         SetClickedFunc(),
				})
			procsTbl.SetCell(
				0, devices.User, &tview.TableCell{
					//Reference:       nil,
					Text:            header.User.Text,
					Align:           header.User.TextAlign,
					MaxWidth:        header.User.MaxWidth,
					Expansion:       header.User.MinWidth,
					Color:           HEADER_TEXT_COLOR,
					BackgroundColor: HEADER_BKGD_COLOR,
					Transparent:     false,
					//Attributes:      0,
					//NotSelectable:   false,
					//Clicked:         nil,
				})
			procsTbl.SetCell(
				0, devices.Exec, &tview.TableCell{
					//Reference:       nil,
					Text:            header.Exec.Text,
					Align:           header.Exec.TextAlign,
					MaxWidth:        header.Exec.MaxWidth,
					Expansion:       header.Exec.MinWidth,
					Color:           HEADER_TEXT_COLOR,
					BackgroundColor: HEADER_BKGD_COLOR,
					Transparent:     false,
					//Attributes:      0,
					//NotSelectable:   false,
					//Clicked:         nil,
				})
			procsTbl.SetCell(
				0, devices.Cpu, &tview.TableCell{
					//Reference:       nil,
					Text:            header.Cpu.Text,
					Align:           header.Cpu.TextAlign,
					MaxWidth:        header.Cpu.MaxWidth,
					Expansion:       header.Cpu.MinWidth,
					Color:           HEADER_TEXT_COLOR,
					BackgroundColor: HEADER_BKGD_COLOR,
					Transparent:     false,
					//Attributes:      0,
					//NotSelectable:   false,
					//Clicked:         nil,
				})
			procsTbl.SetCell(
				0, devices.Mem, &tview.TableCell{
					//Reference:       nil,
					Text:            header.Mem.Text,
					Align:           header.Mem.TextAlign,
					MaxWidth:        header.Mem.MaxWidth,
					Expansion:       header.Mem.MinWidth,
					Color:           HEADER_TEXT_COLOR,
					BackgroundColor: HEADER_BKGD_COLOR,
					Transparent:     false,
					//Attributes:      0,
					//NotSelectable:   false,
					//Clicked:         nil,
				})

			// todo: get procs and render them here
			for i := range procs {
				// Start at _row == 1 and not zero because we don't want to
				//	overwrite the header row!
				_row := i + 1
				procsTbl.
					SetCell(_row, devices.Pid, &tview.TableCell{
						Text:        strconv.Itoa(procs[i].Pid),
						Align:       tview.AlignCenter,
						Expansion:   header.Pid.MinWidth,
						MaxWidth:    header.Pid.MaxWidth,
						Transparent: true,
					}).
					SetCell(_row, devices.User, &tview.TableCell{
						Text:        procs[i].User,
						Align:       tview.AlignLeft,
						Expansion:   header.User.MinWidth,
						MaxWidth:    header.User.MaxWidth,
						Transparent: true,
					}).
					SetCell(_row, devices.Exec, &tview.TableCell{
						Text:        procs[i].Name,
						Align:       tview.AlignLeft,
						Expansion:   header.Exec.MinWidth,
						MaxWidth:    header.Exec.MaxWidth,
						Transparent: true,
					}).
					SetCell(_row, devices.Cpu, &tview.TableCell{
						Text:        formatValue(procs[i].Cpu, 2),
						Align:       tview.AlignCenter,
						Expansion:   header.Cpu.MinWidth,
						MaxWidth:    header.Cpu.MaxWidth,
						Transparent: true,
					}).
					SetCell(_row, devices.Mem, &tview.TableCell{
						Text:        formatValue(procs[i].Mem, 2),
						Align:       tview.AlignCenter,
						Expansion:   header.Mem.MinWidth,
						MaxWidth:    header.Mem.MaxWidth,
						Transparent: true,
					}) //.
				//	SetCell(_row, devices.Gpu, &tview.TableCell{
				//		Text:        formatValue(procs[i].Gpu, 2),
				//		Align:       tview.AlignCenter,
				//		Expansion:   header.Gpu.MinWidth,
				//		MaxWidth:    header.Gpu.MaxWidth,
				//		Transparent: true,
				//	})
				//}
			}

			// This helps restore the amount of rows scrolled by the user
			procsTbl.SetOffset(rowOffset, 0)
		})
	}

}
