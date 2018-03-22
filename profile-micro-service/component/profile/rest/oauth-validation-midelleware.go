// Package rest implement the callback required by the profile package
package rest

import (
	"encoding/json"
	"github.com/adriendomoison/apigoboot/api-tool/errorhandling/apihelper"
	"github.com/adriendomoison/apigoboot/profile-micro-service/config"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
)

// ResponseDTOUserInfo is the object to map JSON response body when requesting basic user info
type ResponseDTOUserInfo struct {
	Email  string `json:"email"`
	UserId uint   `json:"user_id"`
}

func askOauthServiceForTokenOwnerUserId(token string) (uint, int, *apihelper.ApiErrors) {

	req, err := http.NewRequest("GET", config.GAppUrl+"/api/private-v1/access-token/"+token+"/get-owner", nil)
	// TODO add client credential access token
	//req.Header.Set("Authorization", "Bearer xxx")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	accessTokenOwner := struct {
		UserId uint `json:"user_id"`
	}{}
	if json.Unmarshal(body, &accessTokenOwner) != nil {
		apiErrors := apihelper.ApiErrors{}
		json.Unmarshal(body, &apiErrors)
		return 0, resp.StatusCode, &apiErrors
	}
	return accessTokenOwner.UserId, 0, nil
}

// ValidateAccessToken check oauth2 access token (middleware)
func (r *rest) ValidateAccessToken(c *gin.Context) {

	authorizationCode := c.Request.Header.Get("Authorization")

	if authorizationCode == "" || len(authorizationCode) <= 7 {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}

	publicId := c.Param("profileId")

	if len(publicId) == 0 {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}

	tokenUserId, _, err := askOauthServiceForTokenOwnerUserId(authorizationCode[7:])
	if err != nil {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}

	if ok, err := r.service.IsThatTheUserId(publicId, tokenUserId); err != nil {
		c.JSON(apihelper.BuildResponseError(err))
		return
	} else if !ok {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}
	c.Next()
}
