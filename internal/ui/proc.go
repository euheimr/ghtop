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

type TableHeaderAttr struct {
	Text          string
	TextAlign     int
	TextColor     tcell.Color
	TextBkgdColor tcell.Color
	MinWidth      int
	MaxWidth      int
}

type TableHeader struct {
	Pid  TableHeaderAttr
	User TableHeaderAttr
	Exec TableHeaderAttr
	Cpu  TableHeaderAttr
	Mem  TableHeaderAttr
	//Gpu TableHeaderAttr
}

const (
	// Table header labels
	DOWN_ARROW         string = "▼"
	UP_ARROW                  = "▲"
	HEADER_TABLE_LABEL        = "[ Processes · "
	HEADER_PID_LABEL          = "PID"
	HEADER_CNT_LABEL          = "CNT"
	HEADER_USER_LABEL         = "USR"
	HEADER_EXEC_LABEL         = "EXEC"
	HEADER_CPU_LABEL          = "CPU%"
	HEADER_MEM_LABEL          = "MEM%"
	//HEADER_GPU_LABEL        = "GPU%"

	// Table header text formatting
	//HEADER_FMT_BLINK         = "[::l]"
	//HEADER_FMT_REVERSE_COLOR = "[::r]"
	HEADER_TEXT_ALIGN      int         = tview.AlignCenter
	HEADER_TEXT_COLOR      tcell.Color = tcell.ColorWhite
	HEADER_TEXT_BKGD_COLOR             = tcell.ColorBlue
)

var (
	lastProcsCount uint64
	lastUpdate     time.Time
	tableTitle     string
)

var (
	header            TableHeader
	defaultHeader     TableHeader
	processes         []devices.Process
	currentColumnSort int
	lastColumnSort    int
	sortDescending    bool
)

func init() {
	// This is the DEFAULT header. The header text is updated/overwritten with updateHeaderText()
	defaultHeader = TableHeader{
		Pid: TableHeaderAttr{
			Text:          HEADER_PID_LABEL,
			TextAlign:     HEADER_TEXT_ALIGN,
			TextColor:     HEADER_TEXT_COLOR,
			TextBkgdColor: HEADER_TEXT_BKGD_COLOR,
			MinWidth:      0,
			MaxWidth:      5,
		},
		User: TableHeaderAttr{
			Text:          HEADER_USER_LABEL,
			TextAlign:     HEADER_TEXT_ALIGN,
			TextColor:     HEADER_TEXT_COLOR,
			TextBkgdColor: HEADER_TEXT_BKGD_COLOR,
			MinWidth:      8,
			MaxWidth:      6,
		},
		Exec: TableHeaderAttr{
			Text:          HEADER_EXEC_LABEL,
			TextAlign:     HEADER_TEXT_ALIGN,
			TextColor:     HEADER_TEXT_COLOR,
			TextBkgdColor: HEADER_TEXT_BKGD_COLOR,
			MinWidth:      22,
			MaxWidth:      9,
		},
		Cpu: TableHeaderAttr{
			Text:          HEADER_CPU_LABEL,
			TextAlign:     HEADER_TEXT_ALIGN,
			TextColor:     HEADER_TEXT_COLOR,
			TextBkgdColor: HEADER_TEXT_BKGD_COLOR,
			MinWidth:      2,
			MaxWidth:      3,
		},
		Mem: TableHeaderAttr{
			Text:          HEADER_MEM_LABEL,
			TextAlign:     HEADER_TEXT_ALIGN,
			TextColor:     HEADER_TEXT_COLOR,
			TextBkgdColor: HEADER_TEXT_BKGD_COLOR,
			MinWidth:      2,
			MaxWidth:      3,
		},
		//GPU: TableHeaderAttr{
		//	Text:          HEADER_GPU_LABEL,
		//	TextAlign:     HEADER_TEXT_ALIGN,
		//	TextColor:     HEADER_TEXT_COLOR,
		//	TextBkgdColor: HEADER_TEXT_BKGD_COLOR,
		//	MinWidth:      4,
		//	MaxWidth:      4,
		//},
	}
	header = defaultHeader

	// seed the default values
	sortDescending = true
	currentColumnSort = devices.Cpu
}

