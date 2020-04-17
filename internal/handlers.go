package internal

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func (rt *Router) PostMovie(w http.ResponseWriter, r *http.Request) {
	var movie MovieHTTP
	if err := json.NewDecoder(r.Body).Decode(&movie); err != nil {
		log.Printf("ERR\t%v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	watchlist := GetWatchlist(movie.Login)

	if stringInSlice(movie.IMDBId, watchlist.SeenMovies) {
		log.Printf("movie %s already in seen list for user %s", movie.IMDBId, movie.Login)
		w.WriteHeader(http.StatusConflict)
		return
	} else {
		watchlist.SeenMovies = append(watchlist.SeenMovies, movie.IMDBId)
		if !UpdateWatchlist(&watchlist) {
			log.Printf("RESP\tPOST\tcannot add movie to seen watchlist, user %s, movie %s", watchlist.Login, movie.IMDBId)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
}

func (rt *Router) DeleteMovie(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World!")
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
