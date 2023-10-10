package internal

import (
	"github.com/knadh/koanf/parsers/toml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
	"log"
	"strings"
	"time"
)

type ConfigVars struct {
	Debug              bool
	IsProduction       bool
	UpdateInterval     time.Duration
	GroupProcesses     bool
	TempScale          string
	FormatMemAsPercent bool
	EnableNvidia       bool
	//mbps               bool
}

var Config ConfigVars

func init() {
	cfg, err := readConfig()
	if err != nil {
		log.Fatal(GetFuncName(), "Failed to read config")
	} else {
		updateConfigVars(cfg)
	}
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

	Config = ConfigVars{
		Debug:              cfg.Bool("Debug"),
		UpdateInterval:        cfg.Duration("UpdateInterval") * time.Millisecond,
		UpdateInterval:     cfg.Duration("UpdateInterval") * time.Millisecond,
		GroupProcesses:     cfg.Bool("GroupProcesses"),
		TempScale:          tempScale,
		FormatMemAsPercent: cfg.Bool("FormatMemoryAsPercent"),
		EnableNvidia:       cfg.Bool("EnableNvidia"),
	}
	return true
}

func readConfig() (cfg *koanf.Koanf, err error) {
	cfg = koanf.New(".")

	err = cfg.Load(file.Provider("config.toml"), toml.Parser())
	if err != nil {
		log.Fatalf(GetFuncName(), "Error loading config - ", err.Error())
		return cfg, err
	}
	return cfg, nil
}

//func WriteConfig() bool {}
