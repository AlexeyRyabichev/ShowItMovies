package internal

import (
	"encoding/json"
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
	login := r.Header.Get("X-Login")
	movieID := r.Header.Get("X-IMDBId")

	watchlist := GetWatchlist(login)

	watchlist.SeenMovies = removeFromSlice(movieID ,watchlist.SeenMovies)
	watchlist.UnseenMovies = removeFromSlice(movieID ,watchlist.UnseenMovies)

	if !UpdateWatchlist(&watchlist) {
		log.Printf("RESP\tPOST\tcannot remove movie from watchlist, user %s, movie %s", watchlist.Login, movieID)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func removeFromSlice(a string, list []string) []string {
	for i, b := range list {
		if b == a {
			list[i] = list[len(list) - 1]
			return list[:len(list) - 1]
		}
	}
	return list
}
