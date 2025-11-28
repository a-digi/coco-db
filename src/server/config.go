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
		Port:      "8080",
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
		_ = os.MkdirAll(logDir, 0755)
	}
	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Fehler beim Öffnen der Logdatei: %v", err)
	}
	log.SetOutput(file)
}
