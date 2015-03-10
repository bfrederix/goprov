package api

import (
	"github.com/gorilla/mux"
    "fmt"
    "net/http"
	"log"
	"text/template"
	"encoding/json"
	"appengine"
    "appengine/user"
)


/////////////////////////////// API Endpoints ///////////////////////////


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


////////////////////// HTML Returning //////////////////////////


type DefaultContext struct {
	ImagePath  string
	CSSPath    string
	JSPath     string
	AuthURL    string
	AuthAction string
	Context    appengine.Context
	User       *user.User
}


var templates = template.Must(template.ParseGlob("templates/*"))


func SetDefaultContext(r *http.Request, dc *DefaultContext) {
	c := appengine.NewContext(r)
    u := user.Current(c)
	dc.ImagePath = "/static/img/"
	dc.CSSPath = "/static/css/"
	dc.JSPath = "/static/js/"
	dc.Context = c
	dc.User = u
	if u == nil {
		dc.AuthURL, _ = user.LoginURL(c, r.URL.Path)
		dc.AuthAction = "Login"
	} else {
		dc.AuthURL, _ = user.LogoutURL(c, r.URL.Path)
		dc.AuthAction = "Logout"
	}
}


func Welcome(rw http.ResponseWriter, r *http.Request) {
    rw.Header().Set("Content-type", "text/html; charset=utf-8")
    dc := DefaultContext{}
	SetDefaultContext(r, &dc)
    if dc.User == nil {
        fmt.Fprintf(rw, `<a href="%s">Sign in or register</a>`, dc.AuthURL)
        return
    }
    fmt.Fprintf(rw, `Welcome, %s! (<a href="%s">sign out</a>)`, dc.User, dc.AuthURL)
}


func Home(rw http.ResponseWriter, r *http.Request) {
	dc := DefaultContext{}
	SetDefaultContext(r, &dc)
    err := templates.ExecuteTemplate(rw, "home", &dc) // nil arg is context
	if err != nil {
        http.Error(rw, err.Error(), http.StatusInternalServerError)
    }
}


func CreateHandler() *mux.Router {
	r := mux.NewRouter()
	// Versioned API
	s := r.PathPrefix("/v1").Subrouter()
	s.HandleFunc("/player/{entityId}/", PlayerAPIGetID)
	s.HandleFunc("/player/", PlayerAPIGetQuery)
	s.HandleFunc("/players/", PlayersAPIGet)
	s.HandleFunc("/shows/", ShowsAPIGet)
	s.HandleFunc("/leaderboard_entries/", LeardboardEntriesAPIGet)

	//HTML Pages
	h := r.PathPrefix("/").Subrouter()
	h.HandleFunc("/", Home)
	h.HandleFunc("/welcome/", Welcome)

	return r
}


func init() {
	http.Handle("/", CreateHandler())
}
