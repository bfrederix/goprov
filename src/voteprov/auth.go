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
	//"appengine/datastore"
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
	Page       string
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
	sessionID, _ := session.Values["id"].(string)
	dc.SessionID = sessionID
	// Try to get user profile by session
	sessionParams := map[string]interface{}{"current_session": sessionID}
	_, userProfile := GetUserProfile(rw, r, false, sessionParams)
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
		dc.UserID = userProfile.UserID
		return
	}
	log.Println("Session data: ", sessionID)
	dc.AuthURL, _ = user.LogoutURL(c, r.URL.Path)
	dc.AuthAction = "Logout"
	dc.IsAdmin = u.Admin
	userProfile = GoogleLogin(rw, r, c, u, dc)
	dc.UserID = userProfile.UserID
	log.Println("User ID: ", dc.UserID)
}


func GoogleLogin(rw http.ResponseWriter, r *http.Request, c appengine.Context, u *user.User, dc *DefaultContext) UserProfile {
	// Try to get user profile by user id
	userIDParams := map[string]interface{}{"user_id": u.ID}
	upKey, userProfile := GetUserProfile(rw, r, false, userIDParams)
	// If the user profile was found, return
	if userProfile != (UserProfile{}) {
		err := UpdateProfileSession(c, upKey, &userProfile, dc.SessionID)
		if err == nil {
			return userProfile
		}
	} else {
		// Try to get user profile by user id
		userEmailParams := map[string]interface{}{"email": u.Email}
		upKey, userProfile := GetUserProfile(rw, r, false, userEmailParams)
		// If the user profile was found, return
		if userProfile != (UserProfile{}) {
			err := UpdateProfileSession(c, upKey, &userProfile, dc.SessionID)
			if err == nil {
				return userProfile
			}
		} else {
			_, userProfile, err := CreateUserProfile(c,
													  u.ID,
													  "",
													  u.Email,
													  "google",
													  dc.SessionID,
													  "")
			if err != nil {
				http.Error(rw, err.Error(), http.StatusInternalServerError)
				return *userProfile
			}
		}
	}

	return UserProfile{}
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
