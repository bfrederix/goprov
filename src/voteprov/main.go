package voteprov

import (
    "net/http"
	"github.com/gorilla/mux"
)


func CreateHandler() *mux.Router {
	r := mux.NewRouter()
	// Versioned API
	s := r.PathPrefix("/api/v1").Subrouter()
	s.HandleFunc("/player/{entityId}/", PlayerAPIGetID)
	s.HandleFunc("/player/", PlayerAPIGetQuery)
	s.HandleFunc("/players/", PlayersAPIGet)
	s.HandleFunc("/show/{entityId}/", ShowAPIGetID)
	s.HandleFunc("/shows/", ShowsAPIGet)
	s.HandleFunc("/leaderboard_entries/", LeaderboardEntriesAPIGet)
	s.HandleFunc("/suggestions/", SuggestionsAPIGet)
	s.HandleFunc("/vote_type/{entityId}/", VoteTypeAPIGetID)
	s.HandleFunc("/user_profiles/", UserProfilesAPIGet)
	s.HandleFunc("/leaderboards/user/{userId:[0-9]+}/", UserLeaderboardStatsGet)
	s.HandleFunc("/medal/{entityId}/", MedalAPIGetID)

	// User Pages
	h := r.PathPrefix("/").Subrouter()
	h.HandleFunc("/", HomePage)
	h.HandleFunc("/leaderboards/show/{show_id:[0-9]+}/", LeaderboardsPage)
	h.HandleFunc("/leaderboards/{start_date:[0-9]+}/{end_date:[0-9]+}/", LeaderboardsPage)
	h.HandleFunc("/leaderboards/", LeaderboardsPage)
	h.HandleFunc("/recap/{show_id:[0-9]+}/", ShowRecapPage)
	h.HandleFunc("/recap/", ShowRecapPage)
	h.HandleFunc("/user/{user_id:[0-9]+}/", UserAccountPage)
	h.HandleFunc("/medals/", MedalsPage)

	// Admin Pages
	a := r.PathPrefix("/admin").Subrouter()
	a.HandleFunc("/show_display/{show_id:[0-9]+}/", ShowDisplayPage)
	a.HandleFunc("/show_controller/{show_id:[0-9]+}/", ShowControllerPage)
	a.HandleFunc("/pre_show/", InstructionPage)
	a.HandleFunc("/create_show/", CreateShowPage)
	a.HandleFunc("/vote_types/", VoteTypesPage)
	a.HandleFunc("/suggestion_pools/", SuggestionPoolPage)
	a.HandleFunc("/create_medals/", CreateMedalsPage)
	a.HandleFunc("/add_players/", AddPlayersPage)
	a.HandleFunc("/delete_tools/", DeleteToolsPage)
	a.HandleFunc("/js_test/", JSTestPage)
	a.HandleFunc("/export_emails/", ExportEmailsPage)


	// Facebook Callback
	h.HandleFunc("/fb_login/", FacebookLogin)

	// Logout
	h.HandleFunc("/logout/", LogoutHandler)

	return r
}


func init() {
	http.Handle("/", CreateHandler())
}
