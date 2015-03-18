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


func InstructionPage(rw http.ResponseWriter, r *http.Request) {
	dc := DefaultContext{}
	AuthContext(rw, r, &dc)
	// Redirect if not admin
	if !dc.IsAdmin {
		AdminRedirect(rw, r)
	}
	dc.Page = "pre_show"
    err := templates.ExecuteTemplate(rw, "pre_show", &dc)
	if err != nil {
        http.Error(rw, err.Error(), http.StatusInternalServerError)
    }
}
