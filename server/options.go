package server

import (
	"crypto/tls"
	"log"
)

// Options represents MQTT server options.
type Options struct {
	TLSConfig *tls.Config
	Logger    *log.Logger

	// DisableMonitor disables embedded disconnenction detector or so.
	DisableMonitor bool
}

// DefaultOptions is used as Server#Options for default.
var DefaultOptions = &Options{}
