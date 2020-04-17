package internal

import (
	"fmt"
	"net/http"
)

func (rt *Router) PostMovie(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World!")
}

func (rt *Router) DeleteMovie(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World!")
}
