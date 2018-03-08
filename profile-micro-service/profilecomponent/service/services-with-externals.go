package service

import (
	"errors"
	"strconv"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"github.com/adriendomoison/apigoboot/errorhandling/apihelper"
	"github.com/adriendomoison/apigoboot/errorhandling/servicehelper"
	"github.com/adriendomoison/apigoboot/profile-micro-service/config"
	"github.com/adriendomoison/apigoboot/profile-micro-service/profilecomponent/rest"
)

var BaseUrl = config.GAppUrl + "/api/private-v1"

func askUserServiceForUserId(email string) (rest.ResponseDTOUserInfo, *servicehelper.Error) {
	if resDTO, apiErrors, statusCode := callGetUserIdService(email); len(apiErrors.Errors) > 0 {
		return rest.ResponseDTOUserInfo{}, &servicehelper.Error{
			Detail:  errors.New(apiErrors.Errors[0].(apihelper.Error).Detail),
			Message: apiErrors.Errors[0].(apihelper.Error).Message,
			Param:   apiErrors.Errors[0].(apihelper.Error).Param,
			Code:    servicehelper.Code(statusCode),
		}
	} else {
		return resDTO, nil
	}
}

func askUserServiceForUserEmail(userId uint) (rest.ResponseDTOUserInfo, *servicehelper.Error) {
	if resDTO, apiErrors, statusCode := callGetUserEmailService(userId); len(apiErrors.Errors) > 0 {
		return rest.ResponseDTOUserInfo{}, &servicehelper.Error{
			Detail:  errors.New(apiErrors.Errors[0].(apihelper.Error).Detail),
			Message: apiErrors.Errors[0].(apihelper.Error).Message,
			Param:   apiErrors.Errors[0].(apihelper.Error).Param,
			Code:    servicehelper.Code(statusCode),
		}
	} else {
		return resDTO, nil
	}
}

// callGetUserIdService ask the user micro service for a user id
func callGetUserIdService(email string) (rest.ResponseDTOUserInfo, apihelper.ApiErrors, int) {

	// call api
	req, err := http.NewRequest("GET", BaseUrl+"/user/email/" + email, nil)

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
	req, err := http.NewRequest("GET", BaseUrl+"/user/id/" + strconv.Itoa(int(userId)), nil)

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