package voteprov

import (
    "net/http"
	"html/template"
	//"log"
)


type BaseContext struct {
	DC        DefaultContext
}


var templates = template.Must(template.ParseGlob("templates/*"))


func HomePage(rw http.ResponseWriter, r *http.Request) {
	bc := BaseContext{}
	AuthContext(rw, r, &bc.DC)
	bc.DC.Page = "home"
	//log.Println("Session data: ", session.Values)
    err := templates.ExecuteTemplate(rw, "home", &bc)
	if err != nil {
        http.Error(rw, err.Error(), http.StatusInternalServerError)
    }
}


func LeaderboardsPage(rw http.ResponseWriter, r *http.Request) {
	bc := BaseContext{}
	AuthContext(rw, r, &bc.DC)
	bc.DC.Page = "leaderboards"
    err := templates.ExecuteTemplate(rw, "leaderboards", &bc)
	if err != nil {
        http.Error(rw, err.Error(), http.StatusInternalServerError)
    }
}


func ShowRecapPage(rw http.ResponseWriter, r *http.Request) {
	bc := BaseContext{}
	AuthContext(rw, r, &bc.DC)
	bc.DC.Page = "recap"
    err := templates.ExecuteTemplate(rw, "recap", &bc)
	if err != nil {
        http.Error(rw, err.Error(), http.StatusInternalServerError)
    }
}


func UserAccountPage(rw http.ResponseWriter, r *http.Request) {
	bc := BaseContext{}
	AuthContext(rw, r, &bc.DC)
	bc.DC.Page = "user"
    err := templates.ExecuteTemplate(rw, "user", &bc)
	if err != nil {
        http.Error(rw, err.Error(), http.StatusInternalServerError)
    }
}


func MedalsPage(rw http.ResponseWriter, r *http.Request) {
	bc := BaseContext{}
	AuthContext(rw, r, &bc.DC)
	bc.DC.Page = "medals"
    err := templates.ExecuteTemplate(rw, "medals", &bc)
	if err != nil {
        http.Error(rw, err.Error(), http.StatusInternalServerError)
    }
}


/////////////////////////// Admin Pages ///////////////////////////


func InstructionPage(rw http.ResponseWriter, r *http.Request) {
	bc := BaseContext{}
	AuthContext(rw, r, &bc.DC)
	// Redirect if not admin
	if !bc.DC.IsAdmin {
		AdminRedirect(rw, r)
	}
	bc.DC.Page = "pre_show"
    err := templates.ExecuteTemplate(rw, "pre_show", &bc)
	if err != nil {
        http.Error(rw, err.Error(), http.StatusInternalServerError)
    }
}


func CreateShowPage(rw http.ResponseWriter, r *http.Request) {
	bc := BaseContext{}
	AuthContext(rw, r, &bc.DC)
	// Redirect if not admin
	if !bc.DC.IsAdmin {
		AdminRedirect(rw, r)
	}
	bc.DC.Page = "create_show"
    err := templates.ExecuteTemplate(rw, "create_show", &bc)
	if err != nil {
        http.Error(rw, err.Error(), http.StatusInternalServerError)
    }
}


func VoteTypesPage(rw http.ResponseWriter, r *http.Request) {
	bc := BaseContext{}
	AuthContext(rw, r, &bc.DC)
	// Redirect if not admin
	if !bc.DC.IsAdmin {
		AdminRedirect(rw, r)
	}
	bc.DC.Page = "vote_types"
    err := templates.ExecuteTemplate(rw, "vote_types", &bc)
	if err != nil {
        http.Error(rw, err.Error(), http.StatusInternalServerError)
    }
}


func SuggestionPoolPage(rw http.ResponseWriter, r *http.Request) {
	bc := BaseContext{}
	AuthContext(rw, r, &bc.DC)
	// Redirect if not admin
	if !bc.DC.IsAdmin {
		AdminRedirect(rw, r)
	}
	bc.DC.Page = "suggestion_pools"
    err := templates.ExecuteTemplate(rw, "suggestion_pools", &bc)
	if err != nil {
        http.Error(rw, err.Error(), http.StatusInternalServerError)
    }
}


func CreateMedalsPage(rw http.ResponseWriter, r *http.Request) {
	bc := BaseContext{}
	AuthContext(rw, r, &bc.DC)
	// Redirect if not admin
	if !bc.DC.IsAdmin {
		AdminRedirect(rw, r)
	}
	bc.DC.Page = "create_medals"
    err := templates.ExecuteTemplate(rw, "create_medals", &bc)
	if err != nil {
        http.Error(rw, err.Error(), http.StatusInternalServerError)
    }
}


func AddPlayersPage(rw http.ResponseWriter, r *http.Request) {
	bc := BaseContext{}
	AuthContext(rw, r, &bc.DC)
	// Redirect if not admin
	if !bc.DC.IsAdmin {
		AdminRedirect(rw, r)
	}
	bc.DC.Page = "add_players"
    err := templates.ExecuteTemplate(rw, "add_players", &bc)
	if err != nil {
        http.Error(rw, err.Error(), http.StatusInternalServerError)
    }
}


func DeleteToolsPage(rw http.ResponseWriter, r *http.Request) {
	bc := BaseContext{}
	AuthContext(rw, r, &bc.DC)
	// Redirect if not admin
	if !bc.DC.IsAdmin {
		AdminRedirect(rw, r)
	}
	bc.DC.Page = "delete_tools"
    err := templates.ExecuteTemplate(rw, "delete_tools", &bc)
	if err != nil {
        http.Error(rw, err.Error(), http.StatusInternalServerError)
    }
}


func JSTestPage(rw http.ResponseWriter, r *http.Request) {
	bc := BaseContext{}
	AuthContext(rw, r, &bc.DC)
	// Redirect if not admin
	if !bc.DC.IsAdmin {
		AdminRedirect(rw, r)
	}
	bc.DC.Page = "js_test"
    err := templates.ExecuteTemplate(rw, "js_test", &bc)
	if err != nil {
        http.Error(rw, err.Error(), http.StatusInternalServerError)
    }
}


func ExportEmailsPage(rw http.ResponseWriter, r *http.Request) {
	bc := BaseContext{}
	AuthContext(rw, r, &bc.DC)
	// Redirect if not admin
	if !bc.DC.IsAdmin {
		AdminRedirect(rw, r)
	}
	bc.DC.Page = "export_emails"
    err := templates.ExecuteTemplate(rw, "export_emails", &bc)
	if err != nil {
        http.Error(rw, err.Error(), http.StatusInternalServerError)
    }
}
