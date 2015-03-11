package voteprov

import (
    "net/http"
	"html/template"
	"log"
	"appengine"
    "appengine/user"
)


type DefaultContext struct {
	ImagePath  string
	CSSPath    string
	JSPath     string
	AuthURL    string
	AuthAction string
	Context    appengine.Context
	UserID     string
	IsAdmin    bool
	Username   string
}


var templates = template.Must(template.ParseGlob("templates/*"))


func AuthContext(r *http.Request, dc *DefaultContext) {
	c := appengine.NewContext(r)
    u := user.Current(c)
	dc.ImagePath = "/static/img/"
	dc.CSSPath = "/static/css/"
	dc.JSPath = "/static/js/"
	dc.Context = c
	dc.IsAdmin = u.Admin

	if u == nil {
		dc.AuthURL, _ = user.LoginURL(c, r.URL.Path)
		dc.AuthAction = "Login"
	} else {
		dc.AuthURL, _ = user.LogoutURL(c, r.URL.Path)
		dc.AuthAction = "Logout"
	}
}


func Home(rw http.ResponseWriter, r *http.Request) {
	dc := DefaultContext{}
	AuthContext(r, &dc)
	MyHandler(rw, r)
	session, _ := store.Get(r, "session-name")
	log.Println("Session data: ", session.Values["example"])
    err := templates.ExecuteTemplate(rw, "home", &dc) // nil arg is context
	if err != nil {
        http.Error(rw, err.Error(), http.StatusInternalServerError)
    }
}
