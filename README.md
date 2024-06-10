# Real-Time-Forum ⚡
...
## Setup
- ### Golang Source code:
    To build the project from source code:
    1. Install or Update to [Golang ^1.22](https://go.dev/doc/install)
    2. Download this project:
        ```
        git clone https://zone01normandie.org/git/cramanan/real-time-forum
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
    This project comes with a Dockerfile. You will need to run these commands:
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
This API is fully written in Golang and its response are formatted in JSON using [Golang json package](https://pkg.go.dev/encoding/json):

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

To send a register / login request, you must format it into JSON and allow cookies in order to open a session on the server:

The following fields are required by the auth endpoints:
- Register endpoint:
    ```json
        // endpoint : http://localhost:8080/api/register
        {
            "email" : "<example>@<address>.<any>",  
            "name" : "<example>",                   
            "password": "<password>",               
            "age" : "<number>",                     // as a string
            "gender" : "<any>",                     
            "firstName" : "<example>",              
            "lastName" : "<example>"                
        }
    ```
- Login endpoint:
    ```json
        // endpoint : http://localhost:8080/api/login
        {
            "email" : "<example>@<address>.<any>",
            "password": "<secret-password>"
        }
    ```