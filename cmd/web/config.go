package main

import (
	"os"
	"time"
)

type AppConfig struct {
	Origin        string
	DebugMode     bool
	DevEnv        bool
	BaseImgUrl    string
	PasswordTTL   time.Duration
	ActivationTTL time.Duration

	PageSize         int
	MaxSearchResults int
}

func loadConfig(dev, debug bool) AppConfig {
	config := AppConfig{
		DevEnv:           dev,
		DebugMode:        debug || dev,
		PasswordTTL:      30 * time.Minute,
		ActivationTTL:    24 * time.Hour,
		PageSize:         24,
		MaxSearchResults: 12,
	}

	if dev {
		config.Origin = os.Getenv("DEV_ORIGIN")
		config.BaseImgUrl = os.Getenv("DEV_IMG_URL")
	} else {
		config.Origin = os.Getenv("ORIGIN")
		config.BaseImgUrl = os.Getenv("IMG_URL")
	}

	return config
}
