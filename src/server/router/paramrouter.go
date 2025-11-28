package router

import (
	"net/http"
	"strings"
)

// ParamHandlerFunc ist die Handler-Signatur für parameterisierte Routen
// params enthält die extrahierten Pfadparameter
// Beispiel: func(w http.ResponseWriter, r *http.Request, params map[string]string)
type ParamHandlerFunc func(http.ResponseWriter, *http.Request, map[string]string)

// routeEntry speichert das Pattern und den Handler
type routeEntry struct {
	method  string
	pattern string
	handler ParamHandlerFunc
}

// ParamRouter ist ein eigener Router für parameterbasierte URLs
// Er unterstützt Patterns wie /api/databases/{dbname}/tables/{tname}
type ParamRouter struct {
	routes []routeEntry
}

// NewParamRouter erzeugt einen neuen ParamRouter
func NewParamRouter() *ParamRouter {
	return &ParamRouter{routes: []routeEntry{}}
}

// HandleFunc registriert einen Handler für ein Pattern und eine HTTP-Methode
func (pr *ParamRouter) HandleFunc(method, pattern string, handler ParamHandlerFunc) {
	pr.routes = append(pr.routes, routeEntry{method, pattern, handler})
}

// ServeHTTP implementiert http.Handler und matched die Route
func (pr *ParamRouter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for _, route := range pr.routes {
		if r.Method != route.method {
			continue
		}
		params, ok := matchPattern(route.pattern, r.URL.Path)
		if ok {
			route.handler(w, r, params)
			return
		}
	}
	http.NotFound(w, r)
}

// matchPattern prüft, ob path zum pattern passt und extrahiert Parameter
func matchPattern(pattern, path string) (map[string]string, bool) {
	pSeg := strings.Split(strings.Trim(pattern, "/"), "/")
	uSeg := strings.Split(strings.Trim(path, "/"), "/")
	if len(pSeg) != len(uSeg) {
		return nil, false
	}
	params := make(map[string]string)
	for i := range pSeg {
		if strings.HasPrefix(pSeg[i], "{") && strings.HasSuffix(pSeg[i], "}") {
			key := pSeg[i][1:len(pSeg[i])-1]
			params[key] = uSeg[i]
		} else if pSeg[i] != uSeg[i] {
			return nil, false
		}
	}
	return params, true
}

