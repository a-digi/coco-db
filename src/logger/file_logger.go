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

// Info schreibt eine Info-Nachricht ins Log
func (f *FileLogger) Info(v ...interface{}) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.logger.SetPrefix("INFO: ")
	f.logger.Output(2, fmt.Sprint(v...))
}

// Error schreibt eine Fehler-Nachricht ins Log
func (f *FileLogger) Error(v ...interface{}) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.logger.SetPrefix("ERROR: ")
	f.logger.Output(2, fmt.Sprint(v...))
}

// Close schließt die Logdatei
func (f *FileLogger) Close() error {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.file.Close()
}
