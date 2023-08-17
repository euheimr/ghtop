package util

import (
	"github.com/knadh/koanf/parsers/toml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
	"log"
	"time"
)

type configVars struct {
	UpdateInterval time.Duration
	GroupProcesses bool
	TempScale      string
	Nvidia         bool
	//mbps           bool
}

var Config *configVars

func init() {
	cfg := readConfig()

	Config = &configVars{
		UpdateInterval: cfg.Duration("UpdateInterval"),
		GroupProcesses: cfg.Bool("GroupProcesses"),
		TempScale:      cfg.String("TempScale"),
		Nvidia:         cfg.Bool("Nvidia"),
		//mbps:           cfg.Bool("mbps"),
	}
}

func readConfig() (k *koanf.Koanf) {
	k = koanf.New(".")

	if err := k.Load(file.Provider("config.toml"), toml.Parser()); err != nil {
		log.Fatalf("Error loading config: %v", err)
	}
	return k
}
