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

	logger.Log("Dies ist eine Log-Nachricht.")
	logger.Debug("Dies ist eine Debug-Nachricht.")
	logger.Info("Dies ist eine Info-Nachricht.")
	logger.Notice("Dies ist eine Notice-Nachricht.")
	logger.Warning("Dies ist eine Warning-Nachricht.")
	logger.Error("Dies ist eine Error-Nachricht.")
	logger.Critical("Dies ist eine Critical-Nachricht.")
	logger.Alert("Dies ist eine Alert-Nachricht.")
	logger.Emergency("Dies ist eine Emergency-Nachricht.")

	// Überprüfe, ob die Datei geschrieben wurde
	if _, err := os.Stat(logPath); err != nil {
		t.Errorf("Logdatei wurde nicht erstellt: %v", err)
	}
}
