package bus

import "time"

type BusConnectionConfig struct {
	Dsn        string
	MaxRetries int
	Timeout    time.Duration
}

type BusConfig struct {
	Connection BusConnectionConfig
}
