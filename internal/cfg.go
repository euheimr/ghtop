package internal

import (
	"github.com/knadh/koanf/parsers/toml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
	"log/slog"
	"time"
)

type ConfigVars struct {
	Debug          bool
	UpdateInterval time.Duration
	Celsius        bool
	EnableNvidia   bool
}

var Cfg *ConfigVars

const CONFIG_FILENAME = "cfg.toml"

func init() {
	cfg := koanf.New(".")

	if err := cfg.Load(file.Provider(CONFIG_FILENAME), toml.Parser()); err != nil {
		slog.Info("Could not load config file at `" + CONFIG_FILENAME + "` ! Using defaults ...")
		Cfg = &ConfigVars{
			Debug:          false,
			UpdateInterval: 100 * time.Millisecond,
			Celsius:        true,
			EnableNvidia:   false,
		}
	} else {
		updateConfigVars(cfg)
		slog.Info("ConfigVars loaded from `" + CONFIG_FILENAME + "`")
	}
}

func updateConfigVars(cfg *koanf.Koanf) {
	Cfg = &ConfigVars{
		Debug:          cfg.Bool("Debug"),
		UpdateInterval: cfg.Duration("UpdateInterval") * time.Millisecond,
		Celsius:        cfg.Bool("Celsius"),
		EnableNvidia:   cfg.Bool("EnableNvidia"),
	}
}