func updateHeaderText(groupProcesses bool) {
	// reset the header to the default header, effectively overwriting header text
	header = defaultHeader

	// This determines how to format the table header depending on which column is being sorted
	//	and if it's descending (high to low) or not (low to high - ascending)
	switch sortDescending {
	case true:
		switch currentColumnSort {
		case devices.Pid:
			if groupProcesses {
				header.Pid.Text = HEADER_CNT_LABEL + DOWN_ARROW
			} else {
				header.Pid.Text = HEADER_PID_LABEL + DOWN_ARROW
			}
		case devices.User:
			header.User.Text = HEADER_USER_LABEL + DOWN_ARROW
		case devices.Exec:
			header.Exec.Text = HEADER_EXEC_LABEL + DOWN_ARROW
		case devices.Cpu:
			header.Cpu.Text = HEADER_CPU_LABEL + DOWN_ARROW
		case devices.Mem:
			header.Mem.Text = HEADER_MEM_LABEL + DOWN_ARROW
		}
	case false:
		switch currentColumnSort {
		case devices.Pid:
			if groupProcesses {
				header.Pid.Text = HEADER_CNT_LABEL + UP_ARROW
			} else {
				header.Pid.Text = HEADER_PID_LABEL + UP_ARROW
			}
		case devices.User:
			header.User.Text = HEADER_USER_LABEL + UP_ARROW
		case devices.Exec:
			header.Exec.Text = HEADER_EXEC_LABEL + UP_ARROW
		case devices.Cpu:
			header.Cpu.Text = HEADER_CPU_LABEL + UP_ARROW
		case devices.Mem:
			header.Mem.Text = HEADER_MEM_LABEL + UP_ARROW
		}
	}
}

func formatTableHeaderBorder(procsTbl *tview.Table, currentProcsCount uint64) {
	procsCountLabel := strconv.FormatUint(currentProcsCount, 10)

	// this logic formats the table border and title with colors and processes count
	//	depending on the number of processes increasing (red) or decreasing (green)
	if currentProcsCount < lastProcsCount {
		// Less processes is formatted with an DOWN arrow and GREEN in color
		processesDelta := strconv.FormatUint(lastProcsCount-currentProcsCount, 10)
		tableTitle = HEADER_TABLE_LABEL + internal.GREEN + procsCountLabel + "(-" + processesDelta + ")" +
			DOWN_ARROW + internal.WHITE + " ]"
		procsTbl.SetBorderColor(tcell.ColorGreen).SetTitle(tableTitle)
		lastUpdate = time.Now().UTC()
	} else if currentProcsCount > lastProcsCount {
		// More processes is formatted with a UP arrow and RED in color
		processesDelta := strconv.FormatUint(currentProcsCount-lastProcsCount, 10)
		tableTitle = HEADER_TABLE_LABEL + internal.RED + procsCountLabel + "(+" + processesDelta + ")" +
			UP_ARROW + internal.WHITE + " ]"
		procsTbl.SetBorderColor(tcell.ColorRed).SetTitle(tableTitle)
		lastUpdate = time.Now().UTC()
	} else if currentProcsCount == lastProcsCount {
		lastUpdateDelta := time.Now().UTC().Sub(lastUpdate)
		// if the count of current processes are the same for longer than X seconds,
		//	then reset the formatting of the processes box to DASH + WHITE
		if lastUpdateDelta.Seconds() >= 20 {
			tableTitle = HEADER_TABLE_LABEL + internal.WHITE + procsCountLabel + " ]"
			procsTbl.SetBorderColor(tcell.ColorWhite).SetTitle(tableTitle)
		}
	}
}

func formatValue(value interface{}, precision int) (val string) {
	switch val := value.(type) {
	case float64:
		ratio := math.Pow(10, float64(precision))
		roundedFloat := math.Round(val*ratio) / ratio
		return strconv.FormatFloat(roundedFloat, 'g', 3, 64)
	case uint64:
		return strconv.FormatUint(val, 10)
	}
	return val
}

