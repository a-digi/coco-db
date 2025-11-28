package logger

// Logger ist ein Interface für Logging-Implementierungen
// (z.B. FileLogger, später auch ConsoleLogger, etc.)
type Logger interface {
	Info(v ...interface{})
	Error(v ...interface{})
	Close() error
}

// Ensure FileLogger implements Logger
var _ Logger = (*FileLogger)(nil)

// Beispiel für die Nutzung:
// logger, _ := logger.NewFileLogger("log.txt")
// logger.Info("Server gestartet")
// logger.Error("Fehler: ", err)
// logger.Close()
