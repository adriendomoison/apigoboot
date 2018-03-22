// Package rest implement the callback required by the user package
package rest

import (
	"encoding/json"
	"github.com/adriendomoison/apigoboot/api-tool/errorhandling/apihelper"
	"github.com/adriendomoison/apigoboot/user-micro-service/config"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
)

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

func (r *rest) ValidateAccessToken(c *gin.Context) {

	authorizationCode := c.Request.Header.Get("Authorization")

	if authorizationCode == "" || len(authorizationCode) <= 7 {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}

	email := c.Param("email")

	if len(email) == 0 {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}

	tokenUserId, _, err := askOauthServiceForTokenOwnerUserId(authorizationCode[7:])
	if err != nil {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}

	if ok, err := r.service.IsThatTheUserId(email, tokenUserId); err != nil {
		c.JSON(apihelper.BuildResponseError(err))
		return
	} else if !ok {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}
	c.Next()
}
