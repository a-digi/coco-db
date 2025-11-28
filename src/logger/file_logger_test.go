package logger

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestFileLogger(t *testing.T) {
	baseDir := "testlogs"
	logger, err := NewFileLogger(baseDir)
	if err != nil {
		t.Fatalf("Fehler beim Erstellen des FileLogger: %v", err)
	}
	defer func() {
		logger.Close()
		os.RemoveAll(baseDir)
	}()

	logger.Log("Dies ist eine Log-Nachricht.")
	logger.Debug("Dies ist eine Debug-Nachricht.")
	logger.Info("Dies ist eine Info-Nachricht.")
	logger.Notice("Dies ist eine Notice-Nachricht.")
	logger.Warning("Dies ist eine Warning-Nachricht.")
	logger.Error("Dies ist eine Error-Nachricht.")
	logger.Critical("Dies ist eine Critical-Nachricht.")
	logger.Alert("Dies ist eine Alert-Nachricht.")
	logger.Emergency("Dies ist eine Emergency-Nachricht.")

	now := time.Now()
	dir := filepath.Join(baseDir,
		formatInt(now.Year(), 4),
		formatInt(int(now.Month()), 2),
		formatInt(now.Day(), 2),
	)
	dbLogPath := filepath.Join(dir, "db.log")
	errorLogPath := filepath.Join(dir, "error.log")

	// Prüfe, ob db.log existiert und normale Logs enthält
	if _, err := os.Stat(dbLogPath); err != nil {
		t.Errorf("db.log wurde nicht erstellt: %v", err)
	} else {
		content, _ := os.ReadFile(dbLogPath)
		if !strings.Contains(string(content), "INFO") || !strings.Contains(string(content), "LOG") {
			t.Error("db.log enthält keine normalen Logs")
		}
	}

	// Prüfe, ob error.log existiert und Fehler enthält
	if _, err := os.Stat(errorLogPath); err != nil {
		t.Errorf("error.log wurde nicht erstellt: %v", err)
	} else {
		content, _ := os.ReadFile(errorLogPath)
		if !strings.Contains(string(content), "ERROR") || !strings.Contains(string(content), "CRITICAL") {
			t.Error("error.log enthält keine Fehler-Logs")
		}
	}
}

func formatInt(i, width int) string {
	s := "0000" + strconv.Itoa(i)
	return s[len(s)-width:]
}
