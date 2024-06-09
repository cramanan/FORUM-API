# Real-Time-Forum ⚡
...

## The API
This API is fully written in Golang and its response are formatted in JSON using the [json package](https://pkg.go.dev/encoding/json):

Example:
```json
    // http://localhost:8080/api/users
    [
        {
            "id": "c1ec9487-5c5c-46a4-a72c-947f71bd5c88",
            "email": "wally@mail.com",
            "name": "Wally",
            "created": "2024-06-08T19:16:19.247313513Z"
        },
    ]
```

A lot of endpoints are Protected so you will need to register / login first to open a Session and access the API.

To send a register / login request, you must format it into JSON:

Request Body (Make sure to allow Cookies):

```json
    // http://localhost:8080/api/login
    {
        "email" : "<example>@<address>.<any>",
        "password": "<secret-password>"
    }
```