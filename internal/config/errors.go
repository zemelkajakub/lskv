package config

import (
	"errors"
)

var (
	ErrNoActiveProfile    = errors.New("no active profile: please set an active profile using 'lskv profile switch [alias]' or create a new profile using 'lskv profile init [alias]...'")
	ErrConfigFileNotFound = errors.New("no config file found: please recreate config file using 'lskv profile switch [alias]' or create a new profile using 'lskv profile init [alias]...'")
)
