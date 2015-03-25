package voteprov

import (
    "net/http"
	"log"
	"encoding/json"
	"github.com/gorilla/mux"
)


func PlayerAPIGetID(rw http.ResponseWriter, r *http.Request) {
	// Get the player entity
	_, player := GetPlayer(r, true, nil)
	log.Println("PlayerAPIGet: ", "ID!")
	// Encode the player into JSON output
	json.NewEncoder(rw).Encode(player)
    //fmt.Fprint(rw, "Welcome! ", player.Name)
}


func PlayerAPIGetQuery(rw http.ResponseWriter, r *http.Request) {
	// Get the player entity
	_, player := GetPlayer(r, false, nil)
	// Encode the player into JSON output
	json.NewEncoder(rw).Encode(player)
}


func PlayersAPIGet(rw http.ResponseWriter, r *http.Request) {
	// Get the player entity
	_, players := GetPlayers(r, nil)
	log.Println("Logging: ", "Hello!")
	// Encode the player into JSON output
	json.NewEncoder(rw).Encode(players)
    //fmt.Fprint(rw, "Welcome! ", player.Name)
}


func ShowsAPIGet(rw http.ResponseWriter, r *http.Request) {
	_, shows := GetShows(r, nil)
	json.NewEncoder(rw).Encode(shows)
}


func LeardboardEntriesAPIGet(rw http.ResponseWriter, r *http.Request) {
	_, leaderboardEntries := GetLeaderboardEntries(r, nil)
	json.NewEncoder(rw).Encode(leaderboardEntries)
}


type UserDataStruct struct {
	ShowEntries      []LeaderboardEntry
	UserSuggestions  []Suggestion
	LeaderboardStats UserTotals
	UserProfile      UserProfile
}


func UserDataAPI(rw http.ResponseWriter, r *http.Request) {
	// Get user id path variable
	vars := mux.Vars(r)
	userId := vars["userId"]

	// Get the show leaderboard entries by user id
	showLeaderboardParams := map[string]interface{}{
		"user_id": userId,
		"order_by_show_date": "True"}
	_, showLeaderboardEntries := GetLeaderboardEntries(r, showLeaderboardParams)

	// Get the suggestions by the user id
	suggestionParams := map[string]interface{}{"user_id": userId}
	_, userSuggestions := GetSuggestions(r, suggestionParams)

	// Get the leaderboard stats for the user
	leaderboardStats := GetLeaderboardStats(r, userId, nil, nil)

	// Get the user profile by user id
	userProfileParams := map[string]interface{}{"user_id": userId}
	_, userProfile := GetUserProfile(r, false, userProfileParams)

	// Create the response json structure
	uds := UserDataStruct{
		ShowEntries:      showLeaderboardEntries,
		UserSuggestions:  userSuggestions,
		LeaderboardStats: leaderboardStats,
		UserProfile:      userProfile,
	}
	json.NewEncoder(rw).Encode(uds)
}
