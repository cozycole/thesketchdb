package main

import (
	"time"

	"github.com/cozycole/thesketchdb/internal/web"
)

func loadConfig(dev, debug bool) web.AppConfig {
	config := web.AppConfig{
		DevEnv:           dev,
		DebugMode:        debug || dev,
		PasswordTTL:      30 * time.Minute,
		ActivationTTL:    24 * time.Hour,
		PageSize:         24,
		MaxSearchResults: 12,
		Origin:           env("ORIGIN", dev),
		BaseImgUrl:       env("IMG_URL", dev),
	}
	return config
}
