package internal

import (
	"github.com/knadh/koanf/parsers/toml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
	"log/slog"
	"strings"
	"time"
)

type ConfigVars struct {
	Debug                 bool
	TempScale             string
	UpdateInterval        time.Duration
	EnableNvidia          bool
	EnableUIButtons       bool
	GroupProcesses        bool
	SelectedViewOverride  int
	ShowOnlyUserProcesses bool
	//FormatMemAsPercent bool
	//mbps               bool
}

var Config ConfigVars

func init() {
	cfg, err := readConfig()
	//Log.Debug("Read config ...")
	if err != nil {
		slog.Error(GetFuncName(), err)
	}
	updateConfigVars(cfg)
	slog.Debug("Init config.go")
}

func updateConfigVars(cfg *koanf.Koanf) bool {
	var tempScale string
	var scale = strings.ToUpper(cfg.String("TempScale"))
	switch scale {
	default:
		tempScale = "C"
	case "C":
		tempScale = "C"
	case "F":
		tempScale = "F"
	}

	viewOverride := cfg.Int("SelectedViewOverride")
	// This sets the default Selected View to 0 if the config variable is out of range
	if viewOverride < 0 || viewOverride > 2 {
		viewOverride = 0
	}

	Config = ConfigVars{
		Debug:                 cfg.Bool("Debug"),
		UpdateInterval:        cfg.Duration("UpdateInterval") * time.Millisecond,
		EnableNvidia:          cfg.Bool("EnableNvidia"),
		EnableUIButtons:       cfg.Bool("EnableUIButtons"),
		GroupProcesses:        cfg.Bool("GroupProcesses"),
		SelectedViewOverride:  viewOverride,
		TempScale:             tempScale,
		ShowOnlyUserProcesses: cfg.Bool("ShowOnlyUserProcesses"),
		//FormatMemAsPercent: cfg.Bool("FormatMemoryAsPercent"),
	}
	return true
}

func readConfig() (cfg *koanf.Koanf, err error) {
	cfg = koanf.New(".")

	err = cfg.Load(file.Provider("config.toml"), toml.Parser())
	if err != nil {
		slog.Error(GetFuncName(), "Error loading config - ", err.Error())
		return cfg, err
	}
	return cfg, nil
}

//func WriteConfig() bool {}
