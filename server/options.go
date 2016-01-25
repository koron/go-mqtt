package server

import "log"

// Options represents MQTT server options.
type Options struct {
	Logger *log.Logger
}

// DefaultOptions is used as Server#Options for default.
var DefaultOptions = &Options{}
