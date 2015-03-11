package voteprov

import (
    "net/http"
	"log"
	"encoding/json"
)

func PlayerAPIGetID(rw http.ResponseWriter, r *http.Request) {
	// Get the player entity
	player := GetPlayer(rw, r, true)
	log.Println("PlayerAPIGet: ", "ID!")
	// Encode the player into JSON output
	json.NewEncoder(rw).Encode(player)
    //fmt.Fprint(rw, "Welcome! ", player.Name)
}


func PlayerAPIGetQuery(rw http.ResponseWriter, r *http.Request) {
	// Get the player entity
	player := GetPlayer(rw, r, false)
	// Encode the player into JSON output
	json.NewEncoder(rw).Encode(player)
}


func PlayersAPIGet(rw http.ResponseWriter, r *http.Request) {
	// Get the player entity
	players := GetPlayers(rw, r)
	log.Println("Logging: ", "Hello!")
	// Encode the player into JSON output
	json.NewEncoder(rw).Encode(players)
    //fmt.Fprint(rw, "Welcome! ", player.Name)
}


func ShowsAPIGet(rw http.ResponseWriter, r *http.Request) {
	shows := GetShows(rw, r)
	json.NewEncoder(rw).Encode(shows)
}


func LeardboardEntriesAPIGet(rw http.ResponseWriter, r *http.Request) {
	leaderboardEntries := GetLeaderboardEntries(rw, r)
	json.NewEncoder(rw).Encode(leaderboardEntries)
}
