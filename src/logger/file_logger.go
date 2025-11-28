package logger

import (
	"fmt"
	"log"
	"os"
	"sync"
)

// FileLogger ist ein thread-sicherer Logger, der in eine Datei schreibt
// und das Standard-Log-Interface von Go nutzt.
type FileLogger struct {
	file   *os.File
	logger *log.Logger
	mu     sync.Mutex
}

// NewFileLogger öffnet (oder erstellt) die Logdatei und gibt einen FileLogger zurück
func NewFileLogger(path string) (*FileLogger, error) {
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}
	return &FileLogger{
		file:   file,
		logger: log.New(file, "", log.LstdFlags|log.Lshortfile),
	}, nil
}

// Log schreibt eine Log-Nachricht ins Log
func (f *FileLogger) Log(v ...interface{}) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.logger.SetPrefix("LOG: ")
	f.logger.Output(2, fmt.Sprint(v...))
}

// Debug schreibt eine Debug-Nachricht ins Log
func (f *FileLogger) Debug(v ...interface{}) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.logger.SetPrefix("DEBUG: ")
	f.logger.Output(2, fmt.Sprint(v...))
}

// Info schreibt eine Info-Nachricht ins Log
func (f *FileLogger) Info(v ...interface{}) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.logger.SetPrefix("INFO: ")
	f.logger.Output(2, fmt.Sprint(v...))
}

// Notice schreibt eine Notice-Nachricht ins Log
func (f *FileLogger) Notice(v ...interface{}) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.logger.SetPrefix("NOTICE: ")
	f.logger.Output(2, fmt.Sprint(v...))
}

// Warning schreibt eine Warnung ins Log
func (f *FileLogger) Warning(v ...interface{}) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.logger.SetPrefix("WARNING: ")
	f.logger.Output(2, fmt.Sprint(v...))
}

// Error schreibt eine Fehler-Nachricht ins Log
func (f *FileLogger) Error(v ...interface{}) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.logger.SetPrefix("ERROR: ")
	f.logger.Output(2, fmt.Sprint(v...))
}

// Critical schreibt eine kritische Nachricht ins Log
func (f *FileLogger) Critical(v ...interface{}) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.logger.SetPrefix("CRITICAL: ")
	f.logger.Output(2, fmt.Sprint(v...))
}

// Alert schreibt einen Alarm ins Log
func (f *FileLogger) Alert(v ...interface{}) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.logger.SetPrefix("ALERT: ")
	f.logger.Output(2, fmt.Sprint(v...))
}

// Emergency schreibt eine Notfall-Nachricht ins Log
func (f *FileLogger) Emergency(v ...interface{}) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.logger.SetPrefix("EMERGENCY: ")
	f.logger.Output(2, fmt.Sprint(v...))
}

// Close schließt die Logdatei
func (f *FileLogger) Close() error {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.file.Close()
}
