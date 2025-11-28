package logger

type Logger interface {
	Log(v ...interface{})
	Debug(v ...interface{})
	Info(v ...interface{})
	Notice(v ...interface{})
	Warning(v ...interface{})
	Error(v ...interface{})
	Critical(v ...interface{})
	Alert(v ...interface{})
	Emergency(v ...interface{})
}

type NoopLogger struct{}

func (n *NoopLogger) Log(v ...interface{})      {}
func (n *NoopLogger) Debug(v ...interface{})    {}
func (n *NoopLogger) Info(v ...interface{})     {}
func (n *NoopLogger) Notice(v ...interface{})   {}
func (n *NoopLogger) Warning(v ...interface{})  {}
func (n *NoopLogger) Error(v ...interface{})    {}
func (n *NoopLogger) Critical(v ...interface{}) {}
func (n *NoopLogger) Alert(v ...interface{})    {}
func (n *NoopLogger) Emergency(v ...interface{}) {}

