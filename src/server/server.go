// Einstiegspunkt für den Server
// Verantwortlich für die Initialisierung und das Starten des HTTP-Servers

package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"
)

// StartServerWithConfig initialisiert und startet den HTTP-Server mit Konfiguration
func StartServerWithConfig(cfg ServerConfig) {
	// Logdateipfad ggf. mit LogFolder kombinieren, falls nicht absolut
	logFile := cfg.ServerLog
	if !filepath.IsAbs(logFile) && cfg.LogFolder != "" {
		logFile = filepath.Join(cfg.LogFolder, logFile)
	}
	InitLogging(logFile)
	addr := ":" + cfg.Port
	mux := SetupRouter()

	srv := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	// Graceful Shutdown vorbereiten
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-quit
		log.Println("Shutdown signal received, shutting down server...")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			log.Fatalf("Server Shutdown Failed: %v", err)
		}
		log.Println("Server gracefully stopped")
	}()

	log.Printf("Server läuft auf http://localhost%s\n", addr)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
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
