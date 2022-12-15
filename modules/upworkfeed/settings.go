package upworkfeed

import (
	"github.com/olebedev/config"
	"github.com/wtfutil/wtf/cfg"
)

const (
	defaultFocusable = true
	defaultTitle     = "Upwork Feed"
)

// Settings defines the configuration properties for this module
type Settings struct {
	*cfg.Common

	// Define your settings attributes here
	feedUrl string `help:"URL of the Upwork RSS feed"`
}

// NewSettingsFromYAML creates a new settings instance from a YAML config block
func NewSettingsFromYAML(name string, ymlConfig *config.Config, globalConfig *config.Config) *Settings {
	common := cfg.NewCommonSettingsFromModule(name, defaultTitle, defaultFocusable, ymlConfig, globalConfig)
	settings := Settings{
		Common: common,

		// Configure your settings attributes here. See http://github.com/olebedev/config for type details
		feedUrl: ymlConfig.UString("feedUrl", ""),
	}

	return &settings
}
