# Real-Time-Forum ⚡
<sub>From [01-edu](https://github.com/01-edu/public) public course</sub>

...

Note: This is a project not a package

## Libraries used

- For the Database:

    [database/sql](https://pkg.go.dev/database/sql) : Golang's standard library for SQL database that requires a driver.

    [mattn/go-sqlite3](https://pkg.go.dev/github.com/mattn/go-sqlite3) : This library provides the sqlite3 driver needed by database/sql.

    [context](https://pkg.go.dev/context) : Golang's standard library Context type used to begin SQL transactions and handle concurrency.

- For the HTTP Server:
    [net/http](https://pkg.go.dev/net/http) :
    I created my own API type that inherits from Golang's http.Server. I implemented a built-in Session Storage and a DB Storage. [See API struct](/api/api.go)

    [encoding/json](https://pkg.go.dev/encoding/json) : 
    The API only communicates in JSON. HTTP requests bodies are decoded from raw JSON and responses are encoded in JSON. [See examples](#the-api)

    [gorilla/websocket](https://pkg.go.dev/github.com/gorilla/websocket) : 
    The API also provide a WebSocket connection for live-chatting. The messages are also encoded/decoded in JSON.



## Setup
- ### Golang Source code:
    To build the project from source code:
    1. Install or Update to [Golang ^1.22](https://go.dev/doc/install)
    2. Download this project:
        ```
        git clone https://github.com/cramanan/FORUM-API
        ```
    3. Download the dependencies:
        ```
        go mod download
        ```
    4. Run the program:
        ```
        go run .
        ```
- ### Docker:
    This project comes with a Dockerfile. You will need to run these commands (you might need to run it as root):
    1. Build the image:
        ```
        docker image build -t <image-name> .
        ```
    2. Run the container:
        ```
        docker run --name <container-name> -p 8080:8080 <image-name>
        ```
    3. Remove them when you're done:
        ```
        docker rm <container-name>
        docker rmi <image-name>
        ```

## The API
This API is fully written in Golang and its response are formatted in JSON:

Example:

Request URL : http://localhost:8080/api/users

Response Body:
```json
    [
        {
            "id": "c1ec9487-5c5c-46a4-a72c-947f71bd5c88",
            "email": "wally@mail.com",
            "name": "Wally",
            "created": "2024-06-08T19:16:19.247313513Z"
        },

        ...
    ]
```

A lot of endpoints are Protected so you will need to register / login first to open a Session and access the API.

To register or login you can either:

- Open the connect page and fill the forms (allow cookies keep the session ID)
    
- Send an HTTP POST request with a raw JSON body

The following fields are required by the auth endpoints:
- Register endpoint:

    Request URL: http://localhost:8080/api/login

    Request Body:
    ```json
        {
            "email" : "<example>@<address>.<any>",  
            "name" : "<example>",                   
            "password": "<password>",
            "age" : "<number>",
            "gender" : "<any>",
            "firstName" : "<example>",              
            "lastName" : "<example>"                
        }
    ```
- Login endpoint:

    Request URL: http://localhost:8080/api/login

    Request Body:
    ```json
        {
            "email" : "<example>@<address>.<any>",
            "password": "<secret-password>"
        }
    ```
