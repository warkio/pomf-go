package pomf

import (
	"net/http"
)

func writeNotFound(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusNotFound)
	w.Write(nil)

	return nil
}

func writeMethodNotAllowed(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusMethodNotAllowed)
	w.Write(nil)

	return nil
}
