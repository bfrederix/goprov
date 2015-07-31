package api

import (
    "net/http"
	"log"
	"strconv"
	"encoding/json"
	"voteprov/data"
	"github.com/gorilla/mux"
)

func Test(rw http.ResponseWriter, r *http.Request) {
	db := data.DBConnection()
	log.Println("Test: ", db)
	json.NewEncoder(rw).Encode("Hello There!")
}


func PlayerAPIGetID(rw http.ResponseWriter, r *http.Request) {
	db := data.DBConnection()
	// Get URL path variables
	vars := mux.Vars(r)
	idString := vars["Id"]
	// Entity ID string to int
	id, _ := strconv.ParseInt(idString, 0, 64)
	player, _ := data.GetPlayer(db, id)
	//log.Println("PlayerAPIGet: ", id)
	// Encode the player into JSON output
	json.NewEncoder(rw).Encode(player)
    //fmt.Fprint(rw, "Welcome! ", player.Name)
}


/*
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


func ShowAPIGetID(rw http.ResponseWriter, r *http.Request) {
	_, show := GetShow(r, true, nil)
	json.NewEncoder(rw).Encode(show)
}


func VoteTypeAPIGetID(rw http.ResponseWriter, r *http.Request) {
	_, voteType := GetVoteType(r, true, nil)
	json.NewEncoder(rw).Encode(voteType)
}


func ShowsAPIGet(rw http.ResponseWriter, r *http.Request) {
	_, shows := GetShows(r, nil)
	json.NewEncoder(rw).Encode(shows)
}


func LeaderboardEntriesAPIGet(rw http.ResponseWriter, r *http.Request) {
	//user_id + order_by_show_date available
	_, leaderboardEntries := GetLeaderboardEntries(r, nil)
	json.NewEncoder(rw).Encode(leaderboardEntries)
}


func SuggestionsAPIGet(rw http.ResponseWriter, r *http.Request) {
	_, userSuggestions := GetSuggestions(r, nil)
	json.NewEncoder(rw).Encode(userSuggestions)
}


func UserProfilesAPIGet(rw http.ResponseWriter, r *http.Request) {
	_, userProfiles := GetUserProfiles(r, false, nil)
	json.NewEncoder(rw).Encode(userProfiles)
}


func MedalAPIGetID(rw http.ResponseWriter, r *http.Request) {
	_, player := GetMedal(r, true, nil)
	json.NewEncoder(rw).Encode(player)
}


func UserLeaderboardStatsGet(rw http.ResponseWriter, r *http.Request) {
	// Get user id path variable
	vars := mux.Vars(r)
	userId := vars["userId"]
	// Get the leaderboard stats for the user
	leaderboardStats := GetLeaderboardStats(r, userId, nil, nil)
	json.NewEncoder(rw).Encode(leaderboardStats)
}
*/