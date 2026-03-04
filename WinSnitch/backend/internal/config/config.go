package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type Settings struct {
	ProgramDataDir   string `mapstructure:"program_data_dir"`
	ListenAddress    string `mapstructure:"listen_address"`
	BlocklistsURL    string `mapstructure:"blocklists_url"`
	NodeID           string `mapstructure:"node_id"`
	RequireElevation bool   `mapstructure:"require_elevation"`
}

func defaultProgramDataDir() string {
	if p := os.Getenv("ProgramData"); p != "" {
		return filepath.Join(p, "WinSnitch")
	}
	return `C:\ProgramData\WinSnitch`
}

func Load() (Settings, error) {
	v := viper.New()
	v.SetConfigName("winsnitch")
	v.SetConfigType("yaml")
	v.AddConfigPath(".")
	v.AddConfigPath(defaultProgramDataDir())

	v.SetDefault("program_data_dir", defaultProgramDataDir())
	v.SetDefault("listen_address", "127.0.0.1:47777")
	v.SetDefault("blocklists_url", "https://someonewhocares.org/hosts/hosts")
	v.SetDefault("node_id", "local-node")
	v.SetDefault("require_elevation", true)

	_ = v.ReadInConfig()

	var cfg Settings
	if err := v.Unmarshal(&cfg); err != nil {
		return Settings{}, fmt.Errorf("unmarshal config: %w", err)
	}
	return cfg, nil
}
