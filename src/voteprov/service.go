package voteprov


import (
	"time"
	"strings"
	"strconv"
	"github.com/gorilla/mux"
	"appengine"
    "appengine/datastore"
	"net/http"
	//"errors"
	"log"
)


func GetEntityKeyByIDs(c appengine.Context, modelType string, entityIdString string) (*datastore.Key) {
	// Entity ID string to int
	entityId, err := strconv.ParseInt(entityIdString, 0, 64)
	log.Println("GetModelEntity: ", entityId)
	if err != nil {
		// Couldn't find the entity by ID
		// Decode the key
    	entityKey, keyErr := datastore.DecodeKey(entityIdString)
    	if keyErr != nil {
    		// Couldn't decode the key
        	log.Fatal(keyErr)
    	}
    	return entityKey
    }
	// Get the key based on the entity ID
	entityKey := datastore.NewKey(c, modelType, "", entityId, nil)

	return entityKey
}


func GetEntityKeyByURLIDs(rw http.ResponseWriter, r *http.Request, modelType string) (appengine.Context, *datastore.Key) {
	// Get URL path variables
	vars := mux.Vars(r)
	entityIdString := vars["entityId"]
	c := appengine.NewContext(r)
	entityKey := GetEntityKeyByIDs(c, modelType, entityIdString)

	return c, entityKey
}


func GetModelEntities(rw http.ResponseWriter, r *http.Request, modelType string, limit int, params map[string]interface{}) (appengine.Context, *datastore.Query) {
	c := appengine.NewContext(r)
	q := datastore.NewQuery(modelType)
	log.Println("Query Params: ", params)
	// If empty query parameters, return the full query results
	if params == nil {
		return c, q
	}
	// If a name was specified
	if name, ok := params["name"]; ok {
		q = q.Filter("name =", name)
	}
	// If a session id was specified
	if sessionID, ok := params["current_session"]; ok {
		q = q.Filter("current_session =", sessionID)
	}
	// If a show id/key was specified
	if showID, ok := params["show"]; ok {
		showKey := GetEntityKeyByIDs(c, "Show", showID.(string))
		q = q.Filter("show =", showKey)
	}
	// If archived was specified
	if params["archived"] != nil {
		archived, err := strconv.ParseBool(params["archived"].(string))
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}
		q = q.Filter("archived =", archived)
	}
	// If ordering on created date was specified
	if params["order_by_created"] != nil {
		orderByCreated, err := strconv.ParseBool(params["order_by_created"].(string))
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}
		if orderByCreated == true {
			q = q.Order("created")
		}
	}
	// If we want to return only one item
	if limit != 0 {
		q = q.Limit(limit)
	}

	return c, q
}


func WebQueryEntities(rw http.ResponseWriter, r *http.Request, modelType string, limit int) (appengine.Context, *datastore.Query) {
	queryParams := r.URL.Query()
	params := make(map[string]interface{}, len(queryParams))
	if queryParams != nil {
		for k, v := range queryParams {
			params[k] = v[0]
		}
	}
	c, q := GetModelEntities(rw, r, modelType, limit, params)
	return c, q
}


// Need to add a new function that can create the query
// With just context and queryParams so that we can use it for setting properties

///////////////////////// Single Item Get /////////////////////////////////


func GetPlayer(rw http.ResponseWriter, r *http.Request, hasID bool, params map[string]interface{}) (*datastore.Key, Player) {
	if hasID == true {
		var player Player
		c, playerKey := GetEntityKeyByURLIDs(rw, r, "Player")

		// Try to load the data into the Player struct model
		if err := datastore.Get(c, playerKey, &player); err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}
		// Make sure the image path is set
		player.SetProperties()
		return playerKey, player
	} else {
		var players []Player
		var c appengine.Context
		var q *datastore.Query
		// If parameters were specified
		if params != nil {
			c, q = GetModelEntities(rw, r, "Player", 1, params)
		} else {
			// Otherwise use query params from url
			c, q = WebQueryEntities(rw, r, "Player", 1)
		}
		keys, err := q.GetAll(c, &players)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}
		// If nothing was found
		if players == nil {
			return &datastore.Key{}, Player{}
		}
		// Set the non-model fields
		players[0].SetProperties()
		return keys[0], players[0]
	}
}


