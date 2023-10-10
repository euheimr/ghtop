package devices

import (
	"github.com/shirou/gopsutil/v3/host"
	"log/slog"
	"os/user"
	"strconv"
	"strings"
)

type SysInfo struct {
	Hostname     string
	User         string
	SocketsCores string
	Threads      string
}

var SystemInfo *SysInfo

func init() {
	// Get Sysinfo data - this isn't in the for loop because this doesn't change
	//	during the lifetime of the program, thus we only get it once
	var hostInfo, _ = host.Info()

	u, _ := user.Current()
	usr := u.Username
	if len(strings.Split(u.Username, "\\")) == 1 {
		usr = strings.Split(u.Username, "\\")[0]
	}
	if len(strings.Split(u.Username, "\\")) > 1 {
		usr = strings.Split(u.Username, "\\")[1]
	}

	// We only need to initialize this once
	SystemInfo = &SysInfo{
		Hostname: strings.Split(hostInfo.Hostname, ".")[0],
		User:     usr,
		SocketsCores: strconv.FormatInt(int64(len(CpuInfo)), 10) + "/" +
			strconv.FormatInt(int64(CpuInfo[0].Cores)*int64(len(CpuInfo)), 10),
		Threads: strconv.FormatInt(int64(CpuInfo[0].Threads), 10),
	}
	SystemInfo.Hostname = strings.Trim(SystemInfo.Hostname, "\t ")
	SystemInfo.User = strings.Trim(SystemInfo.User, "\t ")
	SystemInfo.SocketsCores = strings.Trim(SystemInfo.SocketsCores, "\t ")
	SystemInfo.Threads = strings.Trim(SystemInfo.Threads, "\t ")
	slog.Debug("Init SysInfo (devices/sys.go)")
}
