package internal

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

func (rt *Router) PostWatchlist(w http.ResponseWriter, r *http.Request) {
	var movie MovieHTTP
	if err := json.NewDecoder(r.Body).Decode(&movie); err != nil {
		log.Printf("ERR\t%v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	watchlist := GetWatchlist(movie.Login)

	if stringInSlice(movie.IMDBId, watchlist.SeenMovies) {
		log.Printf("RESP\tPOST\tmovie %s already in seen list for user %s", movie.IMDBId, movie.Login)
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

func (rt *Router) DeleteWatchlist(w http.ResponseWriter, r *http.Request) {
	login := r.Header.Get("X-Login")
	movieID := r.Header.Get("X-IMDBId")
	fromSeen, _ := strconv.ParseBool(r.Header.Get("X-Seen"))
	fromUnseen, _ := strconv.ParseBool(r.Header.Get("X-Unseen"))

	watchlist := GetWatchlist(login)

	if fromSeen {
		watchlist.SeenMovies = removeFromSlice(movieID, watchlist.SeenMovies)
	}
	if fromUnseen{
		watchlist.UnseenMovies = removeFromSlice(movieID, watchlist.UnseenMovies)
	}

	if !UpdateWatchlist(&watchlist) {
		log.Printf("RESP\tPOST\tcannot remove movie from watchlist, user %s, movie %s", watchlist.Login, movieID)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (rt *Router) GetWatchlist(w http.ResponseWriter, r *http.Request) {
	login := r.Header.Get("X-Login")

	watchlist := GetWatchlist(login)

	js, err := json.Marshal(&watchlist)
	if err != nil {
		log.Printf("ERR\tcannot parse watchlist to json, %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	if bytes, err := w.Write(js); err != nil {
		log.Printf("ERR\t%v", err)
		http.Error(w, "ERR\tcannot write json to response: "+err.Error(), http.StatusInternalServerError)
		return
	} else {
		log.Printf("RESP\tGET\twritten %d bytes in response", bytes)
	}
}

func (rt *Router) GetMovie(w http.ResponseWriter, r *http.Request) {
	login := r.Header.Get("X-Login")
	movieID := r.Header.Get("X-IMDBId")

	watchlist := GetWatchlist(login)

	movieInfo := MovieHTTP{
		Login:    login,
		Password: "",
		IMDBId:   "",
		Seen:     stringInSlice(movieID, watchlist.SeenMovies),
		Unseen:   stringInSlice(movieID, watchlist.UnseenMovies),
	}

	js, err := json.Marshal(&movieInfo)
	if err != nil {
		log.Printf("ERR\tcannot parse movie to json, %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	if bytes, err := w.Write(js); err != nil {
		log.Printf("ERR\t%v", err)
		http.Error(w, "ERR\tcannot write json to response: "+err.Error(), http.StatusInternalServerError)
		return
	} else {
		log.Printf("RESP\tGET\twritten %d bytes in response", bytes)
	}
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
			list[i] = list[len(list)-1]
			return list[:len(list)-1]
		}
	}
	return list
}
