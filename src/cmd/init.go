package cmd

import (
	"flag"
	"fmt"
	"os"
)

// InitDatabase initialisiert das Datenverzeichnis für coco-db.
// Sie prüft, ob das Verzeichnis existiert, legt es ggf. an und prüft die Schreibrechte.
func InitDatabase(dataDir string) error {
	if dataDir == "" {
		return fmt.Errorf("Datenverzeichnis (--data-dir) muss angegeben werden")
	}
	info, err := os.Stat(dataDir)
	if os.IsNotExist(err) {
		err = os.MkdirAll(dataDir, 0755)
		if err != nil {
			return fmt.Errorf("Konnte Datenverzeichnis nicht anlegen: %w", err)
		}
		fmt.Println("Datenverzeichnis angelegt:", dataDir)
	} else if err != nil {
		return fmt.Errorf("Fehler beim Zugriff auf Datenverzeichnis: %w", err)
	} else if !info.IsDir() {
		return fmt.Errorf("Pfad existiert, ist aber kein Verzeichnis: %s", dataDir)
	}
	// Teste Schreibrechte
	testFile := dataDir + "/.coco-db-init-test"
	f, err := os.Create(testFile)
	if err != nil {
		return fmt.Errorf("Keine Schreibrechte im Datenverzeichnis: %w", err)
	}
	f.Close()
	os.Remove(testFile)
	fmt.Println("Datenverzeichnis ist schreibbar.")
	return nil
}

// RunInitCommand parst die Flags und führt die Initialisierung aus.
func RunInitCommand() {
	dataDir := flag.String("data-dir", "", "Pfad zum Datenverzeichnis")
	flag.Parse()
	if err := InitDatabase(*dataDir); err != nil {
		fmt.Fprintln(os.Stderr, "Fehler bei der Initialisierung:", err)
		os.Exit(1)
	}
	fmt.Println("Initialisierung erfolgreich abgeschlossen.")
}
