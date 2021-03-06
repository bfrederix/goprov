package data


/*
import (
	"fmt"
	"time"
	"strings"
	"strconv"
	"sort"
	"github.com/gorilla/mux"
	"appengine"
    "appengine/datastore"
	"net/http"
	//"errors"
	"log"
)
*/


import (
	"database/sql"
    "fmt"
    //"log"
    _ "github.com/lib/pq"
)


const (
    DB_HOST_URL = "devvoteprov.crtjwt7ubwk0.us-west-2.rds.amazonaws.com:5432"
    DB_USER     = "voteprovprod"
    DB_PASSWORD = "pr0dpr0v"
    DB_NAME     = "voteprov_prod"
)


func DBConnection() (db *sql.DB) {
    dbUrl := fmt.Sprintf("postgresql://%s:%s@%s/%s",
        DB_USER, DB_PASSWORD, DB_HOST_URL, DB_NAME)
    db, err := sql.Open("postgres", dbUrl)
    if err != nil {
        panic(err)
    }
    return db
}


func GetPlayer(db *sql.DB, id int64) (*Player, error) {
    const query = `SELECT id,name,photo_filename,star,created from players_player where id = $1 `
    var player Player
    err := db.QueryRow(query, id).Scan(&player.Id,
                                       &player.Name,
                                       &player.PhotoFilename,
                                       &player.Star,
                                       &player.Created)
    player.SetProperties()
    return &player, err
}


