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


func CreateShowPage(rw http.ResponseWriter, r *http.Request) {
	dc := DefaultContext{}
	AuthContext(rw, r, &dc)
	// Redirect if not admin
	if !dc.IsAdmin {
		AdminRedirect(rw, r)
	}
	dc.Page = "create_show"
    err := templates.ExecuteTemplate(rw, "create_show", &dc)
	if err != nil {
        http.Error(rw, err.Error(), http.StatusInternalServerError)
    }
}


func VoteTypesPage(rw http.ResponseWriter, r *http.Request) {
	dc := DefaultContext{}
	AuthContext(rw, r, &dc)
	// Redirect if not admin
	if !dc.IsAdmin {
		AdminRedirect(rw, r)
	}
	dc.Page = "vote_types"
    err := templates.ExecuteTemplate(rw, "vote_types", &dc)
	if err != nil {
        http.Error(rw, err.Error(), http.StatusInternalServerError)
    }
}


func SuggestionPoolPage(rw http.ResponseWriter, r *http.Request) {
	dc := DefaultContext{}
	AuthContext(rw, r, &dc)
	// Redirect if not admin
	if !dc.IsAdmin {
		AdminRedirect(rw, r)
	}
	dc.Page = "suggestion_pools"
    err := templates.ExecuteTemplate(rw, "suggestion_pools", &dc)
	if err != nil {
        http.Error(rw, err.Error(), http.StatusInternalServerError)
    }
}


func CreateMedalsPage(rw http.ResponseWriter, r *http.Request) {
	dc := DefaultContext{}
	AuthContext(rw, r, &dc)
	// Redirect if not admin
	if !dc.IsAdmin {
		AdminRedirect(rw, r)
	}
	dc.Page = "create_medals"
    err := templates.ExecuteTemplate(rw, "create_medals", &dc)
	if err != nil {
        http.Error(rw, err.Error(), http.StatusInternalServerError)
    }
}


func AddPlayersPage(rw http.ResponseWriter, r *http.Request) {
	dc := DefaultContext{}
	AuthContext(rw, r, &dc)
	// Redirect if not admin
	if !dc.IsAdmin {
		AdminRedirect(rw, r)
	}
	dc.Page = "add_players"
    err := templates.ExecuteTemplate(rw, "add_players", &dc)
	if err != nil {
        http.Error(rw, err.Error(), http.StatusInternalServerError)
    }
}


func DeleteToolsPage(rw http.ResponseWriter, r *http.Request) {
	dc := DefaultContext{}
	AuthContext(rw, r, &dc)
	// Redirect if not admin
	if !dc.IsAdmin {
		AdminRedirect(rw, r)
	}
	dc.Page = "delete_tools"
    err := templates.ExecuteTemplate(rw, "delete_tools", &dc)
	if err != nil {
        http.Error(rw, err.Error(), http.StatusInternalServerError)
    }
}


func JSTestPage(rw http.ResponseWriter, r *http.Request) {
	dc := DefaultContext{}
	AuthContext(rw, r, &dc)
	// Redirect if not admin
	if !dc.IsAdmin {
		AdminRedirect(rw, r)
	}
	dc.Page = "js_test"
    err := templates.ExecuteTemplate(rw, "js_test", &dc)
	if err != nil {
        http.Error(rw, err.Error(), http.StatusInternalServerError)
    }
}


func ExportEmailsPage(rw http.ResponseWriter, r *http.Request) {
	dc := DefaultContext{}
	AuthContext(rw, r, &dc)
	// Redirect if not admin
	if !dc.IsAdmin {
		AdminRedirect(rw, r)
	}
	dc.Page = "export_emails"
    err := templates.ExecuteTemplate(rw, "export_emails", &dc)
	if err != nil {
        http.Error(rw, err.Error(), http.StatusInternalServerError)
    }
}
