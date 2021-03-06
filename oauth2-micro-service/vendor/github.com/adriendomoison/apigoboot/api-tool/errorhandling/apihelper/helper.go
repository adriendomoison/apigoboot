// Package apihelper automate construction of http response
package apihelper

import (
	"errors"
	"github.com/adriendomoison/apigoboot/api-tool/errorhandling/servicehelper"
	"github.com/gin-gonic/gin"
	"gopkg.in/go-playground/validator.v8"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

// ApiError interface for all API error messages
type ApiError interface {
}

// Error is the default error message structure the API returns
type Error struct {
	apiError ApiError
	Param    string `json:"param"`
	Detail   string `json:"detail"`
	Message  string `json:"message"`
}

// ApiErrors carry the list of errors returned by the API from a request
type ApiErrors struct {
	Errors []ApiError
}

// BuildRequestError build a usable JSON error object from an error string generated by the structure validator
func BuildRequestError(err error) (int, ApiErrors) {
	var apiErrors ApiErrors
	switch err.(type) {
	case validator.ValidationErrors:
		for _, v := range err.(validator.ValidationErrors) {
			var validationError Error
			validationError.Param = toSnakeCase(v.Field)
			validationError.Detail = "Field validation for " + toSnakeCase(v.Field) + " failed on the " + v.Tag + " tag."
			if v.Tag == "required" {
				validationError.Message = "This field is required"
			}
			if v.Tag == "email" {
				validationError.Message = "Invalid email address. Valid e-mail can contain only latin letters, numbers, '@' and '.'"
			}
			if v.Tag == "url" {
				validationError.Message = "Invalid URL address. Valid URL start with http:// or https://"
			}
			apiErrors.Errors = append(apiErrors.Errors, validationError)
		}
		return http.StatusBadRequest, apiErrors
	default:
		apiErrors.Errors = append(apiErrors.Errors, Error{
			Detail: err.Error(),
		})
		return http.StatusBadRequest, apiErrors
	}
}

// BuildResponseError apply the right status to the http response and build the error JSON object
func BuildResponseError(err *servicehelper.Error) (status int, apiErrors ApiErrors) {
	apiErrors.Errors = append(apiErrors.Errors, Error{
		Detail:  err.Detail.Error(),
		Message: err.Message,
		Param:   err.Param,
	})
	return int(err.Code), apiErrors
}

// BuildResponseError apply the right status to the http response and build the error JSON object
func BuildChatfuelResponseError(err *servicehelper.Error) (int, gin.H) {
	var apiError ApiErrors
	apiError.Errors = append(apiError.Errors, Error{
		Detail:  err.Detail.Error(),
		Message: err.Message,
		Param:   err.Param,
	})
	return int(err.Code), gin.H{"set_attributes": apiError.Errors[0]}
}

// GetBoolQueryParam allow to retrieve a boolean query parameter.
// It takes the gin contect as param to build error if queryParam is not formatted correctly and a default value to set value if the parameter is optional and not set.
func GetBoolQueryParam(c *gin.Context, value *bool, queryParam string, defaultValue bool) bool {
	var err error
	if c.Query(queryParam) != "" {
		if *value, err = strconv.ParseBool(c.Query(queryParam)); err != nil {
			c.JSON(BuildRequestError(
				errors.New("query parameter '" + queryParam + "' value should be true or false (omit the key and default value '" + strconv.FormatBool(defaultValue) + "' will be applied)")),
			)
			return false
		}
	} else {
		*value = defaultValue
	}
	return true
}

// toSnakeCase change a string to it's snake case version
func toSnakeCase(str string) string {
	snake := regexp.MustCompile("(.)([A-Z][a-z]+)").ReplaceAllString(str, "${1}_${2}")
	snake = regexp.MustCompile("([a-z0-9])([A-Z])").ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}
