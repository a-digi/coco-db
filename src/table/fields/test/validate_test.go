package fields_test

import (
	"testing"
	"github.com/a-digi/coco-db/src/table/fields/types"
	fieldsmod "github.com/a-digi/coco-db/src/table/fields"
)

func TestValidateEntry_StringConstraints(t *testing.T) {
	meta := types.TableMeta{
		Fields: []types.FieldMeta{{
			Name:      "username",
			Type:      "string",
			Required:  true,
			MinLength: intPtr(3),
			MaxLength: intPtr(8),
			Pattern:   "^[a-z]+$",
			Enum:      []interface{}{ "alice", "bob", "carol" },
		}},
	}
	// Valid
	entry := map[string]interface{}{ "username": "alice" }
	errs := fieldsmod.ValidateEntry(entry, meta)
	if len(errs) != 0 {
		t.Errorf("Valid string failed: %v", errs)
	}
	// Too short
	entry = map[string]interface{}{ "username": "al" }
	errs = fieldsmod.ValidateEntry(entry, meta)
	if len(errs) == 0 || errs[0].Code != "ERR_MIN_LENGTH" {
		t.Errorf("Expected ERR_MIN_LENGTH, got %v", errs)
	}
	// Too long
	entry = map[string]interface{}{ "username": "alicebobx" }
	errs = fieldsmod.ValidateEntry(entry, meta)
	if len(errs) == 0 || errs[0].Code != "ERR_MAX_LENGTH" {
		t.Errorf("Expected ERR_MAX_LENGTH, got %v", errs)
	}
	// Pattern
	entry = map[string]interface{}{ "username": "Alice" }
	errs = fieldsmod.ValidateEntry(entry, meta)
	if len(errs) == 0 || errs[0].Code != "ERR_PATTERN" {
		t.Errorf("Expected ERR_PATTERN, got %v", errs)
	}
	// Enum
	entry = map[string]interface{}{ "username": "eve" }
	errs = fieldsmod.ValidateEntry(entry, meta)
	if len(errs) == 0 || errs[0].Code != "ERR_ENUM" {
		t.Errorf("Expected ERR_ENUM, got %v", errs)
	}
}

func intPtr(i int) *int { return &i }

