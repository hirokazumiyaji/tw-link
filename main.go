package main

import (
  "encoding/base64"
  "encoding/json"
  "fmt"
  "net/http"
  "net/url"
  "os"
  "strings"

  "github.com/naoya/go-pit"
)

const Endpoint = "https://api.twitter.com"

type accessToken struct {
  TokenType string `json:"token_type"`
  AccessToken string `json:"access_token"`
}

type timeline struct {
  Coordinates interface {} `json:"coordinates"`
  Truncated bool `json:"truncated"`
  CreatedAt string `json:"created_at"`
  Favorited bool `json:"favorited"`
  IdStr string `json:"id_str"`
  Text string `json:"text"`
}

func getToken(consumerKey, consumerSecret string) (*accessToken, error) {
  bearerToken := base64.StdEncoding.EncodeToString([]byte(consumerKey + ":" + consumerSecret))

  values := url.Values{}
  values.Add("grant_type", "client_credentials")
  req, err := http.NewRequest("POST", Endpoint + "/oauth2/token", strings.NewReader(values.Encode()))
  req.Header.Add("Authorization", "Basic " + bearerToken)
  req.Header.Add("Content-Type", "application/x-www-form-urlencoded;charset=UTF-8")
  client := http.Client{}
  res, err := client.Do(req)
  if err != nil {
    fmt.Println("Error\ttype:request token\terr:%v", err)
    return nil, err
  }

  defer res.Body.Close()

  decoder := json.NewDecoder(res.Body)
  var t accessToken
  if err = decoder.Decode(&t); err != nil {
    fmt.Println("Error\ttype:decode token\terr:%v", err)
    return nil, err
  }

  return &t, nil
}

func getTimeline(token string, max_id int) (, error) {
  values := url.Values{}
  values.Add("count", 200)
  if max_id != 0 {
    values.Add("max_id", max_id)
  }
  req, _ := http.NewRequest("GET", Endpoint + "/1.1/statuses/home_timeline.json", nil)
  req.Header.Add("Authorization", "Bearer " + token)
  req.Header.Add("Content-Type", "application/json; charset=utf-8")
  req.URL.RawQuery = values.Encode()
  client := http.Client{}

  res, err := client.Do(req)
  if err != nil {
    fmt.Println("Error\ttype:request timeline\terr:%v", err)
    return nil, err
  }
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
  fmt.Println(token.AccessToken)
}
