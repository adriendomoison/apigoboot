# ApiGoBoot - WIP - Looking for inputs & code reviews! 🚀❤️️🚀❤️🚀

[![Build Status](https://travis-ci.org/adriendomoison/apigoboot.svg?branch=master)](https://travis-ci.org/adriendomoison/apigoboot) [![Coverage Status](https://coveralls.io/repos/github/adriendomoison/apigoboot/badge.svg?branch=master)](https://coveralls.io/github/adriendomoison/apigoboot?branch=master) [![Go Report Card](https://goreportcard.com/badge/github.com/adriendomoison/apigoboot)](https://goreportcard.com/report/github.com/adriendomoison/apigoboot) [![GoDoc](https://godoc.org/github.com/adriendomoison/apigoboot?status.svg)](https://godoc.org/github.com/adriendomoison/apigoboot) [![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

This is a playground for a boilerplate API in Go using a micro service oriented architecture. The architecture will be updated depending on what I'll be reading/learning on the web.

### Set up

- Add `GOOGLE_LOGIN_API_CIENT_ID` and `GOOGLE_LOGIN_API_SECRET_ID` to your environment to allow oauth2 to process sign in/up with google (generated at [https://console.developers.google.com/apis/credentials/oauthclient](https://console.developers.google.com/apis/credentials/oauthclient)
- Add the line `127.0.0.1      api.go.boot` to your `/etc/hosts`
- Create a facebook app and add `api.go.boot:4200` in the field `Site URL` in basic settings ([https://developers.facebook.com/apps](https://developers.facebook.com/apps))

#### Start

```
docker-compose up
```

### Return values

The API return user friendly error message that can be printed directly client-side.
The errors are always using the same nomenclature.

```
{
    "Errors": [
        {
            "param": "last_name",
            "detail": "The field last_name is required",
            "message": "Please complete this field"
        },
        {
            "param": "email",
            "detail": "The field email is required",
            "message": "Please complete this field"
        }
    ]
}
```
An `Errors` array is returned composed of objects including the faulty `param` and a `detail` human readable sentence to describe each error.

When an error happen during the request process but that the submitted request pass the validation, a more detailed error will be returned:

```
{
    "Errors": [
        {
            "status_code": 1003,
            "title": "Operation scheduled_date malformed",
            "detail": "The field scheduled_date need to use the 2006-12-31 format",
            "message": "Please use the 2006-12-31 format"
       }
    ]
}
```
