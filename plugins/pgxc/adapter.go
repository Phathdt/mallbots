package pgxc

import (
	"context"
	"github.com/jackc/pgx/v5/tracelog"
	sctx "github.com/phathdt/service-context"
	"strings"
)

const (
	ansiReset       = "\033[0m"
	ansiBrightBlue  = "\033[94m"
	ansiBrightGreen = "\033[92m"
	ansiBrightCyan  = "\033[96m"
)

type PgxLogAdapter struct {
	logger sctx.Logger
}

func colorizeQuery(msg string) string {
	msgLower := strings.ToLower(msg)
	if strings.Contains(msgLower, "select") {
		return ansiBrightBlue + msg + ansiReset
	} else if strings.Contains(msgLower, "insert") {
		return ansiBrightGreen + msg + ansiReset
	} else if strings.Contains(msgLower, "update") {
		return "\033[93m" + msg + ansiReset // Yellow
	} else if strings.Contains(msgLower, "delete") {
		return "\033[91m" + msg + ansiReset // Red
	}
	return ansiBrightCyan + msg + ansiReset
}

func (l *PgxLogAdapter) Log(ctx context.Context, level tracelog.LogLevel, msg string, data map[string]interface{}) {
	// Skip if message contains "prepare" (case insensitive)
	if strings.Contains(strings.ToLower(msg), "prepare") {
		return
	}

	coloredMsg := colorizeQuery(msg)

	switch level {
	case tracelog.LogLevelTrace:
		l.logger.Debugf("%s %v", coloredMsg, data)
	case tracelog.LogLevelDebug:
		l.logger.Debugf("%s %v", coloredMsg, data)
	case tracelog.LogLevelInfo:
		l.logger.Infof("%s %v", coloredMsg, data)
	case tracelog.LogLevelWarn:
		l.logger.Warnf("%s %v", coloredMsg, data)
	case tracelog.LogLevelError:
		l.logger.Errorf("%s %v", coloredMsg, data)
	default:
		l.logger.Infof("%s %v", coloredMsg, data)
	}
}
