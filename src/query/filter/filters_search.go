package filter

// matchesAllFiltersSearch prüft, ob ein Eintrag alle Filterbedingungen für Suchanfragen erfüllt (AND-Logik)
func MatchesAllFiltersSearch(entry map[string]interface{}, filter map[string]interface{}, operatorFuncs map[string]func(interface{}, interface{}) bool, isEqual func(interface{}, interface{}) bool) bool {
	for field, cond := range filter {
		val, ok := entry[field]
		if !ok {
			return false
		}
		switch c := cond.(type) {
		case map[string]interface{}:
			for op, opVal := range c {
				fn, found := operatorFuncs[op]
				if !found {
					continue
				}
				if !fn(val, opVal) {
					return false
				}
			}
			continue
		default:
			if !isEqual(val, c) {
				return false
			}
		}
	}
	return true
}

