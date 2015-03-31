package voteprov


import (
	"fmt"
	//"log"
    "time"
	"strconv"
	"net/http"
	"appengine"
	"appengine/datastore"
)

type Player struct {
	Name          string    `datastore:"name" json:"name,omitempty"`
	PhotoFilename string    `datastore:"photo_filename,noindex" json:"photo_filename,omitempty"`
	Star          bool      `datastore:"star" json:"star,omitempty"`
	DateAdded     time.Time `datastore:"date_added" json:"-"`
	IMGPath       string    `json:"img_path,omitempty"`
}


func (p *Player) SetProperties() {
    p.IMGPath = fmt.Sprintf("/static/img/players/%s", p.PhotoFilename)
}


type SuggestionPool struct {
	Name          string    `datastore:"name" json:"name,omitempty"`
	DisplayName   string    `datastore:"display_name" json:"display_name,omitempty"`
	Description   string    `datastore:"description" json:"description,omitempty"`
	Created       time.Time `datastore:"created" json:"created,omitempty"`
}


type VoteType struct {
	// Defined at creation
	Name                string         `datastore:"name" json:"name,omitempty"`
	DisplayName         string         `datastore:"display_name" json:"display_name,omitempty"`
	SuggestionPool      *datastore.Key `datastore:"suggestion_pool" json:"suggestion_pool,omitempty"`
	PreshowVoted        bool           `datastore:"preshow_voted" json:"preshow_voted,omitempty"`
	HasIntervals        bool           `datastore:"has_intervals" json:"has_intervals,omitempty"`
	IntervalUsesPlayers bool           `datastore:"interval_uses_players" json:"interval_uses_players,omitempty"`
	Intervals           []int64        `datastore:"intervals" json:"intervals,omitempty"`
	Style               string         `datastore:"style" json:"style,omitempty"`
	Occurs              string         `datastore:"occurs" json:"occurs,omitempty"`
	Ordering            int64          `datastore:"ordering" json:"ordering,omitempty"`
	Options             int64          `datastore:"options" json:"options,omitempty"`
	RandomizeAmount     int64          `datastore:"randomize_amount" json:"randomize_amount,omitempty"`
	ButtonColor         string         `datastore:"button_color" json:"button_color,omitempty"`

	// Dynamic
	CurrentInterval     int64          `datastore:"current_interval" json:"current_interval,omitempty"`
	CurrentInit         time.Time      `datastore:"current_init" json:"current_init,omitempty"`
}

// Option display for voting
type OptionDisplay struct {
	OptionKey string `json:"option_key,omitempty"`
	PlayerKey string `json:"player_key,omitempty"`
}

// Used to capture the current state of the show
type CurrentState struct {
	State              string           `json:"state,omitempty"`
	Display            string           `json:"display,omitempty"`
	DisplayName        string           `json:"display_name,omitempty"`
	Style              string           `json:"style,omitempty"`
	Interval           int64            `json:"interval,omitempty"`
	CurrentTime        time.Time        `json:"current_time,omitempty"`
	DisplayEndTime     time.Time        `json:"display_end_time,omitempty"`
	Speedup            bool             `json:"speedup,omitempty"`
	RemainingIntervals map[string]int64 `json:"remaining_intervals,omitempty"`
	OptionDisplay      []OptionDisplay  `json:"options,omitempty"`
}

type Show struct {
	ID                  int64            `datastore:"ID" json:"id,omitempty"`
	VoteLength          int64            `datastore:"vote_length" json:"vote_length,omitempty"`
	ResultLength        int64            `datastore:"result_length" json:"result_length,omitempty"`
	VoteOptions         int64            `datastore:"vote_options" json:"vote_options,omitempty"`
	Timezone            string           `datastore:"timezone" json:"timezone,omitempty"`
	VoteTypes           []*datastore.Key `datastore:"vote_types" json:"vote_types,omitempty"`
	Players             []*datastore.Key `datastore:"players" json:"players,omitempty"`
	PlayerPool          []*datastore.Key `datastore:"player_pool" json:"player_pool,omitempty"`
	Created             time.Time        `datastore:"created" json:"created,-"`
	Archived            bool             `datastore:"archived" json:"archived,omitempty"`
	CurrentVoteType     *datastore.Key   `datastore:"current_vote_type" json:"current_vote_type,omitempty"`
	CurrentVoteInit     time.Time        `datastore:"current_vote_init" json:"-"`
	RecapType           *datastore.Key   `datastore:"recap_type" json:"recap_type,omitempty"`
	RecapInit           time.Time        `datastore:"recap_init" json:"-"`
	Locked              bool             `datastore:"locked" json:"locked,omitempty"`
	CurrentShowInterval *datastore.Key   `json:"current_show_interval,omitempty"`
	CurrentState        CurrentState     `json:"current_state,omitempty"`
}


