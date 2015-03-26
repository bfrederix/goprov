package voteprov

import (
	"fmt"
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


// Handle logging out by resetting the session and forwarding to logout
func LogoutHandler(rw http.ResponseWriter, r *http.Request) {
	referer := r.Referer()
	logoutPath := fmt.Sprintf("/_ah/login?action=logout&continue=%s", referer)
	ResetSession(rw, r)
	http.Redirect(rw, r, logoutPath, 302)
}


func AdminRedirect(rw http.ResponseWriter, r *http.Request) {
	referer := r.Referer()
	logoutPath := fmt.Sprintf("/_ah/login?continue=%s", referer)
	ResetSession(rw, r)
	http.Redirect(rw, r, logoutPath, 302)
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
	_, userProfile := GetUserProfiles(r, false, sessionParams)
	log.Println("userProfile: ", userProfile)
	// If the user doesn't exists, return
	if u == nil {
		dc.AuthURL, _ = user.LoginURL(c, r.URL.Path)
		dc.AuthAction = "Login"
		dc.UserID = ""
		dc.IsAdmin = false
		dc.Username = ""
		return
	}
	log.Println("User ID: ", u.ID)
	// If the user profile was found, return
	if userProfile.UserID != "" {
		dc.AuthURL = "/logout/"
		dc.AuthAction = "Logout"
		dc.IsAdmin = u.Admin
		dc.UserID = userProfile.UserID
		log.Println("userProfile FOUND: ", userProfile)
		return
	}
	dc.AuthURL = "/logout/"
	dc.AuthAction = "Logout"
	dc.IsAdmin = u.Admin
	userProfile = GoogleLogin(r, c, u, dc)
	log.Println("Google Profile: ", userProfile)
	dc.UserID = userProfile.UserID
}


func GoogleLogin(r *http.Request, c appengine.Context, u *user.User, dc *DefaultContext) UserProfile {
	// Try to get user profile by user id
	userIDParams := map[string]interface{}{"user_id": u.ID}
	upKey, userProfile := GetUserProfiles(r, false, userIDParams)
	log.Println("ID Based: ", userProfile)
	// If the user profile was found, return
	if userProfile.UserID != "" {
		err := UpdateProfileSession(c, upKey, &userProfile, dc.SessionID)
		if err == nil {
			return userProfile
		}
	} else {
		// Try to get user profile by user id
		userEmailParams := map[string]interface{}{"email": u.Email}
		upKey, userProfile := GetUserProfiles(r, false, userEmailParams)
		log.Println("Email Based: ", userProfile)
		// If the user profile was found, return
		if userProfile.UserID != "" {
			err := UpdateProfileSession(c, upKey, &userProfile, dc.SessionID)
			if err == nil {
				return userProfile
			}
		} else {
			log.Println("Create New Profile")
			// Create the new user profile
			_, userProfile, err := CreateUserProfile(c,
													  u.ID,
													  "",
													  u.Email,
													  "google",
													  dc.SessionID,
													  "")
			if err != nil {
				panic(err.Error())
			}
			return *userProfile
		}
	}

	return UserProfile{}
}


func SessionHandler(rw http.ResponseWriter, r *http.Request) *sessions.Session {
    // Get a session. We're ignoring the error resulted from decoding an
    // existing session: Get() always returns a session, even if empty.
    session := GetSession(r)
	if session.Values["id"] == nil {
		NewSessionID(rw, r, session)
	}
	return session
}


func GetSession(r *http.Request) *sessions.Session {
	session, err := cookieStore.Get(r, "session")
	if err != nil {
		panic(err.Error())
	}
	return session
}


func NewSessionID(rw http.ResponseWriter, r *http.Request, session *sessions.Session) {
	rando := rand.New(rand.NewSource(time.Now().UnixNano()))
	// Create the session id in the cookie if it doesn't exist
	session.Values["id"] = strconv.FormatInt(rando.Int63(), 10)
	session.Save(r, rw)
}


func ResetSession(rw http.ResponseWriter, r *http.Request) {
	session := GetSession(r)
	NewSessionID(rw, r, session)
}
