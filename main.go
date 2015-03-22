package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"

	"github.com/garyburd/go-oauth/oauth"
	"github.com/naoya/go-pit"
)

const endpoint = "https://api.twitter.com"

var httpRegexp = regexp.MustCompile("https?://.*")

type tweet struct {
	Text string `json:"text"`
	Id   string `json:"id_str"`
	User struct {
		ScreenName      string `json:"screen_name"`
		ProfileImageUrl string `json:"profile_image_url"`
	} `json:"user"`
}

var oauthClient = oauth.Client{
	TemporaryCredentialRequestURI: endpoint + "/oauth/request_token",
	ResourceOwnerAuthorizationURI: endpoint + "/oauth/authenticate",
	TokenRequestURI:               endpoint + "/oauth/access_token",
}

type Options struct {
	Count              int
	SinceId            int
	MaxId              int
	TrimUser           bool
	ExcludeReplies     bool
	ContributorDetails bool
	IncludeEntities    bool
}

func (o *Options) Values() url.Values {
	values := url.Values{}
	if o.Count != 0 {
		values.Add("count", strconv.Itoa(o.Count))
	}
	if o.SinceId != 0 {
		values.Add("since_id", strconv.Itoa(o.SinceId))
	}
	if o.MaxId != 0 {
		values.Add("max_id", strconv.Itoa(o.MaxId))
	}
	if o.TrimUser != true {
		values.Add("trim_user", strconv.FormatBool(o.TrimUser))
	}
	if o.ExcludeReplies != true {
		values.Add("exclude_replies", strconv.FormatBool(o.ExcludeReplies))
	}
	if o.ContributorDetails != true {
		values.Add("contributor_details", strconv.FormatBool(o.ContributorDetails))
	}
	if o.IncludeEntities != false {
		values.Add("include_rts", strconv.FormatBool(o.IncludeEntities))
	}
	return values
}

func getHomeTimeline(token *oauth.Credentials, options *Options) ([]tweet, error) {
	url_ := endpoint + "/1.1/statuses/home_timeline.json"
	values := options.Values()
	oauthClient.SignParam(token, "GET", url_, values)
	url_ = url_ + "?" + values.Encode()
	res, err := http.Get(url_)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	if res.StatusCode != 200 {
		return nil, err
	}

	var timeline []tweet
	if err = json.NewDecoder(res.Body).Decode(&timeline); err != nil {
		return nil, err
	}
	return timeline, nil
}

func filterTimeline(timeline []tweet) []tweet {
	_timeline := make([]tweet, len(timeline))

	for _, t := range timeline {
		if httpRegexp.MatchString(t.Text) {
			_timeline = append(_timeline, t)
		}
	}
	return _timeline
}

func main() {
	config, err := pit.Get("twitter.com")
	if err != nil {
		fmt.Println("pit Error %v", err)
		os.Exit(1)
	}
	oauthClient.Credentials.Token = config["consumer_key"]
	oauthClient.Credentials.Secret = config["consumer_secret"]

	oauthToken, isFoundOAuthToken := config["access_token"]
	oauthSecret, isFoundOAuthSecret := config["access_token_secret"]
	var token *oauth.Credentials
	if isFoundOAuthToken && isFoundOAuthSecret {
		token = &oauth.Credentials{oauthToken, oauthSecret}
	} else {
		token, err = oauthClient.RequestTemporaryCredentials(http.DefaultClient, "", nil)
		if err != nil {
			fmt.Println("failed to request temporary credentials:", err)
			os.Exit(1)
		}
	}

	timeline, err := getHomeTimeline(token, &Options{Count: 200})
	if err != nil {
		os.Exit(1)
	}

	timeline = filterTimeline(timeline)
	fmt.Println(timeline)
}
