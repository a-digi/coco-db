package main

import (
	"fmt"
	"github.com/a-digi/coco-db/src/index"
	"github.com/a-digi/coco-db/src/server"
	"os"
	"os/exec"
	"strconv"
	"syscall"
	"time"
)

// main ist der Einstiegspunkt für die coco-db CLI.
func main() {
	if len(os.Args) > 1 && os.Args[1] == "init" {
		// Initialisiere das Datenverzeichnis
		configPath := "config.json"
		if _, err := os.Stat(configPath); os.IsNotExist(err) {
			fmt.Println("[INFO] config.json nicht gefunden, lege Default-Konfiguration an:", configPath)
			if err := server.WriteDefaultConfig(configPath); err != nil {
				fmt.Println("[FATAL] Konnte config.json nicht anlegen:", err)
				os.Exit(1)
			}
		}
		fmt.Println("Initialisierung abgeschlossen.")
		return
	}

	if len(os.Args) > 1 && os.Args[1] == "start" {
		startServerBackground()
		return
	}

	if len(os.Args) > 1 && os.Args[1] == "stop" {
		stopServer()
		return
	}

	fmt.Println("Verwendung: coco-db <init|start|stop>")
}

func startServerBackground() {
	config := server.LoadConfig("config.json")
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
	time.Sleep(1 * time.Second)
	if err := cmd.Process.Signal(syscall.Signal(0)); err != nil {
		fmt.Println("Warnung: Serverprozess ist direkt nach dem Start nicht mehr aktiv (Port belegt oder Fehler). PID-Datei wird entfernt.")
		os.Remove(config.PidFile)
		os.Exit(1)
	}
	os.Exit(0)
}

func stopServer() {
	config := server.LoadConfig("config.json")
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
	for i := 0; i < 10; i++ {
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

func processExists(pid int) bool {
	proc, err := os.FindProcess(pid)
	if err != nil {
		return false
	}
	return proc.Signal(syscall.Signal(0)) == nil
}

func init() {
	if len(os.Args) > 1 && os.Args[1] == "_run" {
		config := server.LoadConfig("config.json")
		server.StartServerWithConfig(config)

		// Nach dem Laden der Indizes: RAM-Debug-Ausgabe für wichtige Indizes
		index.GetRegistry().DebugPrintIndex("poseidon.users.age", 10)
		index.GetRegistry().DebugPrintIndex("poseidon.users.date", 10)

		os.Exit(0)
	}
}
