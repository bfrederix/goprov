package voteprov

import (
	"time"
    "net/http"
	"encoding/json"
	"math/rand"
	"appengine/user"
	"github.com/gorilla/appengine/sessions"
)

// Setup memcached session
var memStore = sessions.NewMemcacheStore("", []byte("85f8djd8s0sx0"))
var cookieStore = sessions.NewCookieStore([]byte("22737468r8fs9"))

// Facebook Login Endpoint
func FacebookLogin(rw http.ResponseWriter, r *http.Request) {
	// Encode the player into JSON output
	json.NewEncoder(rw).Encode(true)
}

func SessionHandler(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
    // Get a session. We're ignoring the error resulted from decoding an
    // existing session: Get() always returns a session, even if empty.
    session, _ := cookieStore.Get(r, "session-id")
	if session.Values["id"] == nil {
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		// Create the session id in the cookie if it doesn't exist
		session.Values["id"] = r.Int63()
	}

    // Set some session values.
    session.Values["example"] = queryParams["name"]
    // Save it.
    session.Save(r, w)
}


func Authenticate(r *http.Request, u user.User) (string, string){
	return "", ""
}
