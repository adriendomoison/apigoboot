# apigoboot

A Go API boilplate using a micro-service architecture and integrating Oauth2 secured login, payment, and emails.

### Set up

You need to add in your environment a `GOOGLE_LOGIN_API_CIENT_ID` and `GOOGLE_LOGIN_API_SECRET_ID` to allow oauth2 to process login with google

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
An `Errors` array is returned composed of objects including the problematic `param` and a `detail` human readable sentence to describe each error.

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
