package main

import (
    "net/http"
    "log"
    "github.com/gorilla/mux"
    "voteprov/api"
)


func CreateHandler() *mux.Router {
	r := mux.NewRouter()
	// Versioned API
	s := r.PathPrefix("/v1").Subrouter()
	s.HandleFunc("/test/", api.Test)
	s.HandleFunc("/player/{Id}/", api.PlayerAPIGetID)
	/*
	s.HandleFunc("/player/", api.PlayerAPIGetQuery)
	s.HandleFunc("/players/", PlayersAPIGet)
	s.HandleFunc("/show/{entityId}/", ShowAPIGetID)
	s.HandleFunc("/shows/", ShowsAPIGet)
	s.HandleFunc("/leaderboard_entries/", LeaderboardEntriesAPIGet)
	s.HandleFunc("/suggestions/", SuggestionsAPIGet)
	s.HandleFunc("/vote_type/{entityId}/", VoteTypeAPIGetID)
	s.HandleFunc("/user_profiles/", UserProfilesAPIGet)
	s.HandleFunc("/leaderboards/user/{userId:[0-9]+}/", UserLeaderboardStatsGet)
	s.HandleFunc("/medal/{entityId}/", MedalAPIGetID)

	// Facebook Callback
	h.HandleFunc("/fb_login/", FacebookLogin)

	// Logout
	h.HandleFunc("/logout/", LogoutHandler)
	*/

	return r
}


func main() {
	http.Handle("/", CreateHandler())
	log.Fatal(http.ListenAndServe(":4000", nil))
}
