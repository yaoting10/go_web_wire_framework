package config

import (
	"goboot/pkg/twitter"
)

type Twitter struct {
	twitter.Setting `mapstructure:",squash"`
	RedirectUrl     string `mapstructure:"redirectUrl" yaml:"redirectUrl"`
}
