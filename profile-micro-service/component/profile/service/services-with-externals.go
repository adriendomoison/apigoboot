// Package service implement the services required by the rest package
package service

import (
	"encoding/json"
	"errors"
	"github.com/adriendomoison/apigoboot/errorhandling/apihelper"
	"github.com/adriendomoison/apigoboot/errorhandling/servicehelper"
	"github.com/adriendomoison/apigoboot/profile-micro-service/component/profile/rest"
	"io/ioutil"
	"net/http"
	"strconv"
)

// TODO make it "http://api.user.apigoboot/api/private-v1/user"
var userBaseUrl = "http://api.user:4200/api/private-v1/user"

func askUserServiceForUserId(email string) (rest.ResponseDTOUserInfo, *servicehelper.Error) {
	resDTO, apiErrors, statusCode := callGetUserIdService(email)
	if len(apiErrors.Errors) > 0 {
		return rest.ResponseDTOUserInfo{}, &servicehelper.Error{
			Detail:  errors.New(apiErrors.Errors[0].(apihelper.Error).Detail),
			Message: apiErrors.Errors[0].(apihelper.Error).Message,
			Param:   apiErrors.Errors[0].(apihelper.Error).Param,
			Code:    servicehelper.Code(statusCode),
		}
	}
	return resDTO, nil
}

func askUserServiceForUserEmail(userId uint) (rest.ResponseDTOUserInfo, *servicehelper.Error) {
	resDTO, apiErrors, statusCode := callGetUserEmailService(userId)
	if len(apiErrors.Errors) > 0 {
		return rest.ResponseDTOUserInfo{}, &servicehelper.Error{
			Detail:  errors.New(apiErrors.Errors[0].(apihelper.Error).Detail),
			Message: apiErrors.Errors[0].(apihelper.Error).Message,
			Param:   apiErrors.Errors[0].(apihelper.Error).Param,
			Code:    servicehelper.Code(statusCode),
		}
	}
	return resDTO, nil
}

// callGetUserIdService ask the user micro service for a user id
func callGetUserIdService(email string) (rest.ResponseDTOUserInfo, apihelper.ApiErrors, int) {

	// call api
	req, err := http.NewRequest("GET", userBaseUrl+"/email/"+email, nil)

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

	return userInfo, apiErrors, resp.StatusCode
}

// callGetUserIdService ask the user micro service for a user id
func callGetUserEmailService(userId uint) (rest.ResponseDTOUserInfo, apihelper.ApiErrors, int) {

	// call api
	req, err := http.NewRequest("GET", userBaseUrl+"/id/"+strconv.Itoa(int(userId)), nil)

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

	return userInfo, apiErrors, resp.StatusCode
}
