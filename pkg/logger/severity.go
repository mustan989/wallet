package logger

import (
	"encoding/json"
)

type Severity uint8

const (
	_ Severity = iota
	Debug
	Info
	Warning
	Error
	Fatal
)

func (s Severity) String() string {
	return severities[s]
}

func (s Severity) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}

func (s Severity) MarshalText() (text []byte, err error) {
	return []byte(severities[s]), nil
}

var severities = map[Severity]string{
	Debug:   "DEBUG",
	Info:    "INFO",
	Warning: "WARN",
	Error:   "ERROR",
	Fatal:   "FATAL",
}
