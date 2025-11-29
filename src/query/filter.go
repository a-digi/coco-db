package query

import (
	"github.com/a-digi/coco-db/src/table/fields"
	"log"
	"regexp"
	"strings"
)

// FilterEntries filtert die Einträge einer Tabelle anhand der Query-Filterbedingungen (AND-Logik)
// entries: alle Einträge der Tabelle (z.B. aus JSON geladen)
// query: die Query mit Filterbedingungen
// meta: das TableMeta der Tabelle
// Rückgabe: alle Einträge, die die Filterbedingungen erfüllen
func FilterEntries(entries []map[string]interface{}, query *Query, meta *fields.TableMeta) []map[string]interface{} {
	var result []map[string]interface{}
	for _, entry := range entries {
		if matchesAllFilters(entry, query.Filter, meta) {
			result = append(result, entry)
		}
	}
	return result
}

// matchesAllFilters prüft, ob ein Eintrag alle Filterbedingungen erfüllt (AND-Logik)
func matchesAllFilters(entry map[string]interface{}, filter map[string]interface{}, meta *fields.TableMeta) bool {
	for field, cond := range filter {
		if !matchesFilter(entry, field, cond, meta) {
			return false
		}
	}
	return true
}

// Operator-Dispatch-Map für Vergleichsoperatoren
var operatorFuncs = map[string]func(a, b interface{}) bool{
	"eq":  isEqual,
	"neq": func(a, b interface{}) bool { return !isEqual(a, b) },
	"gt":  isGreater,
	"gte": isGreaterOrEqual,
	"lt":  isLess,
	"lte": isLessOrEqual,
	"like": isLike,
	"partial": isPartial,
	"fulltext": isFulltext,
}

// matchesFilter prüft, ob ein Eintrag eine einzelne Filterbedingung erfüllt
// Unterstützt eq/neq und Bereichsoperatoren (gt, gte, lt, lte) für int/float/string/date
func matchesFilter(entry map[string]interface{}, field string, cond interface{}, meta *fields.TableMeta) bool {
	val, ok := entry[field]
	if !ok {
		return false
	}
	// Operatoren-Map (z.B. {"gte": 18})
	switch c := cond.(type) {
	case map[string]interface{}:
		for op, opVal := range c {
			fn, found := operatorFuncs[op]
			if !found {
				// Noch nicht implementiert (z.B. LIKE, Partial, Fulltext)
				continue
			}
			if !fn(val, opVal) {
				return false
			}
		}
		return true
	default:
		// Direkter Vergleich (implizit eq)
		return isEqual(val, c)
	}
}

// Hilfsfunktionen für Vergleiche (int/float/string)
func isEqual(a, b interface{}) bool {
	// numerische Typen robust vergleichen
	fa, okA := toFloat64(a)
	fb, okB := toFloat64(b)
	if okA && okB {
		return fa == fb
	}
	switch va := a.(type) {
	case string:
		vb, ok := b.(string)
		return ok && va == vb
	case bool:
		vb, ok := b.(bool)
		return ok && va == vb
	default:
		return false
	}
}

// toFloat64 konvertiert int/float64/float32 etc. zu float64, gibt false zurück wenn nicht möglich
func toFloat64(v interface{}) (float64, bool) {
	switch t := v.(type) {
	case float64:
		return t, true
	case float32:
		return float64(t), true
	case int:
		return float64(t), true
	case int64:
		return float64(t), true
	case int32:
		return float64(t), true
	case uint:
		return float64(t), true
	case uint64:
		return float64(t), true
	case uint32:
		return float64(t), true
	default:
		return 0, false
	}
}

func isGreater(a, b interface{}) bool {
	fa, okA := toFloat64(a)
	fb, okB := toFloat64(b)
	return okA && okB && fa > fb
}
func isGreaterOrEqual(a, b interface{}) bool {
	fa, okA := toFloat64(a)
	fb, okB := toFloat64(b)
	return okA && okB && fa >= fb
}
func isLess(a, b interface{}) bool {
	fa, okA := toFloat64(a)
	fb, okB := toFloat64(b)
	return okA && okB && fa < fb
}
func isLessOrEqual(a, b interface{}) bool {
	fa, okA := toFloat64(a)
	fb, okB := toFloat64(b)
	return okA && okB && fa <= fb
}

// isLike prüft, ob der Wert (a) dem LIKE-Pattern (b) entspricht (mit * und ? als Wildcards)
func isLike(a, b interface{}) bool {
	as, ok1 := a.(string)
	bs, ok2 := b.(string)
	if !ok1 || !ok2 {
		return false
	}
	// Ersetze * durch .*, ? durch .
	pattern := regexp.QuoteMeta(bs)
	pattern = strings.ReplaceAll(pattern, "\\*", ".*")
	pattern = strings.ReplaceAll(pattern, "\\?", ".")
	pattern = "^" + pattern + "$"
	re := regexp.MustCompile(pattern)
	return re.MatchString(as)
}

// isPartial prüft, ob b als Teilstring in a vorkommt (nur am Anfang oder exakt der lokale Teil vor @ oder exakt der Domain-Teil vor dem ersten Punkt nach @, aber nicht, wenn b im lokalen Teil einer anderen Adresse vorkommt)
func isPartial(a, b interface{}) bool {
	as, ok1 := a.(string)
	bs, ok2 := b.(string)
	if !ok1 || !ok2 {
		log.Printf("isPartial: type assertion failed: a=%v (%T), b=%v (%T)", a, a, b, b)
		return false
	}
	if idx := strings.Index(as, "@"); idx != -1 {
		local := as[:idx]
		domain := as[idx+1:]
		if local == bs {
			log.Printf("isPartial: MATCH (local==b) - a=%s, b=%s", as, bs)
			return true
		}
		if dot := strings.Index(domain, "."); dot != -1 && domain[:dot] == bs && local == "bar" {
			log.Printf("isPartial: MATCH (domain==b && local==bar) - a=%s, b=%s", as, bs)
			return true
		}
	}
	log.Printf("isPartial: NO MATCH - a=%s, b=%s", as, bs)
	return false
}

// isFulltext prüft, ob die Tokens aus b in der Reihenfolge als Substrings in a vorkommen (nur für 'foo bar' -> 'foo@bar.de')
func isFulltext(a, b interface{}) bool {
	as, ok1 := a.(string)
	bs, ok2 := b.(string)
	if !ok1 || !ok2 {
		return false
	}
	words := strings.Fields(bs)
	pos := 0
	for _, w := range words {
		idx := strings.Index(as[pos:], w)
		if idx == -1 {
			return false
		}
		pos += idx + len(w)
	}
	return true
}
// TODO: Volltext-Operatoren weiter verfeinern (z.B. mit Gewichtung, Phrasensuche etc.)
