package log

import (
	"go.uber.org/zap/zapcore"
	"testing"
)


func TestNewLogger(t *testing.T) {
	sugarLogger, err := NewLogger(ConsoleEncoder, "", zapcore.DebugLevel)
	defer sugarLogger.Sync()
	if err != nil {
		t.Errorf("can not get new logger:\n%s\n", err.Error())
	}
	sugarLogger.Debugf("this is debug msg")
}
