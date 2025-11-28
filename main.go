package main

import (
	"fmt"
	"github.com/a-digi/coco-db/src/cmd"
	"os"
)

// main ist der Einstiegspunkt für die coco-db CLI.
func main() {
	if len(os.Args) > 1 && os.Args[1] == "init" {
		// Entferne das Subkommando aus den Argumenten für die Flag-Parse
		os.Args = append([]string{os.Args[0]}, os.Args[2:]...)
		cmd.RunInitCommand()
		return
	}

	fmt.Println("coco-db gestartet. (Verwende 'init' für die Initialisierung des Datenverzeichnisses)")
}
