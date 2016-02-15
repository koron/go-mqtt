package server

import (
	"crypto/tls"
	"log"
)

// Options represents MQTT server options.
type Options struct {
	TLSConfig *tls.Config
	Logger    *log.Logger
}

// DefaultOptions is used as Server#Options for default.
var DefaultOptions = &Options{}
