package voteprov

import (
	"log"
	"time"
    "net/http"
	"strconv"
	"encoding/json"
	"math/rand"
	"appengine"
	"appengine/user"
	"github.com/gorilla/sessions"
)

// Setup memcached session
var cookieStore = sessions.NewCookieStore([]byte("22737468r8fs9"))

// Facebook Login Endpoint
func FacebookLogin(rw http.ResponseWriter, r *http.Request) {
	// Encode the player into JSON output
	json.NewEncoder(rw).Encode(true)
}


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
	SessionID  string
}


func AuthContext(rw http.ResponseWriter, r *http.Request, dc *DefaultContext) {
	c := appengine.NewContext(r)
    u := user.Current(c)
	session := SessionHandler(rw, r)
	dc.SessionID, _ = session.Values["id"].(string)
	dc.ImagePath = "/static/img/"
	dc.CSSPath = "/static/css/"
	dc.JSPath = "/static/js/"
	dc.Context = c
	log.Println("User data: ", session.Values["id"])
	// If the user doesn't exists
	if u == nil {
		dc.AuthURL, _ = user.LoginURL(c, r.URL.Path)
		dc.AuthAction = "Login"
		dc.UserID = ""
		dc.IsAdmin = false
		dc.Username = ""
		//dc.SessionID = session.Values["id"]
	} else {
		dc.AuthURL, _ = user.LogoutURL(c, r.URL.Path)
		dc.AuthAction = "Logout"
		dc.IsAdmin = u.Admin
		// Try to get profile by session
		sessionParams := map[string]interface{}{"current_session": dc.SessionID}
		up := GetUserProfile(rw, r, false, sessionParams)
		log.Println("User data 2: ", up)
	}
}


func SessionHandler(rw http.ResponseWriter, r *http.Request) *sessions.Session {
    // Get a session. We're ignoring the error resulted from decoding an
    // existing session: Get() always returns a session, even if empty.
    session, err := cookieStore.Get(r, "session")
	if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}
	if session.Values["id"] == nil {
		rando := rand.New(rand.NewSource(time.Now().UnixNano()))
		// Create the session id in the cookie if it doesn't exist
		session.Values["id"] = strconv.FormatInt(rando.Int63(), 64)
    	// Save it.
    	session.Save(r, rw)
	}
	return session
}


func Authenticate(r *http.Request, u user.User) (string, string){
	return "", ""
}
