package internal

type MovieHTTP struct {
	Login    string `json:"login"`
	Password string `json:"password"`
	IMDBId   string `json:"imdb_id"`
}

type MovieWatchlistHTTP struct {
	Login        string   `json:"login"`
	SeenMovies   []string `json:"seen_movies"`
	UnseenMovies []string `json:"unseen_movies"`
}