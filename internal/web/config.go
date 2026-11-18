package web

import (
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
