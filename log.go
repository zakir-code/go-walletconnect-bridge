package main

import (
	"os"

	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("go-walletconnect-bridge")

// Example format string. Everything except the message has a custom color
// which is dependent on the log level. Many fields have a custom output
// formatting too, eg. the time returns the hour down to the milli second.
var _format = logging.MustStringFormatter(`%{color}%{time:2006-01-02 15:04:05.000} ▶ %{level:.4s} %{id:03d}%{color:reset} %{message}`)

// For demo purposes, create two backend for os.Stderr.
var _backend = logging.NewLogBackend(os.Stderr, "", 0)

// For messages written to backend2 we want to add some additional
// information to the output, including the used log level and the name of
// the function.
var BackendFormatter = logging.NewBackendFormatter(_backend, _format)
