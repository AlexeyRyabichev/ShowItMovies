package internal

import (
	"database/sql"
	"fmt"
	"github.com/lib/pq"
	"log"
	"os"
)

var username = os.Getenv("DBUSER")
var password = os.Getenv("DBPASS")

func GetWatchlist(login string) MovieWatchlistHTTP {
	connStr := fmt.Sprintf("user=%s password=%s dbname=showit sslmode=disable", username, password)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Printf("ERR\t%v", err)
	}
	defer db.Close()

	var watchlist MovieWatchlistHTTP
	err = db.QueryRow(fmt.Sprintf(`select seen_movies, unseen_movies from movies where login='%s'`, login)).Scan((*pq.StringArray)(&watchlist.SeenMovies), (*pq.StringArray)(&watchlist.UnseenMovies))
	if err != nil {
		log.Printf("ERR\tcannot get watchlist from db: %v", err)
		return MovieWatchlistHTTP{}
	}

	watchlist.Login = login
	return watchlist
}

func UpdateWatchlist(watchlist *MovieWatchlistHTTP) bool {
	connStr := fmt.Sprintf("user=%s password=%s dbname=showit sslmode=disable", username, password)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Printf("ERR\t%v", err)
	}
	defer db.Close()

	if _, err := db.Exec("update movies set seen_movies=$1, unseen_movies=$2 where login=$3", pq.Array(watchlist.SeenMovies), pq.Array(watchlist.UnseenMovies), watchlist.Login); err != nil{
		log.Printf("cannot insert watchlist in table: %v", err)
		return false
	}
	return true
}
