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

type FieldMeta struct {
	Name      string        `json:"name"`
	Type      string        `json:"type"`
	MinLength *int          `json:"minLength,omitempty"`
	MaxLength *int          `json:"maxLength,omitempty"`
	Min       *float64      `json:"min,omitempty"`
	Max       *float64      `json:"max,omitempty"`
	Enum      []interface{} `json:"enum,omitempty"`
	Pattern   string        `json:"pattern,omitempty"`
	Required  bool          `json:"required,omitempty"`
}

type TableMeta struct {
	Fields []FieldMeta `json:"fields"`
}

func randomString() string {
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, 8)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func randomStringWithLength(minLen, maxLen int) string {
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	ln := minLen
	if maxLen > minLen {
		ln = rand.Intn(maxLen-minLen+1) + minLen
	}
	b := make([]rune, ln)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func randomValue(f FieldMeta, i int) interface{} {
	// Enum hat Vorrang
	if len(f.Enum) > 0 {
		return f.Enum[rand.Intn(len(f.Enum))]
	}
	minLen, maxLen := 3, 12
	if f.MinLength != nil {
		minLen = *f.MinLength
	}
	if f.MaxLength != nil {
		maxLen = *f.MaxLength
	}
	switch f.Type {
	case "string":
		if strings.Contains(strings.ToLower(f.Name), "email") {
			return fmt.Sprintf("%s%d@foo.de", randomStringWithLength(3, maxLen-8), i)
		}
		if strings.Contains(strings.ToLower(f.Name), "name") {
			return randomStringWithLength(minLen, maxLen)
		}
		return randomStringWithLength(minLen, maxLen)
	case "int", "integer":
		min, max := 0, 100
		if f.Min != nil {
			min = int(*f.Min)
		}
		if f.Max != nil {
			max = int(*f.Max)
		}
		if max > min {
			return rand.Intn(max-min+1) + min
		}
		return min
	case "float":
		min, max := 0.0, 100.0
		if f.Min != nil {
			min = *f.Min
		}
		if f.Max != nil {
			max = *f.Max
		}
		if max > min {
			return min + rand.Float64()*(max-min)
		}
		return min
	case "bool":
		return rand.Intn(2) == 0
	case "date":
		return time.Now().Add(-time.Duration(rand.Intn(1000*24))*time.Hour).Format(time.RFC3339)
	default:
		return nil
	}
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

func main() {
	rand.Seed(time.Now().UnixNano())
	tableArg := flag.String("table", "testdb/users", "Tabelle im Format dbname/tablename")
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

	parts := strings.SplitN(*tableArg, "/", 2)
	if len(parts) != 2 {
		fmt.Println("--table muss im Format dbname/tablename angegeben werden")
		os.Exit(1)
	}
	dbName, tableName := parts[0], parts[1]
	metaPath := filepath.Join(dataDir, dbName, tableName, "meta.json")
	metaBytes, err := ioutil.ReadFile(metaPath)
	if err != nil {
		fmt.Printf("meta.json nicht gefunden: %v\n", err)
		os.Exit(1)
	}
	var meta TableMeta
	if err := json.Unmarshal(metaBytes, &meta); err != nil {
		fmt.Printf("meta.json ungültig: %v\n", err)
		os.Exit(1)
	}

	url := fmt.Sprintf("%s/api/databases/%s/tables/%s/entries", apiURL, dbName, tableName)
	success, fail := 0, 0
	var wg sync.WaitGroup
	jobs := make(chan int, *amountArg)
	results := make(chan bool, *amountArg)
	maxWorkers := 100
	if *amountArg < maxWorkers {
		maxWorkers = *amountArg
	}

	// Worker-Funktion
	worker := func() {
		for i := range jobs {
			entry := make(map[string]interface{})
			for _, f := range meta.Fields {
				if strings.ToLower(f.Name) == "id" {
					continue
				}
				if f.Required || rand.Float64() < 0.8 {
					entry[f.Name] = randomValue(f, i)
				}
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

	// Starte Worker
	for w := 0; w < maxWorkers; w++ {
		wg.Add(1)
		go worker()
	}

	// Jobs verteilen
	for i := 1; i <= *amountArg; i++ {
		jobs <- i
	}
	close(jobs)

	// Ergebnisse sammeln
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
