package vlog

import "testing"

func TestInitLog(t *testing.T) {
	Debug("test before initialized")
	InitLog(false)
	NoFormatLog("no prefix message")
	NoFormatLogf("no prefix message, %s", "test")
	Debug("debug message")
	Debugf("debug message with args, %s", "test")
	Info("info message")
	Infof("info message with args, %s", "test")
	Warn("warn message")
	Warnf("warn message with args, %s", "test")
	Error("error message")
	Errorf("error message with args, %s", "test")
	Fatal("fatal message")
}
