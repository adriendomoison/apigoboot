# apigoboot

A Go API boilplate using a micro-service architecture and integrating Oauth2 secured login, payment, and emails.

### Return values

The API return user friendly error message that can be printed directly client-side.
The errors are always using the same nomenclature.

```
{
    "Errors": [
        {
            "param": "last_name",
            "detail": "The field last_name is required",
            "type": "https://api.apigoboot.com/validation-error/profile#last_name"
        },
        {
            "param": "email",
            "detail": "The field email is required",
            "type": "https://api.apigoboot.com/validation-error/profile#email"
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
            "type": "https://api.apigoboot.com/validation-error/transaction#scheduled_date"
        }
    ]
}
```
