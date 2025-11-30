package query

import "github.com/a-digi/coco-db/src/index"

// GetTestRegistry gibt die zentrale IndexRegistry für Tests zurück
func GetTestRegistry() *index.IndexRegistry {
	return index.GetRegistry()
}

