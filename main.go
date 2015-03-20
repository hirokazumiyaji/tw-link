package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/naoya/go-pit"
)

const Endpoint = "https://api.twitter.com"

type AccessToken struct {
	TokenType   string `json:"token_type"`
	AccessToken string `json:"access_token"`
}

type Timeline struct {
	Coordinates interface{} `json:"coordinates"`
	Truncated   bool        `json:"truncated"`
	CreatedAt   string      `json:"created_at"`
	Favorited   bool        `json:"favorited"`
	IdStr       string      `json:"id_str"`
	Text        string      `json:"text"`
}

func getToken(consumerKey, consumerSecret string) (*AccessToken, error) {
	bearerToken := base64.StdEncoding.EncodeToString([]byte(consumerKey + ":" + consumerSecret))

	values := url.Values{}
	values.Add("grant_type", "client_credentials")
	req, err := http.NewRequest("POST", Endpoint+"/oauth2/token", strings.NewReader(values.Encode()))
	req.Header.Add("Authorization", "Basic "+bearerToken)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded;charset=UTF-8")
	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println("Error\ttype:request token\terr:%v", err)
		return nil, err
	}

	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)
	var t AccessToken
	if err = decoder.Decode(&t); err != nil {
		fmt.Println("Error\ttype:decode token\terr:%v", err)
		return nil, err
	}

	return &t, nil
}

type TimelineOptions struct {
	UserId             int
	ScreenName         string
	SinceId            int
	Count              int
	MaxId              int
	TrimUser           bool
	ExcludeReplies     bool
	ContributorDetails bool
	IncludeRts         bool
}

func getTimeline(token string, options TimelineOptions) (*Timeline, error) {
	values := url.Values{}
	if options.UserId != 0 {
		values.Add("user_id", strconv.Itoa(options.UserId))
	}
	if options.ScreenName != "" {
		values.Add("screen_name", options.ScreenName)
	}
	if options.SinceId != 0 {
		values.Add("since_id", strconv.Itoa(options.SinceId))
	}
	if options.Count != 0 {
		values.Add("count", strconv.Itoa(options.Count))
	}
	if options.MaxId != 0 {
		values.Add("max_id", strconv.Itoa(options.MaxId))
	}
	if options.TrimUser != true {
		values.Add("trim_user", strconv.FormatBool(options.TrimUser))
	}
	if options.ExcludeReplies != true {
		values.Add("exclude_replies", strconv.FormatBool(options.ExcludeReplies))
	}
	if options.ContributorDetails != true {
		values.Add("contributor_details", strconv.FormatBool(options.ContributorDetails))
	}
	if options.IncludeRts != false {
		values.Add("include_rts", strconv.FormatBool(options.IncludeRts))
	}

	req, _ := http.NewRequest("GET", Endpoint+"/1.1/statuses/user_timeline.json", nil)
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("Content-Type", "application/json; charset=utf-8")
	req.URL.RawQuery = values.Encode()
	client := http.Client{}

	res, err := client.Do(req)
	if err != nil {
		fmt.Println("Error\ttype:request timeline\terr:%v", err)
		return nil, err
	}

	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)
	var timeline interface{}
	if err = decoder.Decode(&timeline); err != nil {
		fmt.Println("Error\ttype:decode timeline\terr:%v", err)
		return nil, err
	}
	fmt.Println(timeline)
	return &Timeline{}, nil
}

func main() {
	config, err := pit.Get("twitter.com")
	if err != nil {
		fmt.Println("pit Error %v", err)
		os.Exit(1)
	}
	consumerKey := config["consumer_key"]
	consumerSecret := config["consumer_secret"]

	token, err := getToken(consumerKey, consumerSecret)
	if err != nil {
		os.Exit(1)
	}

	timeline, err := getTimeline(token.AccessToken, TimelineOptions{ScreenName: "hirokazu_miyaji", Count: 200})
	if err != nil {
		os.Exit(1)
	}
	fmt.Println(timeline)
}
