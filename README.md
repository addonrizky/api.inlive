# api.inlive.app
The backend API for inlive.app.

## API Documentation
We're using Swagger with Swag library to automate the API documentation. Check this URL for avaiable endpoints.

## About The Repo
### Directory Structure
We're following the Golang directory best practice from [here](https://github.com/golang-standards/project-layout). And these are the current structure we're use:
```
- internal # all the application code that only used by this service
    - middlerware # middleware used for router
    - models # all models and data related code
    - routes # API routes structured with the file path
        - v1 # version for the API
            - auth # routes for /api/v1/auth/*
            - stream # routes for /api/v1/stream/*
- pkg # package or module code that can be used by other service
- tests # test directory
- sdpcollection # directory for auto generated SDP file
```

### Library Used
Libraries used here are selected based on the Asumsi project best practices. It should be modular, easy to replace or upgrade, not a framework that tight couple with lot of components.

- Router: Mux
- Database : Gorm for ORM, Golang Migrate for database versioning
- Documentation Helper: Swag
- Streaming : Pion for WebRTC library, FFMPEG for video encoding


### Running in local
You need to install docker in your local first, then disable buildkit in preference>docker engine in feature section.

1. Clone this repo
2. Run `go mod vendor` to install modules locally
3. Copy `.env.example` to `.env` and set the environment variables in `.env` (ask other engineers if you don't already know the variables)
4. run `docker-compose up` will run air command that enable hot reload with volume mapping to your local source directory.
5. the api will be available to addres http://localhost:8080. If port 8080 is not availble, you can change it through `.env` on PORT variable.



### How To Contribute
### Adding new endpoint
Each new endpoint should be placed inside `internal/routes/`. Each endpoint  should have it's own directory and each method of endpoint will have it's own file. Then inside that directory, it should have `index.go` with `GetRoutesHandlers()` function to return all endpoint paths and method like below.

```
package stream

import (
	"github.com/asumsi/api.inlive/pkg/api"
)

func GetRoutesHandlers() []router.Route {
	routes := []router.Route{ 
        {Path: "/streams/{id}/init", Handler: Init, Method: http.MethodPost},
		{Path: "/streams/{id}start", Handler: Start, Method: http.MethodPost},
        {Path: "/streams/{id}/end", Handler: End, Method: http.MethodPost},
	}
	return routes
}

```

The path will automatically append to `API_PREFIX`  set at `internal/main.go`. Currently it set to `/v1` So if you declare the path as `/stream/init` the real path will be `/v1/stream/init`

For example if we want to create new endpoint `/v1/hello` then this is what we do:
1. Create directory `internal/routes/v1/hello`
2. Create index file `internal/routes/v1/hello/index.go` 
3. Create method file like below

    ```
    package hello

    import(
        "net/http"

        "github.com/asumsi/api.inlive/pkg/api"
    )

    type Hello struct{
        Message string
    }

    func World(w http.ResponseWriter, r *http.Request) {
        api.RespondJSON(w, api.Response{Code: http.StatusOK, Message: "Hello World", Data: Hello{Message:"World"}})
    }
    ```

    We use API helper function `api.RespondJSON(http.ResponseWriter,Response)` to respond the request with JSON.

3. Declare the endpoint route to `index.go`
    ```
    import (
        "net/http"

        "github.com/asumsi/api.inlive/pkg/router"
    )

    func GetRoutesHandlers() map[string]router.Route {
        routes := map[string]router.Route{
            "/hello/world":  {Handler: World, Method: http.MethodGet},
        }
        return routes
    }
    ```

#### Add New Package
If you like to add new package either on models, or pkg. Below are some convetions that you need to follow:
1. Make sure the package directory name is represent it's function.
2. Make the same package and directory name.
3. Create `types.go` file for all types declaration.
4. Create the same filename with package name to put all functions. For example package name `stream` will have directory name `stream` and inside the directory there will be `stream.go` file. 
5. Consider to split the function into seperate file if it's contain some private functions.

### Testing
We have `tests` directory that can be run by calling `go test ./tests` or verbose mode ` go test -v ./tests`. We have example in `tests/streams_test.go` how to test the stream endpoint.

TODO:
Create all stream test by comparing the value in database with the response from http request

### How to test the feature
