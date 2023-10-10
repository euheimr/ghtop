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
	ColumnTextColorSelected     tcell.Color = tcell.ColorYellow
	HeaderTextColor                         = tcell.ColorWhite
	HeaderTextBkgdColor                     = tcell.ColorRed
	HeaderTextColorSelected                 = tcell.ColorBlue
	HeaderTextBkgdColorSelected             = tcell.ColorWhite
	HeaderAlign                 int         = tview.AlignCenter
	HeaderPidLabel              string      = "PID"
	HeaderCntLabel                          = "CNT"
	HeaderUserLabel                         = "USER"
	HeaderExecLabel                         = "EXEC"
	HeaderCpuLabel                          = "CPU%"
	HeaderMemLabel                          = "MEM%"
	//HEADER_GPU_LABEL              = "GPU%"

	HeaderFormatBlink        = "[::l]"
	HeaderFormatReverseColor = "[::r]"
)

const DownArrow string = "▼" // descending
const UpArrow = "▲"          // ascending

var procsLabel = "Processes"
var baseHeader = Header{
	Pid: HeaderAttr{HeaderPidLabel, HeaderAlign,
		HeaderTextColor, HeaderTextBkgdColor,
		3, 4},
	User: HeaderAttr{HeaderUserLabel, HeaderAlign,
		HeaderTextColor, HeaderTextBkgdColor,
		5, 8},
	Exec: HeaderAttr{HeaderExecLabel, HeaderAlign,
		HeaderTextColor, HeaderTextBkgdColor,
		12, 16},
	Cpu: HeaderAttr{HeaderCpuLabel, HeaderAlign,
		HeaderTextColor, HeaderTextBkgdColor,
		4, 4},
	Mem: HeaderAttr{HeaderMemLabel, HeaderAlign,
		HeaderTextColor, HeaderTextBkgdColor,
		4, 4},
	//Gpu: HeaderAttr{HEADER_GPU_LABEL, HEADER_ALIGN,
	//  HEADER_TEXT_COLOR, HEADER_BKGD_COLOR,
	//	4, 4},
}

var (
	cfg                *internal.ConfigVars
	sortColumn         int
	sortColumnPrevious int
	sortDescending     bool
	header             Header
	//headerNames        [5]string
	processes []devices.Process
)
var (
	lastProcsCount int64
	lastTitle      string
	lastUpdate     time.Time
)

func init() {
	cfg = &internal.Config
	sortDescending = true
	sortColumn = devices.Cpu

	processes, _ = devices.GetProcs(cfg.GroupProcesses)
	header = updateHeaderText()
}

