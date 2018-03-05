package service

import (
	"os"
	"context"
	"io/ioutil"
	"encoding/json"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"github.com/huandu/facebook"
	"github.com/elithrar/simple-scrypt"
)

// CheckCredentials redirect user authentication to the right method depending of the authType
func (s *service) CheckCredentials(username string, password string, authType string) (userId uint, ok bool) {
	if authType == "password" {
		return checkCredentialsForPasswordAuth(s, username, password)
	} else if authType == "facebook" {
		return checkCredentialsForFacebookAuth(s, username, password)
	} else if authType == "google" {
		return checkCredentialsForGoogleAuth(s, username, password)
	}
	return 0, false
}

// CheckCredentialsForPasswordAuth check user credentials in database
func checkCredentialsForPasswordAuth(s *service, email string, password string) (userId uint, ok bool) {
	if entity, err := s.repo.FindByEmail(email); err != nil {
		return 0, false
	} else if err := scrypt.CompareHashAndPassword([]byte(entity.Password), []byte(password)); err != nil {
		return 0, false
	} else {
		return entity.ID, true
	}
}

// CheckCredentialsForPasswordAuth check user credentials by contacting facebook GraphAPI
func checkCredentialsForFacebookAuth(s *service, facebookUserId string, accessToken string) (userId uint, ok bool) {
	res, err := facebook.Get("/"+facebookUserId, facebook.Params{
		"fields":       "email",
		"access_token": accessToken,
	})
	if err != nil {
		return 0, false
	}
	if user, err := s.repo.FindByEmail(res["email"].(string)); err != nil {
		return 0, false
	} else {
		return user.ID, true
	}
}

// CheckCredentialsForGoogleAuth check user credentials by contacting google API
func checkCredentialsForGoogleAuth(s *service, _ string, accessToken string) (userId uint, ok bool) {
	conf := &oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_LOGIN_API_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_LOGIN_API_SECRET_ID"),
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
		},
		Endpoint: google.Endpoint,
	}

	tok := oauth2.Token{
		AccessToken: accessToken,
	}

	client := conf.Client(context.Background(), &tok)
	email, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil {
		return 0, false
	}
	defer email.Body.Close()
	data, _ := ioutil.ReadAll(email.Body)

	var dat map[string]interface{}
	if err := json.Unmarshal(data, &dat); err != nil {
		return 0, false
	}

	if user, err := s.repo.FindByEmail(dat["email"].(string)); err != nil {
		return 0, false
	} else {
		return user.ID, true
	}
}
