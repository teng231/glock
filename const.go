package glock

import "time"

const (
	StatusLocked         = "locked"
	StatusNotPersist     = "error_persist"
	StatusInvalidCounter = "invalid_counter"
	StatusCountdownDone  = "countdown_done"
	Empty                = "empty"
	CantParse            = "can't parse"
)

type ConnectConfig struct {
	RedisAddr string
	RedisPw   string
	Prefix    string
	Timelock  time.Duration
	RedisDb   int
}