func updateHeaderText() (h Header) {
	h = baseHeader

	switch sortColumn {
	case devices.Pid:
		switch sortDescending {
		case true:
			if cfg.GroupProcesses {
				h.Pid.Text = HeaderCntLabel + DownArrow
				//h.Pid.TextColor = HeaderTextColorSelected
				//h.Pid.TextBkgdColor = HeaderTextBkgdColorSelected
			} else {
				h.Pid.Text = HeaderPidLabel + DownArrow
				//h.Pid.TextColor = HeaderTextColorSelected
				//h.Pid.TextBkgdColor = HeaderTextBkgdColorSelected
			}
		case false:
			if cfg.GroupProcesses {
				h.Pid.Text = HeaderCntLabel + UpArrow
				//h.Pid.TextColor = HeaderTextColorSelected
				//h.Pid.TextBkgdColor = HeaderTextBkgdColorSelected
			} else {
				h.Pid.Text = HeaderPidLabel + UpArrow
				//h.Pid.TextColor = HeaderTextColorSelected
				//h.Pid.TextBkgdColor = HeaderTextBkgdColorSelected
			}
		}
	case devices.User:
		switch sortDescending {
		case true:
			h.User.Text = HeaderUserLabel + DownArrow
			//h.User.TextColor = HeaderTextColorSelected
			//h.User.TextBkgdColor = HeaderTextBkgdColorSelected
		case false:
			h.User.Text = HeaderUserLabel + UpArrow
			//h.User.TextColor = HeaderTextColorSelected
			//h.User.TextBkgdColor = HeaderTextBkgdColorSelected
		}
	case devices.Exec:
		switch sortDescending {
		case true:
			h.Exec.Text = HeaderExecLabel + DownArrow
			//h.Exec.TextColor = HeaderTextColorSelected
			//h.Exec.TextBkgdColor = HeaderTextBkgdColorSelected
		case false:
			h.Exec.Text = HeaderExecLabel + UpArrow
			//h.Exec.TextColor = HeaderTextColorSelected
			//h.Exec.TextBkgdColor = HeaderTextBkgdColorSelected
		}
	case devices.Cpu:
		switch sortDescending {
		case true:
			h.Cpu.Text = HeaderCpuLabel + DownArrow
			h.Cpu.TextColor = tcell.ColorWhite
			//h.Cpu.TextBkgdColor = HeaderTextBkgdColorSelected
		case false:
			h.Cpu.Text = HeaderCpuLabel + UpArrow
			//h.Cpu.TextColor = HeaderTextColorSelected
			//h.Cpu.TextBkgdColor = HeaderTextBkgdColorSelected
		}
	case devices.Mem:
		switch sortDescending {
		case true:
			h.Mem.Text = HeaderMemLabel + DownArrow
			//h.Mem.TextColor = HeaderTextColorSelected
			//h.Mem.TextBkgdColor = HeaderTextBkgdColorSelected
		case false:
			h.Mem.Text = HeaderMemLabel + UpArrow
			//h.Mem.TextColor = HeaderTextColorSelected
			//h.Mem.TextBkgdColor = HeaderTextBkgdColorSelected
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
		SetBorder(true)
	//procsTbl.SetBackgroundColor(tcell.ColorGray)

	// Set the default sort direction and selected cell (CPU%, descending)
	procsTbl.Select(0, sortColumn)
	// todo: Input capture get and set sortColumn and sortDescending
	procsTbl.SetMouseCapture(
		func(action tview.MouseAction, event *tcell.EventMouse) (tview.MouseAction, *tcell.EventMouse) {
			if action == tview.MouseLeftClick || action == tview.MouseRightClick {
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
		sortColumnPrevious = sortColumn         // We remember previous column sorting
		_, sortColumn = procsTbl.GetSelection() // Overwrite and get current selection
		// If the newly selected column isn't the same, always start off by sortDescending
		if sortColumnPrevious != sortColumn {
			sortDescending = true
		}

		processes, _ = devices.GetProcs(cfg.GroupProcesses)
		procs, _ := devices.SortProcs(processes, sortColumn, sortDescending)
		procsCnt, _ := devices.GetProcsCount()
		procsCntLabel := strconv.FormatInt(procsCnt, 10)

		time.Sleep(update)
		// We want to keep the amount of rows scrolled down then set it last
		//	after everything is drawn.
		rowOffset, _ := procsTbl.GetOffset()

		app.QueueUpdateDraw(func() {
			/// This area is logic for drawing the Processes Title and border colors
			var title string
			if lastProcsCount < 1 {
				// Initialize lastProcsCount if 0 or nil
				lastProcsCount = procsCnt
			}
			// Less processes are usually good - DOWNARROW + GREEN
			if procsCnt < lastProcsCount {
				delta := strconv.FormatInt(lastProcsCount-procsCnt, 10)
				title = "[ " + procsLabel + " · " + GREEN + procsCntLabel + "(-" + delta + ")" + DownArrow + WHITE + " ]"
				procsTbl.SetBorderColor(tcell.ColorGreen).SetTitle(title)
				lastUpdate = time.Now().UTC()
			} else if procsCnt > lastProcsCount {
				// More processes are usually bad - UPARROW + RED
				delta := strconv.FormatInt(procsCnt-lastProcsCount, 10)
				title = "[ " + procsLabel + " · " + RED + procsCntLabel + "(+" + delta + ")" + UpArrow + WHITE + " ]"
				procsTbl.SetBorderColor(tcell.ColorRed).SetTitle(title)
				lastUpdate = time.Now().UTC()
			} else if procsCnt == lastProcsCount {
				sec := time.Now().UTC().Sub(lastUpdate)
				// If the current processes count is the same as the last update
				//	for more than X (if sec.Seconds() >= X) seconds, then
				//	overwrite to DASH + WHITE
				if sec.Seconds() >= 20 {
					title = "[ " + procsLabel + " · " + WHITE + procsCntLabel + "(-)" + WHITE + " ]"
					procsTbl.SetBorderColor(tcell.ColorWhite).SetTitle(title)
					//lastUpdate = time.Now().UTC()
				} // else if sec.Seconds() < 4 {
				//	procsTbl.SetTitle(lastTitle) //.SetBorderColor(tcell.ColorYellow)
				//	//lastUpdate = time.Now().UTC()
				//}
			}
			// Remember the number of processes and previous title for
			//	comparison on next iteration
			lastProcsCount = procsCnt
			lastTitle = title
			/// END PROCS BOX TITLE ////////////////////////////////////////////

			/// START TABLE UI /////////////////////////////////////////////////
			// Reset the table header
			header = updateHeaderText()

			// ** HEADER SETUP **
			procsTbl.SetCell(
				0, devices.Pid, &tview.TableCell{
					//Reference:       nil,
					Text:            header.Pid.Text,
					Align:           header.Pid.TextAlign,
					MaxWidth:        header.Pid.MaxWidth,
					Expansion:       header.Pid.MinWidth,
					Color:           header.Pid.TextColor,
					BackgroundColor: header.Pid.TextBkgdColor,
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
					Color:           header.User.TextColor,
					BackgroundColor: header.User.TextBkgdColor,
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
					Color:           header.Exec.TextColor,
					BackgroundColor: header.Exec.TextBkgdColor,
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
					Color:           header.Cpu.TextColor,
					BackgroundColor: header.Cpu.TextBkgdColor,
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
					Color:           header.Mem.TextColor,
					BackgroundColor: header.Mem.TextBkgdColor,
					Transparent:     false,
					//Attributes:      0,
					//NotSelectable:   false,
					//Clicked:         nil,
				})

			// todo: get procs and render them here
			for i := range procs {
				// Start at _row == 1 and not zero because we don't want to
				//	overwrite the header row!
				row := i + 1
				procsTbl.
					SetCell(row, devices.Pid, &tview.TableCell{
						//Color:       header.Pid.TextColor,
						Text:        strconv.Itoa(procs[i].Pid),
						Align:       tview.AlignCenter,
						Expansion:   header.Pid.MinWidth,
						MaxWidth:    header.Pid.MaxWidth,
						Transparent: true,
					}).
					SetCell(row, devices.User, &tview.TableCell{
						//Color:       header.User.TextColor,
						Text:        procs[i].User,
						Align:       tview.AlignLeft,
						Expansion:   header.User.MinWidth,
						MaxWidth:    header.User.MaxWidth,
						Transparent: true,
					}).
					SetCell(row, devices.Exec, &tview.TableCell{
						//Color:       header.Exec.TextColor,
						Text:        procs[i].Name,
						Align:       tview.AlignLeft,
						Expansion:   header.Exec.MinWidth,
						MaxWidth:    header.Exec.MaxWidth,
						Transparent: true,
					}).
					SetCell(row, devices.Cpu, &tview.TableCell{
						//Color:       header.Cpu.TextColor,
						Text:        formatValue(procs[i].Cpu, 2),
						Align:       tview.AlignCenter,
						Expansion:   header.Cpu.MinWidth,
						MaxWidth:    header.Cpu.MaxWidth,
						Transparent: true,
					}).
					SetCell(row, devices.Mem, &tview.TableCell{
						//Color:       header.Mem.TextColor,
						Text:        formatValue(procs[i].Mem, 2),
						Align:       tview.AlignCenter,
						Expansion:   header.Mem.MinWidth,
						MaxWidth:    header.Mem.MaxWidth,
						Transparent: true,
					}) //.
				//	SetCell(row, devices.Gpu, &tview.TableCell{
				//	    //Color:       header.Gpu.TextColor,
				//		Text:        formatValue(procs[i].Gpu, 2),
				//		Align:       tview.AlignCenter,
				//		Expansion:   header.Gpu.MinWidth,
				//		MaxWidth:    header.Gpu.MaxWidth,
				//		Transparent: true,
				//	})
				//}
			}
			header = updateHeaderText()

			// This helps restore the amount of rows scrolled by the user
			procsTbl.SetOffset(rowOffset, 0)
		})
	}

}
