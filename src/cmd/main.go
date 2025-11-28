package cmd

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"syscall"
	"time"

	"coco-db/server"
)

var config server.ServerConfig

func ensureConfigFileExists(configPath string) error {
	if _, err := os.Stat(configPath); errors.Is(err, os.ErrNotExist) {
		fmt.Println("[INFO] config.json nicht gefunden, lege Default-Konfiguration an:", configPath)
		return server.WriteDefaultConfig(configPath)
	}
	return nil
}

func main() {
	execPath, err := os.Executable()
	if err != nil {
		fmt.Println("[FATAL] Kann Pfad zum Executable nicht bestimmen:", err)
		os.Exit(1)
	}
	configPath := filepath.Join(filepath.Dir(execPath), "config.json")
	if err := ensureConfigFileExists(configPath); err != nil {
		fmt.Println("[FATAL] Konnte config.json nicht anlegen:", err)
		os.Exit(1)
	}
	fmt.Println("[DEBUG] Lade Konfiguration aus:", configPath)
	config = server.LoadConfig(configPath)
	fmt.Printf("[DEBUG] Geladene Konfiguration: %+v\n", config)
	if len(os.Args) < 2 {
		fmt.Println("Verwendung: coco-db <start|stop>")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "start":
		startServerBackground()
	case "stop":
		stopServer()
	default:
		fmt.Println("Unbekanntes Kommando. Verwendung: coco-db <start|stop>")
		os.Exit(1)
	}
}

func startServerBackground() {
	if _, err := os.Stat(config.PidFile); err == nil {
		fmt.Println("Server läuft bereits (PID-Datei existiert)")
		os.Exit(1)
	}
	cmd := exec.Command(os.Args[0], "_run")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = nil
	cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
	if err := cmd.Start(); err != nil {
		fmt.Printf("Fehler beim Starten: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Server gestartet (PID %d)\n", cmd.Process.Pid)
	f, err := os.Create(config.PidFile)
	if err == nil {
		fmt.Fprintf(f, "%d", cmd.Process.Pid)
		f.Close()
	}
	// Warte kurz und prüfe, ob der Prozess noch läuft
	time.Sleep(1 * time.Second)
	if err := cmd.Process.Signal(syscall.Signal(0)); err != nil {
		fmt.Println("Warnung: Serverprozess ist direkt nach dem Start nicht mehr aktiv (Port belegt oder Fehler). PID-Datei wird entfernt.")
		os.Remove(config.PidFile)
		os.Exit(1)
	}
	os.Exit(0)
}

func stopServer() {
	pidBytes, err := os.ReadFile(config.PidFile)
	if err != nil {
		fmt.Println("PID-Datei nicht gefunden. Läuft der Server?")
		os.Exit(1)
	}
	pid, err := strconv.Atoi(string(pidBytes))
	if err != nil {
		fmt.Println("Ungültige PID in PID-Datei.")
		os.Remove(config.PidFile)
		os.Exit(1)
	}
	proc, err := os.FindProcess(pid)
	if err != nil {
		fmt.Printf("Prozess nicht gefunden: %v\n", err)
		os.Remove(config.PidFile)
		os.Exit(1)
	}
	if err := proc.Signal(syscall.SIGTERM); err != nil {
		fmt.Printf("Fehler beim Beenden: %v\n", err)
		os.Remove(config.PidFile)
		os.Exit(1)
	}
	// Warte auf tatsächliches Beenden des Prozesses
	for i := 0; i < 10; i++ { // bis zu 2 Sekunden warten
		time.Sleep(200 * time.Millisecond)
		if !processExists(pid) {
			break
		}
	}
	if processExists(pid) {
		fmt.Println("Warnung: Prozess konnte nicht beendet werden.")
		os.Exit(1)
	}
	os.Remove(config.PidFile)
	fmt.Println("Server gestoppt.")
}

// processExists prüft, ob ein Prozess mit der gegebenen PID existiert
func processExists(pid int) bool {
	proc, err := os.FindProcess(pid)
	if err != nil {
		return false
	}
	// Signal 0 prüft Existenz
	return proc.Signal(syscall.Signal(0)) == nil
}

// _run: interner Modus für den eigentlichen Serverprozess
func init() {
	if len(os.Args) > 1 && os.Args[1] == "_run" {
		execPath, err := os.Executable()
		if err != nil {
			fmt.Println("[FATAL] (_run) Kann Pfad zum Executable nicht bestimmen:", err)
			os.Exit(1)
		}
		configPath := filepath.Join(filepath.Dir(execPath), "config.json")
		if err := ensureConfigFileExists(configPath); err != nil {
			fmt.Println("[FATAL] (_run) Konnte config.json nicht anlegen:", err)
			os.Exit(1)
		}
		fmt.Println("[DEBUG] (_run) Lade Konfiguration aus:", configPath)
		config = server.LoadConfig(configPath)
		fmt.Printf("[DEBUG] (_run) Geladene Konfiguration: %+v\n", config)
		server.StartServerWithConfig(config)
		os.Exit(0)
	}
}
