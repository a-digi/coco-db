package logger

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// FileLogger ist ein thread-sicherer Logger, der Fehler und normale Logs trennt
// Fehler gehen in error.log, alles andere in db.log (jeweils nach Jahr/Monat/Tag)
type FileLogger struct {
	baseDir    string
	dbLogger   *log.Logger
	dbFile     *os.File
	errorLogger *log.Logger
	errorFile   *os.File
	mu         sync.Mutex
}

// NewFileLogger nimmt das Basisverzeichnis (z.B. logs/) und erzeugt Logger für db.log und error.log
func NewFileLogger(baseDir string) (*FileLogger, error) {
	logger := &FileLogger{baseDir: baseDir}
	if err := logger.rotate(); err != nil {
		return nil, err
	}
	return logger, nil
}

// rotate öffnet die aktuellen Logdateien für das heutige Datum
func (f *FileLogger) rotate() error {
	f.closeFiles()
	now := time.Now()
	dir := filepath.Join(f.baseDir,
		fmt.Sprintf("%04d", now.Year()),
		fmt.Sprintf("%02d", int(now.Month())),
		fmt.Sprintf("%02d", now.Day()),
	)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	dbPath := filepath.Join(dir, "db.log")
	errorPath := filepath.Join(dir, "error.log")
	dbFile, err := os.OpenFile(dbPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	errorFile, err := os.OpenFile(errorPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		dbFile.Close()
		return err
	}
	f.dbFile = dbFile
	f.errorFile = errorFile
	f.dbLogger = log.New(dbFile, "", log.LstdFlags|log.Lshortfile)
	f.errorLogger = log.New(errorFile, "", log.LstdFlags|log.Lshortfile)
	return nil
}

// closeFiles schließt die Logdateien, falls offen
func (f *FileLogger) closeFiles() {
	if f.dbFile != nil {
		f.dbFile.Close()
		f.dbFile = nil
	}
	if f.errorFile != nil {
		f.errorFile.Close()
		f.errorFile = nil
	}
}

// logToFile schreibt je nach Level in die richtige Datei
func (f *FileLogger) logToFile(level string, v ...interface{}) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.rotate() // Optional: Bei jedem Log-Aufruf prüfen, ob ein neuer Tag ist
	msg := fmt.Sprint(v...)
	var logger *log.Logger
	switch level {
	case "ERROR", "CRITICAL", "ALERT", "EMERGENCY":
		logger = f.errorLogger
	default:
		logger = f.dbLogger
	}
	logger.SetPrefix(level + ": ")
	logger.Output(3, msg)
}

func (f *FileLogger) Log(v ...interface{})      { f.logToFile("LOG", v...) }
func (f *FileLogger) Debug(v ...interface{})    { f.logToFile("DEBUG", v...) }
func (f *FileLogger) Info(v ...interface{})     { f.logToFile("INFO", v...) }
func (f *FileLogger) Notice(v ...interface{})   { f.logToFile("NOTICE", v...) }
func (f *FileLogger) Warning(v ...interface{})  { f.logToFile("WARNING", v...) }
func (f *FileLogger) Error(v ...interface{})    { f.logToFile("ERROR", v...) }
func (f *FileLogger) Critical(v ...interface{}) { f.logToFile("CRITICAL", v...) }
func (f *FileLogger) Alert(v ...interface{})    { f.logToFile("ALERT", v...) }
func (f *FileLogger) Emergency(v ...interface{}) { f.logToFile("EMERGENCY", v...) }

// Close schließt die Logdateien
func (f *FileLogger) Close() error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.closeFiles()
	return nil
}
