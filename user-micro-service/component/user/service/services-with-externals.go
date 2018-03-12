// Package service implement the services required by the rest package
package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/adriendomoison/apigoboot/errorhandling/apihelper"
	"github.com/adriendomoison/apigoboot/errorhandling/servicehelper"
	"github.com/adriendomoison/apigoboot/user-micro-service/component/user/rest"
	"github.com/adriendomoison/apigoboot/user-micro-service/config"
	"io/ioutil"
	"net/http"
)

var baseUrl = config.GAppUrl + "/api/private-v1"

// AddWithProfile set up and create a user with a profile
func (s *service) AddWithProfile(reqDTO rest.RequestDTOWithProfile) (rest.ResponseDTOWithProfile, *servicehelper.Error) {
	if _, err := s.Add(rest.RequestDTO{
		Email:    reqDTO.Email,
		Username: reqDTO.Username,
		Password: reqDTO.Password,
	}); err != nil {
		return rest.ResponseDTOWithProfile{}, err
	}
	resDTO, apiErrors, statusCode := callPostProfileService(reqDTO)
	if len(apiErrors.Errors) > 0 {
		s.Remove(resDTO.Email)
		return rest.ResponseDTOWithProfile{}, &servicehelper.Error{
			Detail:  errors.New(apiErrors.Errors[0].(apihelper.Error).Detail),
			Message: apiErrors.Errors[0].(apihelper.Error).Message,
			Param:   apiErrors.Errors[0].(apihelper.Error).Param,
			Code:    servicehelper.Code(statusCode),
		}
	}
	return resDTO, nil
}

// RetrieveWithProfile retrieve a user with its profile
func (s *service) RetrieveWithProfile(email string) (rest.ResponseDTOWithProfile, *servicehelper.Error) {
	if _, err := s.Retrieve(email); err != nil {
		return rest.ResponseDTOWithProfile{}, err
	}
	resDTO, apiErrors, statusCode := callGetProfileService(email)
	if len(apiErrors.Errors) > 0 {
		s.Remove(resDTO.Email)
		return rest.ResponseDTOWithProfile{}, &servicehelper.Error{
			Detail:  errors.New(apiErrors.Errors[0].(apihelper.Error).Detail),
			Message: apiErrors.Errors[0].(apihelper.Error).Message,
			Param:   apiErrors.Errors[0].(apihelper.Error).Param,
			Code:    servicehelper.Code(statusCode),
		}
	}
	return resDTO, nil
}

// callPostProfileService ask the profile micro service to create a profile for the user
func callPostProfileService(reqDTO rest.RequestDTOWithProfile) (rest.ResponseDTOWithProfile, apihelper.ApiErrors, int) {

	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(reqDTO)

	// call api
	req, err := http.NewRequest("POST", baseUrl+"/profiles", b)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// read response
	body, _ := ioutil.ReadAll(resp.Body)

	userWithProfile := rest.ResponseDTOWithProfile{}
	apiErrors := apihelper.ApiErrors{}
	json.Unmarshal(body, &userWithProfile)
	json.Unmarshal(body, &apiErrors)

	return userWithProfile, apiErrors, resp.StatusCode
}

// callPostProfileService ask the profile micro service to create a profile for the user
func callGetProfileService(email string) (rest.ResponseDTOWithProfile, apihelper.ApiErrors, int) {

	// call api
	req, err := http.NewRequest("GET", baseUrl+"/profiles/"+email, nil)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// read response
	body, _ := ioutil.ReadAll(resp.Body)

	userWithProfile := rest.ResponseDTOWithProfile{}
	apiErrors := apihelper.ApiErrors{}
	json.Unmarshal(body, &userWithProfile)
	json.Unmarshal(body, &apiErrors)

	return userWithProfile, apiErrors, resp.StatusCode
}
