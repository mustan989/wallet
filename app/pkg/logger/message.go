package logger

import (
	"time"
)

type message struct {
	Time     time.Time `json:"time"`
	Severity Severity  `json:"severity"`
	Message  string    `json:"message"`
}
