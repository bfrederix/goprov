package voteprov

import (
    "net/http"
	"github.com/gorilla/mux"
)


func CreateHandler() *mux.Router {
	r := mux.NewRouter()
	// Versioned API
	s := r.PathPrefix("/v1").Subrouter()
	s.HandleFunc("/player/{entityId}/", PlayerAPIGetID)
	s.HandleFunc("/player/", PlayerAPIGetQuery)
	s.HandleFunc("/players/", PlayersAPIGet)
	s.HandleFunc("/shows/", ShowsAPIGet)
	s.HandleFunc("/leaderboard_entries/", LeardboardEntriesAPIGet)

	// HTML Pages
	h := r.PathPrefix("/").Subrouter()
	h.HandleFunc("/", Home)

	// Facebook Callback
	h.HandleFunc("/fb_login/", FacebookLogin)

	return r
}


func init() {
	http.Handle("/", CreateHandler())
}
