package voteprov

import (
    "net/http"
	"html/template"
	//"log"
)


var templates = template.Must(template.ParseGlob("templates/*"))


func HomePage(rw http.ResponseWriter, r *http.Request) {
	dc := DefaultContext{}
	AuthContext(rw, r, &dc)
	dc.Page = "home"
	//log.Println("Session data: ", session.Values)
    err := templates.ExecuteTemplate(rw, "home", &dc)
	if err != nil {
        http.Error(rw, err.Error(), http.StatusInternalServerError)
    }
}


func LeaderboardsPage(rw http.ResponseWriter, r *http.Request) {
	dc := DefaultContext{}
	AuthContext(rw, r, &dc)
	dc.Page = "leaderboards"
    err := templates.ExecuteTemplate(rw, "leaderboards", &dc)
	if err != nil {
        http.Error(rw, err.Error(), http.StatusInternalServerError)
    }
}
