package logger

import (
	"os"
	"testing"
)

func TestFileLogger(t *testing.T) {
	logPath := "testlog.log"
	logger, err := NewFileLogger(logPath)
	if err != nil {
		t.Fatalf("Fehler beim Erstellen des FileLogger: %v", err)
	}
	defer func() {
		logger.Close()
		os.Remove(logPath)
	}()

	logger.Info("Dies ist eine Info-Nachricht.")
	logger.Error("Dies ist eine Fehler-Nachricht.")
}

