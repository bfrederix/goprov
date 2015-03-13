package voteprov


import (
	"github.com/gorilla/mux"
	"strconv"
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
		log.Println("Query Name!")
		q = q.Filter("name =", name)
	}
	// If a session id was specified
	if sessionID, ok := params["id"]; ok {
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


func GetPlayer(rw http.ResponseWriter, r *http.Request, hasID bool, params map[string]interface{}) (Player) {
	if hasID == true {
		var player Player
		c, playerKey := GetEntityKeyByURLIDs(rw, r, "Player")

		// Try to load the data into the Player struct model
		if err := datastore.Get(c, playerKey, &player); err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}
		// Make sure the image path is set
		player.SetProperties()
		return player
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

		if _, err := q.GetAll(c, &players); err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}
		// Set the non-model fields
		players[0].SetProperties()
		return players[0]
	}
}


func GetUserProfile(rw http.ResponseWriter, r *http.Request, hasID bool, params map[string]interface{}) (UserProfile) {
	if hasID == true {
		var userProfile UserProfile
		c, userProfileKey := GetEntityKeyByURLIDs(rw, r, "UserProfile")

		// Try to load the data into the UserProfile struct model
		if err := datastore.Get(c, userProfileKey, &userProfile); err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}
		// Make sure the image path is set
		//userProfile.SetProperties()
		return userProfile
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

		if _, err := q.GetAll(c, &userProfiles); err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}
		// Set the non-model fields
		//userProfiles[0].SetProperties()
		return userProfiles[0]
	}
}


///////////////////////// Multiple Item Queries ////////////////////////////


func GetPlayers(rw http.ResponseWriter, r *http.Request, params map[string]interface{}) ([]Player) {
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
	if _, err := q.GetAll(c, &players); err != nil {
        http.Error(rw, err.Error(), http.StatusInternalServerError)
    }
	// Set the non-model fields
	for i := range players {
	    player := &players[i]
        player.SetProperties()
    }
	return players
}


func GetShows(rw http.ResponseWriter, r *http.Request, params map[string]interface{}) ([]Show) {
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
	if _, err := q.GetAll(c, &shows); err != nil {
        http.Error(rw, err.Error(), http.StatusInternalServerError)
    }

	return shows
}


func GetLeaderboardEntries(rw http.ResponseWriter, r *http.Request, params map[string]interface{}) ([]LeaderboardEntry) {
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
	if _, err := q.GetAll(c, &leaderboardEntries); err != nil {
        http.Error(rw, err.Error(), http.StatusInternalServerError)
    }

	return leaderboardEntries
}
