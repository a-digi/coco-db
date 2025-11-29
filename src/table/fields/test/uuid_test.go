package fields_test

import (
	"regexp"
	"testing"

	fields "github.com/a-digi/coco-db/src/table/fields"
)

func TestNewUUIDv4_GeneratesValidUUID(t *testing.T) {
	uuid, err := fields.NewUUIDv4()
	if err != nil {
		t.Fatalf("Fehler beim Generieren der UUID: %v", err)
	}
	// Prüfe Format: 8-4-4-4-12 hex Zeichen
	re := regexp.MustCompile(`^[a-f0-9]{8}-[a-f0-9]{4}-4[a-f0-9]{3}-[89ab][a-f0-9]{3}-[a-f0-9]{12}$`)
	if !re.MatchString(uuid) {
		t.Errorf("Ungültiges UUIDv4-Format: %s", uuid)
	}
}

func TestNewUUIDv4_Unique(t *testing.T) {
	uuidSet := make(map[string]struct{})
	for i := 0; i < 1000; i++ {
		uuid, err := fields.NewUUIDv4()
		if err != nil {
			t.Fatalf("Fehler beim Generieren der UUID: %v", err)
		}
		if _, exists := uuidSet[uuid]; exists {
			t.Errorf("Doppelte UUID generiert: %s", uuid)
		}
		uuidSet[uuid] = struct{}{}
	}
}