func GetCurrentShowInterval(show *Show, showKey *datastore.Key, r *http.Request) (*datastore.Key, ShowInterval) {
	var showIntervalKey *datastore.Key
	var showInterval ShowInterval
	// Get query params
	qParams := MapQuery(r)
	var voteTypeID string
	var ok bool
	voteTypeInterface := qParams["vote_type"]
	// Get the current vote type if it exists
	if voteTypeID, ok = voteTypeInterface.(string); !ok {
		// Get the vote type id in string format
		voteTypeID = strconv.FormatInt(show.CurrentVoteType.IntID(), 10)
	}
	// Only set current show interval if there's a vote type
	if len(voteTypeID) > 0 {
		var interval *int64
		showIDString := strconv.FormatInt(show.ID, 10)
		// Get the current show interval if it exists in the query
		if intervalInterface, ok := qParams["interval"]; ok {
			intervalString := intervalInterface.(string)
			*interval, _ = strconv.ParseInt(intervalString, 0, 64)
		} else {
			// Otherwise, get the current interval from the vote type
			voteTypeParams := map[string]interface{}{
				"vote_type": voteTypeID,
				"show": showIDString,
			}
			_, voteType := GetVoteType(r, false, voteTypeParams)
			interval = &voteType.CurrentInterval
		}
		// if a current interval exists
		if interval != nil {
			intervalString := strconv.FormatInt(*interval, 10)
			showIntervalParams := map[string]interface{}{"show": showIDString,
				                                         "vote_type": voteTypeID,
												         "interval": intervalString}
			showIntervalKey, showInterval = GetShowInterval(r, false, showIntervalParams)
		}
	}
	return showIntervalKey, showInterval
}


func (show *Show) SetProperties(c appengine.Context, showKey *datastore.Key, r *http.Request) {
	//var currentShowInterval ShowInterval
	var voteType VoteType
	// Get the show id in string format
	show.ID = showKey.IntID()
	// Get the current show interval key and values
	show.CurrentShowInterval, _ = GetCurrentShowInterval(show, showKey, r)
	currentTime := time.Now().UTC()
	// Try to load the data into the vote type struct model
	if err := datastore.Get(c, show.CurrentVoteType, &voteType); err != nil {
		panic(err.Error())
	}
	// Set vote type non-model properties
	// Initialize state variables
	currentDisplay := "default"
	currentShowState := "default"
	var currentDisplayName string
	var currentStyle string
	var currentInterval int64
	var displayEndTime time.Time
	// If the current time is after the current vote init
	if currentTime.After(show.CurrentVoteInit) {
		votingEndTime := show.CurrentVoteInit.Add(time.Duration(show.VoteLength)*time.Second)
		resultEndTime := votingEndTime.Add(time.Duration(show.ResultLength)*time.Second)
		// If we're currently in the voting period
		if currentTime.Before(votingEndTime) {
			currentDisplay = "voting"
			currentShowState = voteType.Name
			currentDisplayName = voteType.DisplayName
			currentStyle = voteType.Style
			currentInterval = voteType.CurrentInterval
			displayEndTime = votingEndTime
		} else if currentTime.Before(resultEndTime) {
			// Otherwise if we're in the display period
			currentDisplay = "result"
			currentShowState = voteType.Name
			currentDisplayName = voteType.DisplayName
			currentStyle = voteType.Style
			currentInterval = voteType.CurrentInterval
			displayEndTime = resultEndTime
		}
	}
	show.CurrentState = CurrentState{State: currentShowState,
									 Display: currentDisplay,
									 DisplayName: currentDisplayName,
									 Style: currentStyle,
									 Interval: currentInterval,
									 CurrentTime: currentTime,
		                             DisplayEndTime: displayEndTime,
									 }
}


type Suggestion struct {
	ID             int64          `datastore:"ID" json:"id,omitempty"`
	Show           *datastore.Key `datastore:"show" json:"show,omitempty"`
	SuggestionPool *datastore.Key `datastore:"suggestion_pool" json:"suggestion_pool,omitempty"`
	Used           bool           `datastore:"used" json:"used,omitempty"`
	VotedOn        bool           `datastore:"voted_on" json:"voted_on,omitempty"`
	AmountVotedOn  int64          `datastore:"amount_voted_on" json:"amount_voted_on,omitempty"`
	Value          string         `datastore:"value" json:"value,omitempty"`
	PreshowValue   int64          `datastore:"preshow_value" json:"preshow_value,omitempty"`
	SessionID      string         `datastore:"session_id" json:"-"`
	UserID         string         `datastore:"user_id" json:"user_id,omitempty"`
	Created        time.Time      `datastore:"created" json:"-"`
}


// Need an md5 of session id to validate pre-show from (so we don't straight up give away the id)
func (suggestion *Suggestion) SetProperties(suggestionKey *datastore.Key) {
	// Get the show id in string format
	suggestion.ID = suggestionKey.IntID()
}


