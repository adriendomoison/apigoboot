// Package service implement the services required by the rest package
package service

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/adriendomoison/apigoboot/errorhandling/servicehelper"
	"github.com/adriendomoison/apigoboot/user-micro-service/component/user/rest"
	"github.com/elithrar/simple-scrypt"
	"github.com/huandu/facebook"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"io/ioutil"
	"os"
)

// CheckCredentials redirect user authentication to the right method depending of the authType
func (s *service) CheckCredentials(reqDTO rest.RequestDTOCheckCredentials) (rest.ResponseDTOUserInfo, *servicehelper.Error) {
	if reqDTO.AuthType == "password" {
		return checkCredentialsForPasswordAuth(s, reqDTO.Username, reqDTO.Password)
	} else if reqDTO.AuthType == "facebook" {
		return checkCredentialsForFacebookAuth(s, reqDTO.Username, reqDTO.Password)
	} else if reqDTO.AuthType == "google" {
		return checkCredentialsForGoogleAuth(s, reqDTO.Username, reqDTO.Password)
	}
	return rest.ResponseDTOUserInfo{}, &servicehelper.Error{
		Detail: errors.New("incorrect username or password"),
	}
}

// CheckCredentialsForPasswordAuth check user credentials in database
func checkCredentialsForPasswordAuth(s *service, email string, password string) (rest.ResponseDTOUserInfo, *servicehelper.Error) {
	if entity, err := s.repo.FindByEmail(email); err != nil {
		return rest.ResponseDTOUserInfo{}, &servicehelper.Error{
			Detail: errors.New("incorrect username or password"),
		}
	} else if err := scrypt.CompareHashAndPassword([]byte(entity.Password), []byte(password)); err != nil {
		return rest.ResponseDTOUserInfo{}, &servicehelper.Error{
			Detail: errors.New("incorrect username or password"),
		}
	} else {
		return rest.ResponseDTOUserInfo{
			UserId: entity.ID,
			Email:  entity.Email,
		}, nil
	}
}

// CheckCredentialsForPasswordAuth check user credentials by contacting facebook GraphAPI
func checkCredentialsForFacebookAuth(s *service, facebookUserId string, accessToken string) (rest.ResponseDTOUserInfo, *servicehelper.Error) {
	res, err := facebook.Get("/"+facebookUserId, facebook.Params{
		"fields":       "email",
		"access_token": accessToken,
	})
	if err != nil {
		return rest.ResponseDTOUserInfo{}, &servicehelper.Error{
			Detail: errors.New("incorrect username or password"),
		}
	}
	if user, err := s.repo.FindByEmail(res["email"].(string)); err != nil {
		return rest.ResponseDTOUserInfo{}, &servicehelper.Error{
			Detail: errors.New("incorrect username or password"),
		}
	} else {
		return rest.ResponseDTOUserInfo{
			UserId: user.ID,
			Email:  user.Email,
		}, nil
	}
}

// CheckCredentialsForGoogleAuth check user credentials by contacting google API
func checkCredentialsForGoogleAuth(s *service, _ string, accessToken string) (rest.ResponseDTOUserInfo, *servicehelper.Error) {
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
		return rest.ResponseDTOUserInfo{}, &servicehelper.Error{
			Detail: errors.New("error while trying to check credentials"),
		}
	}
	defer email.Body.Close()
	data, _ := ioutil.ReadAll(email.Body)

	var dat map[string]interface{}
	if err := json.Unmarshal(data, &dat); err != nil {
		return rest.ResponseDTOUserInfo{}, &servicehelper.Error{
			Detail: errors.New("error while trying to check credentials"),
		}
	}

	if user, err := s.repo.FindByEmail(dat["email"].(string)); err != nil {
		return rest.ResponseDTOUserInfo{}, &servicehelper.Error{
			Detail: errors.New("incorrect username or password"),
		}
	} else {
		return rest.ResponseDTOUserInfo{
			UserId: user.ID,
			Email:  user.Email,
		}, nil
	}
}