func UpdateProcs(app *tview.Application, procsTbl *tview.Table, groupProcesses bool, update time.Duration) {

	procsTbl.SetFixed(1, 0).
		SetSeparator(tview.BoxDrawingsLightVertical).
		SetSelectable(false, true).
		SetEvaluateAllRows(true).
		SetBorder(true)

	procsTbl.Select(0, currentColumnSort)

	procsTbl.SetMouseCapture(
		func(action tview.MouseAction, event *tcell.EventMouse) (tview.MouseAction, *tcell.EventMouse) {
			if action == tview.MouseLeftClick {
				switch sortDescending {
				case true:
					sortDescending = false
					break
				default:
					sortDescending = true
				}
			}
			return action, event
		})

	for {
		// set a timestamp right now to determine how long we should sleep before the next UI update
		timestamp := time.Now().UTC()

		lastColumnSort = currentColumnSort             // remember which column we last sorted by
		_, currentColumnSort = procsTbl.GetSelection() // overwrite and get current selected column

		// This is a bugfix for sorting resetting when clicking a new column to sort by
		if currentColumnSort != lastColumnSort {
			// If the newly selected column isnt the same, always default to sorting processes
			//	by Descending order (high to low).
			sortDescending = true
		}

		processes, _ = devices.GetProcesses(currentColumnSort, sortDescending, groupProcesses, update)

		// initialize the lastProcsCount value if the program just started and is 0
		currentProcsCount, _ := devices.GetProcessesCount()
		if lastProcsCount < 1 {
			lastProcsCount = currentProcsCount
		}

		delta := time.Now().UTC().UnixMilli() - timestamp.UnixMilli()
		if update.Milliseconds() < delta {
			// If it takes longer than our set update interval to get the processes info and etc,
			//	then just don't sleep by breaking and updating UI ASAP.
			break
		} else {
			// If the time delta is lower than the set update interval, just sleep for the set
			//	update interval. We also account for the time it takes to get the current processes
			sleepDelta := update - time.Duration(delta)
			//procsTbl.SetTitle(strconv.FormatInt(int64(update/time.Millisecond), 10) + "ms")
			time.Sleep(sleepDelta)
		}

		// We want to save the amount of rows scrolled down then set it last after everything
		//	else is ready. We do it AFTER the sleep so the scrolling doesn't "rubber-band" backwards
		rowOffset, _ := procsTbl.GetOffset()

		app.QueueUpdateDraw(func() {
			formatTableHeaderBorder(procsTbl, currentProcsCount)

			// Remember the number of processes for the next iteration/update
			lastProcsCount = currentProcsCount

			updateHeaderText(groupProcesses)

			tcHeaderPid := &tview.TableCell{
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
				//Clicked: nil,
			}
			tcHeaderPid.SetClickedFunc(func() bool {
				//tcHeaderPid.Text = "??"
				return false
			})

			tcHeaderUser := &tview.TableCell{
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
				//Clicked: nil,
			}
			tcHeaderUser.SetClickedFunc(func() bool {
				return false
			})

			tcHeaderExec := &tview.TableCell{
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
				//Clicked: nil,
			}
			tcHeaderExec.SetClickedFunc(func() bool {
				return false
			})

			tcHeaderCpu := &tview.TableCell{
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
				//Clicked: nil,
			}
			tcHeaderCpu.SetClickedFunc(func() bool {
				return false
			})

			tcHeaderMem := &tview.TableCell{
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
				//Clicked: nil,
			}
			tcHeaderMem.SetClickedFunc(func() bool {
				return false
			})

			// Setting the actual header cells row with SetCell for EACH COLUMN
			procsTbl.SetCell(0, devices.Pid, tcHeaderPid).
				SetCell(0, devices.User, tcHeaderUser).
				SetCell(0, devices.Exec, tcHeaderExec).
				SetCell(0, devices.Cpu, tcHeaderCpu).
				SetCell(0, devices.Mem, tcHeaderMem)

			// Now we draw the rows of processes here
			for i := range processes {
				row := i + 1 // skip the header row. We don't want to overwrite it!

				tcDataPid := &tview.TableCell{
					//Color:       header.Pid.TextColor,
					Text:        strconv.Itoa(processes[i].Pid),
					Align:       tview.AlignCenter,
					Expansion:   header.Pid.MinWidth,
					MaxWidth:    header.Pid.MaxWidth,
					Transparent: true,
				}
				tcDataUser := &tview.TableCell{
					//Color:       header.User.TextColor,
					Text:        processes[i].User,
					Align:       tview.AlignLeft,
					Expansion:   header.User.MinWidth,
					MaxWidth:    header.User.MaxWidth,
					Transparent: true,
				}
				tcDataExec := &tview.TableCell{
					//Color:       header.Exec.TextColor,
					Text:        processes[i].Name,
					Align:       tview.AlignLeft,
					Expansion:   header.Exec.MinWidth,
					MaxWidth:    header.Exec.MaxWidth,
					Transparent: true,
				}
				tcDataCpu := &tview.TableCell{
					//Color:       header.Cpu.TextColor,
					Text:        formatValue(processes[i].Cpu, 2),
					Align:       tview.AlignCenter,
					Expansion:   header.Cpu.MinWidth,
					MaxWidth:    header.Cpu.MaxWidth,
					Transparent: true,
				}
				tcDataMem := &tview.TableCell{
					//Color:       header.Mem.TextColor,
					Text:        formatValue(processes[i].Mem, 2),
					Align:       tview.AlignCenter,
					Expansion:   header.Mem.MinWidth,
					MaxWidth:    header.Mem.MaxWidth,
					Transparent: true,
				}
				//tcDataGpu := &tview.TableCell{
				//	//Color:       header.Gpu.TextColor,
				//	Text:        formatValue(processes[i].Gpu, 2),
				//	Align:       tview.AlignCenter,
				//	Expansion:   header.Gpu.MinWidth,
				//	MaxWidth:    header.Gpu.MaxWidth,
				//	Transparent: true,
				//}

				procsTbl.SetCell(row, devices.Pid, tcDataPid).
					SetCell(row, devices.User, tcDataUser).
					SetCell(row, devices.Exec, tcDataExec).
					SetCell(row, devices.Cpu, tcDataCpu).
					SetCell(row, devices.Mem, tcDataMem) //.
				// SetCell(row, devices.Gpu, tcDataGpu)
			}

			procsTbl.SetOffset(rowOffset, 0)
		})
	}
}