type PreshowVote struct {
	Show           *datastore.Key `datastore:"show" json:"show,omitempty"`
	Suggestion     *datastore.Key `datastore:"suggestion" json:"suggestion,omitempty"`
	SessionID      string         `datastore:"session_id" json:"session_id,omitempty"`
}


type LiveVote struct {
	Show       *datastore.Key `datastore:"show" json:"show,omitempty"`
	VoteType   *datastore.Key `datastore:"vote_type" json:"vote_type,omitempty"`
	Player     *datastore.Key `datastore:"player" json:"player,omitempty"`
	Suggestion *datastore.Key `datastore:"suggestion" json:"suggestion,omitempty"`
	Interval   int64          `datastore:"interval" json:"interval,omitempty"`
	SessionID  string         `datastore:"session_id" json:"session_id,omitempty"`
	UserID     string         `datastore:"user_id" json:"user_id,omitempty"`
}


type ShowInterval struct {
	Show       *datastore.Key `datastore:"show" json:"show,omitempty"`
	VoteType   *datastore.Key `datastore:"vote_type" json:"vote_type,omitempty"`
	Interval   int64          `datastore:"interval" json:"interval,omitempty"`
	Player     *datastore.Key `datastore:"player" json:"player,omitempty"`
}


type VoteOptions struct {
	Show       *datastore.Key   `datastore:"show" json:"show,omitempty"`
	VoteType   *datastore.Key   `datastore:"vote_type" json:"vote_type,omitempty"`
	Interval   int64            `datastore:"interval" json:"interval,omitempty"`
	OptionList []*datastore.Key `datastore:"option_list" json:"option_list,omitempty"`
}


type VotedItem struct {
	VoteType   *datastore.Key `datastore:"vote_type" json:"vote_type,omitempty"`
	Show       *datastore.Key `datastore:"show" json:"show,omitempty"`
	Suggestion *datastore.Key `datastore:"suggestion" json:"suggestion,omitempty"`
	Player     *datastore.Key `datastore:"player" json:"player,omitempty"`
	Interval   int64          `datastore:"interval" json:"interval,omitempty"`
}


type Medal struct {
	Name          string `datastore:"name" json:"name,omitempty"`
	DisplayName   string `datastore:"display_name" json:"display_name,omitempty"`
	Description   string `datastore:"description" json:"description,omitempty"`
	ImageFilename string `datastore:"image_filename" json:"image_filename,omitempty"`
	IconFilename  string `datastore:"icon_filename" json:"icon_filename,omitempty"`
}


type LeaderboardEntry struct {
	Show        *datastore.Key   `datastore:"show" json:"show,omitempty"`
	ShowDate    time.Time        `datastore:"show_date" json:"-"`
	UserID      string           `datastore:"user_id" json:"user_id,omitempty"`
	Points      int64            `datastore:"points" json:"points,omitempty"`
	Wins        int64            `datastore:"wins" json:"wins,omitempty"`
	Medals      []*datastore.Key `datastore:"medals" json:"medals,omitempty"`
	Username    string           `json:"username"`
	Suggestions int64            `json:"suggestions"`
}


func (le *LeaderboardEntry) SetProperties(r *http.Request) {
	// Try to get user profile by user id
	userProfileParams := map[string]interface{}{"user_id": le.UserID}
	_, userProfile := GetUserProfiles(r, false, userProfileParams)
	// Set the Username
	le.Username = userProfile.Username
	// Get the show id in string format
	showIDString := strconv.FormatInt(le.Show.IntID(), 10)
	suggestionParams := map[string]interface{}{
		"user_id": le.UserID,
		"show": showIDString,
	}
	_, suggestions := GetSuggestions(r, suggestionParams)
	// Set the number of suggestions by the user for the show
	le.Suggestions = int64(len(suggestions))
}


type UserProfile struct {
	UserID         string    `datastore:"user_id" json:"user_id,omitempty"`
	Username       string    `datastore:"username" json:"username,omitempty"`
	StripUsername  string    `datastore:"strip_username" json:"strip_username,omitempty"`
	Email          string    `datastore:"email" json:"-"`
	LoginType      string    `datastore:"login_type" json:"-"`
	CurrentSession string    `datastore:"current_session" json:"-"`
	FBAccessToken  string    `datastore:"fb_access_token" json:"-"`
	Created        time.Time `datastore:"created" json:"-"`
}


type EmailOptOut struct {
	Email string `datastore:"email" json:"email"`
}


type LeaderboardSpan struct {
	Name      string    `datastore:"name" json:"name,omitempty"`
	StartDate time.Time `datastore:"start_date" json:"start_date,omitempty"`
	EndDate   time.Time `datastore:"end_date" json:"end_date,omitempty"`
}
