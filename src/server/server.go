// Einstiegspunkt für den Server
// Verantwortlich für die Initialisierung und das Starten des HTTP-Servers

package server

import (
	"context"
	"fmt"
	"github.com/a-digi/coco-db/src/logger"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"syscall"
	"time"
)

// StartServerWithConfig initialisiert und startet den HTTP-Server mit Konfiguration
func StartServerWithConfig(cfg ServerConfig) {
	fmt.Println("[DEBUG] Starte Server mit Konfiguration:", cfg)
	var fileLogger *logger.FileLogger
	var err error
	// Datenverzeichnis anlegen, falls nicht vorhanden
	if cfg.DataDir != "" {
		err = os.MkdirAll(cfg.DataDir, 0755)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Fehler beim Anlegen des Datenverzeichnisses: %v\n", err)
			os.Exit(1)
		}
	}
	// Logdateipfad ggf. mit LogFolder kombinieren, falls nicht absolut
	logFile := cfg.ServerLog
	if !filepath.IsAbs(logFile) && cfg.LogFolder != "" {
		err = os.MkdirAll(cfg.LogFolder, 0755)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Fehler beim Anlegen des Logverzeichnisses: %v\n", err)
			os.Exit(1)
		}
		logFile = filepath.Join(cfg.LogFolder, logFile)
	}
	fmt.Println("[DEBUG] Logdatei:", logFile)
	fileLogger, err = logger.NewFileLogger(logFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fehler beim Initialisieren des FileLoggers: %v\n", err)
		os.Exit(1)
	}
	defer fileLogger.Close()
	fileLogger.Info("Server-Initialisierung gestartet.")
	fileLogger.Info("Verwende Datenverzeichnis:", cfg.DataDir)
	fileLogger.Info("Verwende Logdatei:", logFile)

	addr := ":" + cfg.Port
	mux := SetupRouter()

	srv := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	// PID-Datei schreiben
	if cfg.PidFile != "" {
		err = writePidFile(cfg.PidFile)
		if err != nil {
			fileLogger.Error("Fehler beim Schreiben der PID-Datei:", err)
		} else {
			fileLogger.Info("PID-Datei geschrieben:", cfg.PidFile)
		}
	}

	// Graceful Shutdown vorbereiten
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-quit
		fileLogger.Info("Shutdown signal received, shutting down server...")
		log.Println("Shutdown signal received, shutting down server...")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			fileLogger.Error("Server Shutdown Failed:", err)
			log.Fatalf("Server Shutdown Failed: %v", err)
		}
		fileLogger.Info("Server gracefully stopped")
		log.Println("Server gracefully stopped")
		if cfg.PidFile != "" {
			removePidFile(cfg.PidFile)
		}
	}()

	fileLogger.Info("Server läuft auf http://localhost" + addr)
	fmt.Printf("[DEBUG] Server läuft auf http://localhost%s\n", addr)
	log.Printf("Server läuft auf http://localhost%s\n", addr)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		fileLogger.Error("ListenAndServe():", err)
		fmt.Fprintf(os.Stderr, "ListenAndServe(): %v\n", err)
		log.Fatalf("ListenAndServe(): %v", err)
	}
}

// getEnv liest eine Umgebungsvariable oder gibt einen Defaultwert zurück
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func writePidFile(pidFile string) error {
	pid := os.Getpid()
	return os.WriteFile(pidFile, []byte(strconv.Itoa(pid)), 0644)
}

func removePidFile(pidFile string) {
	_ = os.Remove(pidFile)
}
