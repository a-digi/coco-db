package validate

import "github.com/a-digi/coco-db/src/response"

// ValidateNoIDField prüft, ob das Feld "ID" oder "id" im Eintrag vorhanden ist.
// Gibt einen passenden APIResponse-Fehler zurück, falls das Feld existiert, sonst nil.
func ValidateNoIDField(entry map[string]interface{}) *response.APIResponse {
	if _, exists := entry["ID"]; exists {
		return response.WriteErrorInternal(400, "ERR_FORBIDDEN_ID_FIELD", "Das Feld 'ID' darf beim Anlegen nicht gesetzt werden. Es wird vom System vergeben.", "")
	}
	if _, exists := entry["id"]; exists {
		return response.WriteErrorInternal(400, "ERR_FORBIDDEN_ID_FIELD", "Das Feld 'id' darf beim Anlegen nicht gesetzt werden. Es wird vom System vergeben.", "")
	}
	return nil
}

