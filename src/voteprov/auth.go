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


func GetUserBySession(rw http.ResponseWriter, r *http.Request, sessionID string) *UserProfile {
	// Try to get profile by session
	sessionParams := map[string]interface{}{"current_session": sessionID}
	up := GetUserProfile(rw, r, false, sessionParams)

	return up
}


func AuthContext(rw http.ResponseWriter, r *http.Request, dc *DefaultContext) {
	c := appengine.NewContext(r)
	dc.ImagePath = "/static/img/"
	dc.CSSPath = "/static/css/"
	dc.JSPath = "/static/js/"
	dc.Context = c
	// Get the user
	u := user.Current(c)
	// Get the session id
	session := SessionHandler(rw, r)
	sessionID, _ = session.Values["id"].(string)
	dc.SessionID = sessionID
	// Get the user profile by session
	userProfile := GetUserBySession(rw, r, sessionID)
	// If the user doesn't exists, return
	if u == nil {
		dc.AuthURL, _ = user.LoginURL(c, r.URL.Path)
		dc.AuthAction = "Login"
		dc.UserID = ""
		dc.IsAdmin = false
		dc.Username = ""
		return
	}
	// If the user profile was found, return
	if userProfile != (UserProfile{}) {
		dc.AuthURL, _ = user.LogoutURL(c, r.URL.Path)
		dc.AuthAction = "Logout"
		dc.IsAdmin = u.Admin
		return
	}
	log.Println("User data: ", sessionID)
	dc.AuthURL, _ = user.LogoutURL(c, r.URL.Path)
	dc.AuthAction = "Logout"
	dc.IsAdmin = u.Admin
	log.Println("User data 2: ", up)
}


func GoogleLogin(*user.User, *DefaultContext) {
	return
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
		session.Values["id"] = strconv.FormatInt(rando.Int63(), 10)
    	// Save it.
    	session.Save(r, rw)
	}
	return session
}


func Authenticate(r *http.Request, u user.User) (string, string){
	return "", ""
}
