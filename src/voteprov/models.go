package voteprov


import (
	"fmt"
	//"log"
    "time"
	"strconv"
	"net/http"
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


type Show struct {
	ID              int64            `datastore:"ID" json:"id,omitempty"`
	VoteLength      int64            `datastore:"vote_length" json:"vote_length,omitempty"`
	ResultLength    int64            `datastore:"result_length" json:"result_length,omitempty"`
	VoteOptions     int64            `datastore:"vote_options" json:"vote_options,omitempty"`
	Timezone        string           `datastore:"timezone" json:"timezone,omitempty"`
	VoteTypes       []*datastore.Key `datastore:"vote_types" json:"vote_types,omitempty"`
	Players         []*datastore.Key `datastore:"players" json:"players,omitempty"`
	PlayerPool      []*datastore.Key `datastore:"player_pool" json:"player_pool,omitempty"`
	Created         time.Time        `datastore:"created" json:"created,-"`
	Archived        bool             `datastore:"archived" json:"archived,omitempty"`
	CurrentVoteType *datastore.Key   `datastore:"current_vote_type" json:"current_vote_type,omitempty"`
	CurrentVoteInit time.Time        `datastore:"current_vote_init" json:"-"`
	RecapType       *datastore.Key   `datastore:"recap_type" json:"recap_type,omitempty"`
	RecapInit       time.Time        `datastore:"recap_init" json:"-"`
	Locked          bool             `datastore:"locked" json:"locked,omitempty"`
}


func (show *Show) SetProperties(showKey *datastore.Key) {
	// Get the show id in string format
	show.ID = showKey.IntID()
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
	//log.Println("Leaderboard Show: ", showIDString)
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
