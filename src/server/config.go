// Konfigurations- und Initialisierungslogik für den Server
// Verantwortlich für das Einlesen von Umgebungsvariablen, Flags und Konfigurationsdateien

package server

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

// ServerConfig hält alle konfigurierbaren Parameter
// Diese Struktur kann bei Bedarf erweitert werden

type ServerConfig struct {
	DataDir   string `json:"data_dir"`
	Port      string `json:"port"`
	LogFolder string `json:"log_folder"`
	ServerLog string `json:"server_log"`
	PidFile   string `json:"pid_file"`
}

// LoadConfig liest die Konfiguration aus einer JSON-Datei oder verwendet Standardwerte
func LoadConfig(path string) ServerConfig {
	config := ServerConfig{
		DataDir:   "./data",
		Port:      "2022",
		LogFolder: "./logs",
		ServerLog: "coco-db.log",
		PidFile:   "coco-db.pid",
	}
	file, err := ioutil.ReadFile(path)
	if err == nil {
		_ = json.Unmarshal(file, &config)
	}
	return config
}

// InitLogging initialisiert das Logging in eine Datei
func InitLogging(logFile string) {
	logDir := filepath.Dir(logFile)
	if logDir != "." && logDir != "" {
		err := os.MkdirAll(logDir, 0755)
		if err != nil {
			log.Fatalf("Fehler beim Anlegen des Logverzeichnisses: %v", err)
		}
	}
	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Fehler beim Öffnen der Logdatei: %v", err)
	}
	log.SetOutput(file)
}

// DefaultConfig gibt die Standardkonfiguration zurück
func DefaultConfig() ServerConfig {
	return ServerConfig{
		DataDir:   "./data",
		Port:      "2022",
		LogFolder: "./logs",
		ServerLog: "coco-db.log",
		PidFile:   "coco-db.pid",
	}
}

// WriteDefaultConfig schreibt die Default-Konfiguration als config.json an den angegebenen Pfad
func WriteDefaultConfig(path string) error {
	cfg := DefaultConfig()
	b, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, b, 0644)
}
