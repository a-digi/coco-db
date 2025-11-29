package validate_test

import (
	"testing"
	validate "github.com/a-digi/coco-db/src/table/fields/validate"
)

func TestValidateNoIDField(t *testing.T) {
	cases := []struct {
		name  string
		entry map[string]interface{}
		hasErr bool
	}{
		{"no id field", map[string]interface{}{"foo": 1}, false},
		{"ID field upper", map[string]interface{}{"ID": 1}, true},
		{"id field lower", map[string]interface{}{"id": 1}, true},
		{"both id fields", map[string]interface{}{"id": 1, "ID": 2}, true},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			resp := validate.ValidateNoIDField(c.entry)
			if c.hasErr && resp == nil {
				t.Errorf("expected error, got nil")
			}
			if !c.hasErr && resp != nil {
				t.Errorf("expected no error, got: %v", resp)
			}
		})
	}
}

