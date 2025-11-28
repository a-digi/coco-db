package logger

// Logger ist ein Enterprise-Interface für Logging-Implementierungen
// mit allen gängigen Log-Leveln.
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
	Close() error
}

// Beispiel für die Nutzung:
// logger, _ := logger.NewFileLogger("log.txt")
// logger.Info("Server gestartet")
// logger.Error("Fehler: ", err)
// logger.Close()
