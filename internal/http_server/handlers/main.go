package handlers

import (
	"fmt"
	"net/http"
)

func DefaultHeandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	Body := "Server running!\n"
	fmt.Fprintf(w, "%s", Body)
}
