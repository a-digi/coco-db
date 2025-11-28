// Konfigurations- und Initialisierungslogik für den Server
// Verantwortlich für das Einlesen von Umgebungsvariablen, Flags und Konfigurationsdateien

package server

import (
	"os"
	"log"
)

// ServerConfig hält alle konfigurierbaren Parameter
// Diese Struktur kann bei Bedarf erweitert werden

type ServerConfig struct {
	DataDir string
	Port    string
	LogFile string
}

// LoadConfig liest die Konfiguration aus Umgebungsvariablen oder Standardwerten
func LoadConfig() ServerConfig {
	config := ServerConfig{
		DataDir: os.Getenv("COCO_DB_DATA_DIR"),
		Port:    os.Getenv("COCO_DB_PORT"),
		LogFile: os.Getenv("COCO_DB_LOG_FILE"),
	}
	if config.DataDir == "" {
		config.DataDir = "./data"
	}
	if config.Port == "" {
		config.Port = "8080"
	}
	if config.LogFile == "" {
		config.LogFile = "coco-db.log"
	}
	return config
}

// InitLogging initialisiert das Logging in eine Datei
func InitLogging(logFile string) {
	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Fehler beim Öffnen der Logdatei: %v", err)
	}
	log.SetOutput(file)
}

