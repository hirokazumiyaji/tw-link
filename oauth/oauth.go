package oauth

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"
)

const Endpoint = "https://api.twitter.com"

var OAuthNonceRegex = regexp.MustCompile("[^a-zA-Z0-9]")

type OAuth struct {
	ConsumerKey     string
	ConsumerSecret  string
	Token           string
	TokenSecret     string
	SignatureMethod string
	Version         string
}

func (o *OAuth) Timestamp() string {
	return fmt.Sprintf("%d", time.Now().Unix())
}

func (o *OAuth) Nonce() string {
	var bytes = make([]byte, 32)
	rand.Read(bytes)
	nonce := base64.StdEncoding.EncodeToString(bytes)
	return OAuthNonceRegex.ReplaceAllString(nonce, "")
}

func New(consumerKey, consumerSecret, token, tokenSecret string) *OAuth {
	return &OAuth{
		consumerKey,
		consumerSecret,
		token,
		tokenSecret,
		"HMAC-SHA1",
		"1.0",
	}
}

func (o *OAuth) Header(method, requestUrl string, values url.Values) string {
	signatureBaseString := method + "&" + percentEncode(url.QueryEscape(requestUrl)) + "&" + percentEncode(values.Encode())
	signingKey := percentEncode(url.QueryEscape(o.ConsumerSecret)) + "&" + percentEncode(url.QueryEscape(o.TokenSecret))
	hmac := hmac.New(sha1.New, []byte(signingKey))
	hmac.Write([]byte(signatureBaseString))
	return base64.StdEncoding.EncodeToString(hmac.Sum(nil))
}

func (o *OAuth) RequestToken() error {
	client := http.Client{}
	req, err := http.NewRequest("POST", Endpoint+"/oauth/request_token", nil)
	if err != nil {
		return err
	}

	values := url.Values{}
	values.Add("oauth_consumer_key", o.ConsumerKey)
	values.Add("oauth_nonce", o.Nonce())
	values.Add("oauth_signature_method", o.SignatureMethod)
	values.Add("oauth_timestamp", o.Timestamp())

	header := o.Header("POST", Endpoint+"/oauth/request_token", values)

	req.Header.Add("Authorization", header)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	fmt.Println(res)
	return nil
}

func percentEncode(str string) string {
	return strings.Replace(strings.Replace(str, "+", "%20", -1), "*", "%2A", -1)
}
