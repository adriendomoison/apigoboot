package rest

import (
	"errors"
	"net/http"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/RangelReale/osin"
	userrepo "github.com/adriendomoison/go-boot-api/user/repo"
	userservice "github.com/adriendomoison/go-boot-api/user/service"
)

func HandleLoginPage(ar *osin.AuthorizeRequest, c *gin.Context) (uint, bool) {
	if c.Request.Method == "POST" {
		userService := userservice.New(userrepo.New())
		c.Request.ParseForm()
		if userId, ok := userService.CheckCredentials(c.Request.Form.Get("username"), c.Request.Form.Get("password"), "password"); ok {
			return userId, true
		}
	}
	c.HTML(http.StatusOK, "authentication.tmpl", gin.H{
		"client_id":     ar.Client.GetId(),
		"authorize_url": c.Request.URL,
	})
	return 0, false
}

func DownloadAccessToken(url string, auth *osin.BasicAuth, output map[string]interface{}) error {
	// download access token
	preq, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return err
	}
	if auth != nil {
		preq.SetBasicAuth(auth.Username, auth.Password)
	}

	pclient := &http.Client{}
	presp, err := pclient.Do(preq)
	if err != nil {
		return err
	}
	if presp.StatusCode != 200 {
		return errors.New("invalid status code")
	}

	jdec := json.NewDecoder(presp.Body)
	err = jdec.Decode(&output)
	return err
}
