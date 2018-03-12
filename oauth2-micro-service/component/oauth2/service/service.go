package service

import (
	"bytes"
	"encoding/json"
	"github.com/adriendomoison/apigoboot/errorhandling/apihelper"
	"github.com/adriendomoison/apigoboot/errorhandling/servicehelper"
	"github.com/adriendomoison/apigoboot/oauth2-micro-service/component/oauth2/rest"
	"github.com/adriendomoison/apigoboot/oauth2-micro-service/config"
	"github.com/go-errors/errors"
	"github.com/jinzhu/gorm"
	"io/ioutil"
	"net/http"
)

type RepoInterface interface {
	FindByAccessToken(token string) (Access, error)
}

// Access database object
type Access struct {
	gorm.Model
	Client       string `gorm:"NOT NULL"`
	UserId       uint   `gorm:"NOT NULL"`
	Authorize    string `gorm:"NOT NULL"`
	Previous     string `gorm:"NOT NULL"`
	AccessToken  string `gorm:"NOT NULL;PRIMARY KEY"`
	RefreshToken string `gorm:"NOT NULL"`
	ExpiresIn    int32  `gorm:"NOT NULL"`
	Scope        string `gorm:"NOT NULL"`
	RedirectUri  string `gorm:"NOT NULL"`
}

// Authorize database object
type Authorize struct {
	gorm.Model
	Client      string `gorm:"NOT NULL"`
	UserId      uint   `gorm:"NOT NULL"`
	Code        string `gorm:"NOT NULL;PRIMARY KEY"`
	ExpiresIn   int32  `gorm:"NOT NULL"`
	Scope       string `gorm:"NOT NULL"`
	RedirectUri string `gorm:"NOT NULL"`
	State       string `gorm:"NOT NULL"`
}

// Client database object
type Client struct {
	UserId      uint   `gorm:"NOT NULL"`
	Id          string `gorm:"NOT NULL;PRIMARY KEY"`
	Secret      string `gorm:"NOT NULL"`
	RedirectUri string `gorm:"NOT NULL"`
}

// Refresh database object
type Refresh struct {
	gorm.Model
	Token  string `gorm:"NOT NULL;PRIMARY KEY"`
	Access string `gorm:"NOT NULL"`
}

var PrivateBaseUrl = config.GAppUrl + "/api/private-v1"

var _ rest.ServiceInterface = (*service)(nil)

type service struct {
	repo RepoInterface
}

func New(repo RepoInterface) *service {
	return &service{repo}
}

func (s *service) AskUserServiceToCheckCredentials(username string, password string, method string) (rest.ResponseDTOUserInfo, *apihelper.ApiErrors) {

	requestBody := rest.RequestDTOUserCredentials{
		Username: username,
		Password: password,
		Method:   method,
	}

	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(requestBody)

	// call api
	req, err := http.NewRequest("POST", PrivateBaseUrl+"/user/check-credentials", b)
	req.Header.Set("Authorization", "Bearer YYY")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// read response
	body, _ := ioutil.ReadAll(resp.Body)

	userInfo := rest.ResponseDTOUserInfo{}
	apiErrors := apihelper.ApiErrors{}
	json.Unmarshal(body, &userInfo)
	json.Unmarshal(body, &apiErrors)

	if len(apiErrors.Errors) == 0 {
		return userInfo, nil
	} else {
		return rest.ResponseDTOUserInfo{}, &apiErrors
	}
}

// GetResourceOwnerId ask database to retrieve the id of the user owning the access entity with this access_token
func (s *service) GetResourceOwnerId(token string) (rest.ResponseDTOUserInfo, *servicehelper.Error) {
	accessToken, err := s.repo.FindByAccessToken(token)
	if err != nil {
		return rest.ResponseDTOUserInfo{}, &servicehelper.Error{
			Param:  "access_token",
			Detail: errors.New("failed to retrieve access token"),
		}
	}
	return rest.ResponseDTOUserInfo{
		UserId: accessToken.UserId,
	}, nil
}
