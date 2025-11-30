package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type RelationEntry struct {
	TargetID  string `json:"target_id"`
	RelatedID string `json:"related_id"`
}

func readIDs(entriesDir string) ([]string, error) {
	files, err := ioutil.ReadDir(entriesDir)
	if err != nil {
		return nil, err
	}
	ids := make([]string, 0, len(files))
	for _, f := range files {
		if f.IsDir() {
			// Suche nach <dir>/<dir>.json
			jsonPath := filepath.Join(entriesDir, f.Name(), f.Name()+".json")
			b, err := ioutil.ReadFile(jsonPath)
			if err != nil {
				continue
			}
			var obj map[string]interface{}
			if err := json.Unmarshal(b, &obj); err != nil {
				continue
			}
			idVal, ok := obj["id"]
			idStr, okStr := idVal.(string)
			if ok && okStr && idStr != "" {
				ids = append(ids, idStr)
			}
		}
	}
	return ids, nil
}

func readConfig(path string) (string, string) {
	cfg := struct {
		DataDir string `json:"data_dir"`
		Port    string `json:"port"`
	}{
		DataDir: "./data",
		Port:    "2022",
	}
	b, err := ioutil.ReadFile(path)
	if err == nil {
		_ = json.Unmarshal(b, &cfg)
	}
	return cfg.DataDir, cfg.Port
}

func parseTableField(arg string) (string, string, error) {
	parts := strings.SplitN(arg, ":", 2)
	if len(parts) != 2 {
		return "", "", fmt.Errorf("Ungültiges Argument: %s. Erwartet: tabelle:feldname", arg)
	}
	return parts[0], parts[1], nil
}

func main() {
	rand.Seed(time.Now().UnixNano())
	targetTableArg := flag.String("targetTable", "user_roles", "Zieltabelle, in die neue Einträge geschrieben werden (z.B. user_roles)")
	relationTableArg := flag.String("relationTable", "users:user_id", "Erste Relationstabelle und Feldname im Format tabelle:feldname (z.B. users:user_id)")
	secondRelationTableArg := flag.String("secondRelationTable", "roles:role_id", "Zweite Relationstabelle und Feldname im Format tabelle:feldname (z.B. roles:role_id)")
	dbArg := flag.String("db", "poseidon", "Datenbankname")
	amountArg := flag.Int("amount", 100, "Anzahl der zu generierenden Einträge")
	apiArg := flag.String("api", "", "API-Basis-URL (z.B. http://localhost:2022)")
	dataDirArg := flag.String("data", "", "Datenverzeichnis (Default aus config.json)")
	configArg := flag.String("config", "./config.json", "Pfad zur config.json")
	flag.Parse()

	dataDir, port := readConfig(*configArg)
	if *dataDirArg != "" {
		dataDir = *dataDirArg
	}
	apiURL := *apiArg
	if apiURL == "" {
		apiURL = "http://localhost:" + port
	}
	dbName := *dbArg
	targetTable := *targetTableArg

	relationTable, relationField, err := parseTableField(*relationTableArg)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	secondRelationTable, secondRelationField, err := parseTableField(*secondRelationTableArg)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	relationEntriesDir := filepath.Join(dataDir, dbName, relationTable, "entries")
	secondEntriesDir := filepath.Join(dataDir, dbName, secondRelationTable, "entries")

	relationIDs, err := readIDs(relationEntriesDir)
	if err != nil {
		fmt.Printf("Fehler beim Lesen der IDs aus %s: %v\n", relationEntriesDir, err)
		os.Exit(1)
	}
	secondIDs, err := readIDs(secondEntriesDir)
	if err != nil {
		fmt.Printf("Fehler beim Lesen der IDs aus %s: %v\n", secondEntriesDir, err)
		os.Exit(1)
	}

	url := fmt.Sprintf("%s/api/databases/%s/tables/%s/entries", apiURL, dbName, targetTable)
	success, fail := 0, 0
	var wg sync.WaitGroup
	jobs := make(chan int, *amountArg)
	results := make(chan bool, *amountArg)
	maxWorkers := 50
	if *amountArg < maxWorkers {
		maxWorkers = *amountArg
	}

	worker := func() {
		for i := range jobs {
			relID := relationIDs[rand.Intn(len(relationIDs))]
			secID := secondIDs[rand.Intn(len(secondIDs))]
			entry := map[string]interface{}{
				relationField:       relID,
				secondRelationField: secID,
			}
			body, _ := json.Marshal(entry)
			resp, err := http.Post(url, "application/json", bytes.NewReader(body))
			if err != nil {
				fmt.Printf("Fehler bei Request %d: %v\n", i, err)
				results <- false
				continue
			}
			respBody, _ := ioutil.ReadAll(resp.Body)
			if resp.StatusCode == 200 || resp.StatusCode == 201 {
				fmt.Printf("%d: OK %s\n", i, string(respBody))
				results <- true
			} else {
				fmt.Printf("Fehler %d: %s\n", i, string(respBody))
				results <- false
			}
			resp.Body.Close()
		}
		wg.Done()
	}

	for w := 0; w < maxWorkers; w++ {
		wg.Add(1)
		go worker()
	}
	for i := 1; i <= *amountArg; i++ {
		jobs <- i
	}
	close(jobs)
	for i := 0; i < *amountArg; i++ {
		if <-results {
			success++
		} else {
			fail++
		}
	}
	wg.Wait()
	fmt.Printf("Fertig: %d erfolgreich, %d Fehler\n", success, fail)
}