/*
func stringInSlice(a string, list []string) bool {
    for _, b := range list {
        if b == a {
            return true
        }
    }
    return false
}


func intPos(slice []int64, value int64) int {
	for p, v := range slice {
		log.Println("pv: ", p, v, value)
        if (v == value) {
            return p
        }
    }
    return -1
}


func GetMostRecentShow(r *http.Request,) (*datastore.Key, Show) {
	showParams := map[string]interface{}{"order_by_created": true}
	showKey, show := GetShow(r, false, showParams)
	return showKey, show
}


func GetEntityKeyByIDs(c appengine.Context, modelType string, entityIdString string) (*datastore.Key) {
	// Entity ID string to int
	entityId, err := strconv.ParseInt(entityIdString, 0, 64)
	if err != nil {
		// Couldn't find the entity by ID
		// Decode the key
    	entityKey, keyErr := datastore.DecodeKey(entityIdString)
    	if keyErr != nil {
    		// Couldn't decode the key
        	panic(keyErr)
    	}
    	return entityKey
    }
	// Get the key based on the entity ID
	entityKey := datastore.NewKey(c, modelType, "", entityId, nil)

	return entityKey
}


func GetEntityKeyByURLIDs(c appengine.Context, r *http.Request, modelType string) *datastore.Key {
	// Get URL path variables
	vars := mux.Vars(r)
	entityIdString := vars["entityId"]
	entityKey := GetEntityKeyByIDs(c, modelType, entityIdString)

	return entityKey
}


func GetModelEntities(c appengine.Context, modelType string, limit int, params map[string]interface{}) *datastore.Query {
	q := datastore.NewQuery(modelType)
	log.Println("Query Params: ", params)
	paramsAllowed := []string{
		"key",
        "name",
		"current_session",
		"user_id",
		"email",
		"show",
		"vote_type",
		"interval",
		"archived",
		"order_by_created",
		"order_by_show_date",
	}
	// Check to make sure the parameter passed is allowed
	for key, _ := range params {
		if ok := stringInSlice(key, paramsAllowed); !ok {
			failure := fmt.Sprintf("Query Parameter Not Allowed: %s", key)
			panic(failure)
		}
	}
	// If empty query parameters, return the full query results
	if params == nil {
		return q
	}
	// If a key was specified
	if key, ok := params["key"]; ok {
		q = q.Filter("__key__ =", key.(*datastore.Key))
	}
	// If a certain params were specified
	if name, ok := params["name"]; ok {
		q = q.Filter("name =", name)
	}
	if sessionID, ok := params["current_session"]; ok {
		q = q.Filter("current_session =", sessionID)
	}
	if userID, ok := params["user_id"]; ok {
		q = q.Filter("user_id =", userID)
	}
	if email, ok := params["email"]; ok {
		q = q.Filter("email =", email)
	}
	if interval, ok := params["interval"]; ok {
		q = q.Filter("interval =", interval)
	}
	// If a show id/key was specified
	if showID, ok := params["show"]; ok {
		showKey := GetEntityKeyByIDs(c, "Show", showID.(string))
		q = q.Filter("show =", showKey)
	}
	// If a vote type id/key was specified
	if voteTypeID, ok := params["vote_type"]; ok {
		voteTypeKey := GetEntityKeyByIDs(c, "VoteType", voteTypeID.(string))
		q = q.Filter("vote_type =", voteTypeKey)
	}
	if params["archived"] != nil {
		archived, err := strconv.ParseBool(params["archived"].(string))
		if err != nil {
			panic(err.Error())
		}
		q = q.Filter("archived =", archived)
	}
	// If ordering was specified
	if params["order_by_created"] != nil {
		orderByCreated, err := strconv.ParseBool(params["order_by_created"].(string))
		if err != nil {
			panic(err.Error())
		}
		if orderByCreated == true {
			q = q.Order("created")
		}
	}
	if params["order_by_show_date"] != nil {
		orderByCreated, err := strconv.ParseBool(params["order_by_show_date"].(string))
		if err != nil {
			panic(err.Error())
		}
		if orderByCreated == true {
			q = q.Order("-show_date")
		}
	}
	// If we want to return only one item
	if limit != 0 {
		q = q.Limit(limit)
	}

	return q
}


func MapQuery(r *http.Request) map[string]interface{} {
	// Map the query parameters
	queryParams := r.URL.Query()
	params := make(map[string]interface{}, len(queryParams))
	if queryParams != nil {
		// Make a map out of the parameters with flexible values
		for k, v := range queryParams {
			params[k] = v[0]
		}
	}
	return params
}


func WebQueryEntities(c appengine.Context, r *http.Request, modelType string, limit int) *datastore.Query {
	// Get the query parameters
	params := MapQuery(r)
	q := GetModelEntities(c, modelType, limit, params)
	return q
}


// Need to add a new function that can create the query
// With just context and queryParams so that we can use it for setting properties

///////////////////////// Single Item Get /////////////////////////////////


func GetPlayer(r *http.Request, hasID bool, params map[string]interface{}) (*datastore.Key, Player) {
	c := appengine.NewContext(r)
	if hasID == true {
		var player Player
		playerKey := GetEntityKeyByURLIDs(c, r, "Player")

		// Try to load the data into the Player struct model
		if err := datastore.Get(c, playerKey, &player); err != nil {
			panic(err.Error())
		}
		// Make sure the image path is set
		player.SetProperties()
		return playerKey, player
	} else {
		var players []Player
		var q *datastore.Query
		// If parameters were specified
		if params != nil {
			q = GetModelEntities(c, "Player", 1, params)
		} else {
			// Otherwise use query params from url
			q = WebQueryEntities(c, r, "Player", 1)
		}
		keys, err := q.GetAll(c, &players)
		if err != nil {
			panic(err.Error())
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


func GetVoteType(r *http.Request, hasID bool, params map[string]interface{}) (*datastore.Key, VoteType) {
	c := appengine.NewContext(r)
	if hasID == true {
		var voteType VoteType
		voteTypeKey := GetEntityKeyByURLIDs(c, r, "VoteType")

		// Try to load the data into the VoteType struct model
		if err := datastore.Get(c, voteTypeKey, &voteType); err != nil {
			panic(err.Error())
		}
		voteType.SetProperties(voteTypeKey, r, nil)
		return voteTypeKey, voteType
	} else {
		var voteTypes []VoteType
		var q *datastore.Query
		// If parameters were specified
		if params != nil {
			q = GetModelEntities(c, "VoteType", 1, params)
		} else {
			// Otherwise use query params from url
			q = WebQueryEntities(c, r, "VoteType", 1)
		}
		keys, err := q.GetAll(c, &voteTypes)
		if err != nil {
			panic(err.Error())
		}
		// If nothing was found
		if voteTypes == nil {
			return &datastore.Key{}, VoteType{}
		}
		// If the show entity was specified
		if showEntity, ok := params["showEntity"]; ok {
			voteTypes[0].SetProperties(keys[0], r, showEntity.(*Show))
		} else {
			voteTypes[0].SetProperties(keys[0], r, nil)
		}
		return keys[0], voteTypes[0]
	}
}


func GetShow(r *http.Request, hasID bool, params map[string]interface{}) (*datastore.Key, Show) {
	c := appengine.NewContext(r)
	if hasID == true {
		var show Show
		showKey := GetEntityKeyByURLIDs(c, r, "Show")

		// Try to load the data into the Show struct model
		if err := datastore.Get(c, showKey, &show); err != nil {
			panic(err.Error())
		}
		show.SetProperties(c, showKey, r)
		return showKey, show
	} else {
		var shows []Show
		var q *datastore.Query
		// If parameters were specified
		if params != nil {
			q = GetModelEntities(c, "Show", 1, params)
		} else {
			// Otherwise use query params from url
			q = WebQueryEntities(c, r, "Show", 1)
		}
		keys, err := q.GetAll(c, &shows)
		if err != nil {
			panic(err.Error())
		}
		// If nothing was found
		if shows == nil {
			return &datastore.Key{}, Show{}
		}
		shows[0].SetProperties(c, keys[0], r)
		return keys[0], shows[0]
	}
}


func GetUserProfiles(r *http.Request, hasID bool, params map[string]interface{}) (*datastore.Key, UserProfile) {
	c := appengine.NewContext(r)
	if hasID == true {
		var userProfile UserProfile
		userProfileKey := GetEntityKeyByURLIDs(c, r, "UserProfile")

		// Try to load the data into the UserProfile struct model
		if err := datastore.Get(c, userProfileKey, &userProfile); err != nil {
			panic(err.Error())
		}
		// Make sure the image path is set
		userProfile.SetProperties()
		return userProfileKey, userProfile
	} else {
		var userProfiles []UserProfile
		var q *datastore.Query
		// If parameters were specified
		if params != nil {
			q = GetModelEntities(c, "UserProfile", 1, params)
		} else {
			// Otherwise use query params from url
			q = WebQueryEntities(c, r, "UserProfile", 1)
		}
		keys, err := q.GetAll(c, &userProfiles)
		if err != nil {
			panic(err.Error())
		}
		// If nothing was found
		if userProfiles == nil {
			return &datastore.Key{}, UserProfile{}
		}
		// Set the non-model fields
		userProfiles[0].SetProperties()
		return keys[0], userProfiles[0]
	}
}


func GetShowInterval(r *http.Request, hasID bool, params map[string]interface{}) (*datastore.Key, ShowInterval) {
	c := appengine.NewContext(r)
	if hasID == true {
		var showInterval ShowInterval
		showIntervalKey := GetEntityKeyByURLIDs(c, r, "ShowInterval")

		// Try to load the data into the ShowInterval struct model
		if err := datastore.Get(c, showIntervalKey, &showInterval); err != nil {
			panic(err.Error())
		}
		return showIntervalKey, showInterval
	} else {
		var showIntervals []ShowInterval
		var q *datastore.Query
		// If parameters were specified
		if params != nil {
			q = GetModelEntities(c, "ShowInterval", 1, params)
		} else {
			// Otherwise use query params from url
			q = WebQueryEntities(c, r, "ShowInterval", 1)
		}
		keys, err := q.GetAll(c, &showIntervals)
		if err != nil {
			panic(err.Error())
		}
		// If nothing was found
		if showIntervals == nil {
			return &datastore.Key{}, ShowInterval{}
		}

		return keys[0], showIntervals[0]
	}
}


func GetVotedItem(r *http.Request, hasID bool, params map[string]interface{}) (*datastore.Key, VotedItem) {
	c := appengine.NewContext(r)
	if hasID == true {
		var votedItem VotedItem
		votedItemKey := GetEntityKeyByURLIDs(c, r, "VotedItem")

		// Try to load the data into the VotedItem struct model
		if err := datastore.Get(c, votedItemKey, &votedItem); err != nil {
			panic(err.Error())
		}
		return votedItemKey, votedItem
	} else {
		var votedItems []VotedItem
		var q *datastore.Query
		// If parameters were specified
		if params != nil {
			q = GetModelEntities(c, "VotedItem", 1, params)
		} else {
			// Otherwise use query params from url
			q = WebQueryEntities(c, r, "VotedItem", 1)
		}
		keys, err := q.GetAll(c, &votedItems)
		if err != nil {
			panic(err.Error())
		}
		// If nothing was found
		if votedItems == nil {
			return &datastore.Key{}, VotedItem{}
		}

		return keys[0], votedItems[0]
	}
}


func GetMedal(r *http.Request, hasID bool, params map[string]interface{}) (*datastore.Key, Medal) {
	c := appengine.NewContext(r)
	if hasID == true {
		var medal Medal
		medalKey := GetEntityKeyByURLIDs(c, r, "Medal")

		// Try to load the data into the Medal struct model
		if err := datastore.Get(c, medalKey, &medal); err != nil {
			panic(err.Error())
		}
		return medalKey, medal
	} else {
		var medals []Medal
		var q *datastore.Query
		// If parameters were specified
		if params != nil {
			q = GetModelEntities(c, "Medal", 1, params)
		} else {
			// Otherwise use query params from url
			q = WebQueryEntities(c, r, "Medal", 1)
		}
		keys, err := q.GetAll(c, &medals)
		if err != nil {
			panic(err.Error())
		}
		// If nothing was found
		if medals == nil {
			return &datastore.Key{}, Medal{}
		}

		return keys[0], medals[0]
	}
}

///////////////////////// Multiple Item Queries ////////////////////////////


func GetPlayers(r *http.Request, params map[string]interface{}) ([]*datastore.Key, []Player) {
	c := appengine.NewContext(r)
	var q *datastore.Query
	// If parameters were specified
	if params != nil {
		q = GetModelEntities(c, "Player", 0, params)
	} else {
		// Otherwise use query params from url
		q = WebQueryEntities(c, r, "Player", 0)
	}
	var players []Player
	keys, err := q.GetAll(c, &players)
	if err != nil {
        panic(err.Error())
    }
	// Set the non-model fields
	for i := range players {
	    player := &players[i]
        player.SetProperties()
    }
	return keys, players
}


func GetShows(r *http.Request, params map[string]interface{}) ([]*datastore.Key, []Show) {
	c := appengine.NewContext(r)
	var q *datastore.Query
	// If parameters were specified
	if params != nil {
		q = GetModelEntities(c, "Show", 0, params)
	} else {
		// Otherwise use query params from url
		q = WebQueryEntities(c, r, "Show", 0)
	}
	var shows []Show
	keys, err := q.GetAll(c, &shows)
	if err != nil {
        panic(err.Error())
    }
	// Set the non-model fields
	for i := range shows {
	    show := &shows[i]

        show.SetProperties(c, keys[i], r)
    }

	return keys, shows
}


func GetLeaderboardEntries(r *http.Request, params map[string]interface{}) ([]*datastore.Key, []LeaderboardEntry) {
	c := appengine.NewContext(r)
	var q *datastore.Query
	// If parameters were specified
	if params != nil {
		q = GetModelEntities(c, "LeaderboardEntry", 0, params)
	} else {
		// Otherwise use query params from url
		q = WebQueryEntities(c, r, "LeaderboardEntry", 0)
	}
	var leaderboardEntries []LeaderboardEntry
	keys, err := q.GetAll(c, &leaderboardEntries)
	if err != nil {
        panic(err.Error())
    }

	// Set the non-model fields
	for i := range leaderboardEntries {
	    leaderboardEntry := &leaderboardEntries[i]
        leaderboardEntry.SetProperties(r)
    }

	return keys, leaderboardEntries
}


func GetSuggestions(r *http.Request, params map[string]interface{}) ([]*datastore.Key, []Suggestion) {
	c := appengine.NewContext(r)
	var q *datastore.Query
	// If parameters were specified
	if params != nil {
		q = GetModelEntities(c, "Suggestion", 0, params)
	} else {
		// Otherwise use query params from url
		q = WebQueryEntities(c, r, "Suggestion", 0)
	}
	var suggestions []Suggestion
	keys, err := q.GetAll(c, &suggestions)
	if err != nil {
        panic(err.Error())
    }

	// Set the non-model fields
	for i := range suggestions {
	    suggestion := &suggestions[i]
        suggestion.SetProperties(keys[i])
    }

	return keys, suggestions
}


type UserTotal struct {
	Username    string
	UserID      string
	Points      int64
	Wins        int64
	Medals      []*datastore.Key
	Suggestions int64
	Level       int64
}


type UserTotals []UserTotal

//Set up how to sort UserTotal
func (slice UserTotals) Len() int {
    return len(slice)
}


func (slice UserTotals) Less(i, j int) bool {
    return slice[i].Points < slice[j].Points;
}


func (slice UserTotals) Swap(i, j int) {
    slice[i], slice[j] = slice[j], slice[i]
}


const levelSize int64 = 30

func AddToUserTotal(userTotals *map[string]*UserTotal, leaderboardEntry LeaderboardEntry) {
	// If there isn't an entry for this user yet
	if _, ok := (*userTotals)[leaderboardEntry.UserID]; !ok {
		var medals []*datastore.Key
		ut := &UserTotal{
			Username: leaderboardEntry.Username,
			UserID: leaderboardEntry.UserID,
			Points: 0,
			Wins: 0,
			Medals: medals,
			Suggestions: 0,
			Level: 0,
		}
		(*userTotals)[leaderboardEntry.UserID] = ut
	}
	// Add the points, wins, medals, and suggestions for this user
	(*userTotals)[leaderboardEntry.UserID].Points = (*userTotals)[leaderboardEntry.UserID].Points + leaderboardEntry.Points
	(*userTotals)[leaderboardEntry.UserID].Wins = (*userTotals)[leaderboardEntry.UserID].Wins + leaderboardEntry.Wins
	(*userTotals)[leaderboardEntry.UserID].Suggestions = (*userTotals)[leaderboardEntry.UserID].Suggestions + leaderboardEntry.Suggestions
	// Add the medal keys to the list of medals
	for i := range leaderboardEntry.Medals {
		(*userTotals)[leaderboardEntry.UserID].Medals = append((*userTotals)[leaderboardEntry.UserID].Medals, leaderboardEntry.Medals[i])
	}
	// Calculate the user level
	(*userTotals)[leaderboardEntry.UserID].Level = ((*userTotals)[leaderboardEntry.UserID].Points / levelSize) + 1
}


func GetLeaderboardStats(r *http.Request, userID interface{}, startDate interface{}, endDate interface{}) UserTotals {
	c := appengine.NewContext(r)
	leaderboardStatParams := make(map[string]interface{})
	// If a user was specified
	if userIDString, ok := userID.(string); ok {
		// Add it to the query params
		leaderboardStatParams["user_id"] = userIDString
	}
	_, leaderboardEntries := GetLeaderboardEntries(r, leaderboardStatParams)

	// Initialize the Stats map
	userTotals := make(map[string]*UserTotal)
	for i := range leaderboardEntries {
	    leaderboardEntry := &leaderboardEntries[i]
		// Get the show key and load data
    	var show Show
    	datastore.Get(c, leaderboardEntry.Show, &show)
		// If start and end date were specified
        if startDateTime, ok := startDate.(time.Time); ok {
			if endDateTime, ok := endDate.(time.Time); ok {
				// If the entry falls within the date span
				if show.Created.After(startDateTime) && show.Created.Before(endDateTime) || show.Created.Equal(startDateTime) {
					log.Println("leaderboardEntry.Show.Created: ", show.Created)
					AddToUserTotal(&userTotals, *leaderboardEntry)
				}
			}
		} else {
			log.Println("leaderboardEntry.Show.Created: ", show.Created)
			AddToUserTotal(&userTotals, *leaderboardEntry)
		}
    }
	//Turn the map into a slice
	var orderedTotals UserTotals
	for _, userTotal := range userTotals {
		orderedTotals = append(orderedTotals, *userTotal)
	}
	// Order the totals by points
	sort.Sort(orderedTotals)
	return orderedTotals
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
*/