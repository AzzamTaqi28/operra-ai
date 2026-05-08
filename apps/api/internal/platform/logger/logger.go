package logger

import (
	"log"
	"os"
)

func New(env string) *log.Logger {
	prefix := "operra-api "
	if env != "" {
		prefix = prefix + "[" + env + "] "
	}

	return log.New(os.Stdout, prefix, log.LstdFlags|log.Lmicroseconds)
}