func GetUserProfile(rw http.ResponseWriter, r *http.Request, hasID bool, params map[string]interface{}) (*datastore.Key, UserProfile) {
	if hasID == true {
		var userProfile UserProfile
		c, userProfileKey := GetEntityKeyByURLIDs(rw, r, "UserProfile")

		// Try to load the data into the UserProfile struct model
		if err := datastore.Get(c, userProfileKey, &userProfile); err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}
		// Make sure the image path is set
		//userProfile.SetProperties()
		return userProfileKey, userProfile
	} else {
		var userProfiles []UserProfile
		var c appengine.Context
		var q *datastore.Query
		// If parameters were specified
		if params != nil {
			c, q = GetModelEntities(rw, r, "UserProfile", 1, params)
		} else {
			// Otherwise use query params from url
			c, q = WebQueryEntities(rw, r, "UserProfile", 1)
		}
		keys, err := q.GetAll(c, &userProfiles)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}
		// If nothing was found
		if userProfiles == nil {
			return &datastore.Key{}, UserProfile{}
		}
		// Set the non-model fields
		//userProfiles[0].SetProperties()
		return keys[0], userProfiles[0]
	}
}


///////////////////////// Multiple Item Queries ////////////////////////////


func GetPlayers(rw http.ResponseWriter, r *http.Request, params map[string]interface{}) ([]*datastore.Key, []Player) {
	var c appengine.Context
	var q *datastore.Query
	// If parameters were specified
	if params != nil {
		c, q = GetModelEntities(rw, r, "Player", 0, params)
	} else {
		// Otherwise use query params from url
		c, q = WebQueryEntities(rw, r, "Player", 1)
	}
	var players []Player
	keys, err := q.GetAll(c, &players)
	if err != nil {
        http.Error(rw, err.Error(), http.StatusInternalServerError)
    }
	// Set the non-model fields
	for i := range players {
	    player := &players[i]
        player.SetProperties()
    }
	return keys, players
}


func GetShows(rw http.ResponseWriter, r *http.Request, params map[string]interface{}) ([]*datastore.Key, []Show) {
	var c appengine.Context
	var q *datastore.Query
	// If parameters were specified
	if params != nil {
		c, q = GetModelEntities(rw, r, "Show", 0, params)
	} else {
		// Otherwise use query params from url
		c, q = WebQueryEntities(rw, r, "Show", 1)
	}
	var shows []Show
	keys, err := q.GetAll(c, &shows)
	if err != nil {
        http.Error(rw, err.Error(), http.StatusInternalServerError)
    }

	return keys, shows
}


func GetLeaderboardEntries(rw http.ResponseWriter, r *http.Request, params map[string]interface{}) ([]*datastore.Key, []LeaderboardEntry) {
	var c appengine.Context
	var q *datastore.Query
	// If parameters were specified
	if params != nil {
		c, q = GetModelEntities(rw, r, "LeaderboardEntry", 0, params)
	} else {
		// Otherwise use query params from url
		c, q = WebQueryEntities(rw, r, "LeaderboardEntry", 1)
	}
	var leaderboardEntries []LeaderboardEntry
	keys, err := q.GetAll(c, &leaderboardEntries)
	if err != nil {
        http.Error(rw, err.Error(), http.StatusInternalServerError)
    }

	return keys, leaderboardEntries
}


//////////////////////////// Updating Entities /////////////////////////////////////////


func UpdateProfileSession(c appengine.Context, upKey *datastore.Key, userProfile *UserProfile, sessionID string) error {
	userProfile.CurrentSession = sessionID
	_, err := datastore.Put(c, upKey, userProfile)
	if err != nil {
		return err
	}
	return nil
}

func CreateUserProfile(c appengine.Context, userID string, username string,
					   email string, loginType string, currentSession string,
					   facebookToken string) (*datastore.Key, *UserProfile, error) {
	if username == "" {
		s := strings.Split(email, "@")
    	username = s[0]
	}
	stripUsername := strings.Replace(username, " ", "", -1)
	stripUsername = strings.ToLower(stripUsername)
	userProfile := UserProfile{
		UserID:         userID,
		Username:       username,
		StripUsername:  stripUsername,
		Email:          email,
		LoginType:      loginType,
		CurrentSession: currentSession,
		FBAccessToken:  facebookToken,
		Created:        time.Now().UTC(),
	}
	key, err := datastore.Put(c, datastore.NewIncompleteKey(c, "UserProfile", nil), &userProfile)
	return key, &userProfile, err
}
