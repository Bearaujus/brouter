# BROUTER
Bearaujus Router (brouter) is a golang router that make your service run even more easier!

# Installation
`go get github.com/Bearaujus/brouter`

# Example: Simple Usecase
- Code
```go
package main

import (
	"github.com/Bearaujus/brouter"
	"github.com/sirupsen/logrus"
	"net/http"
)

func main() {
	// create new brouter instance
	router := brouter.NewBRouter()

	// add route
	router.Route(brouter.StructRoute{
		Pattern:     "/test",
		Methods:     []string{http.MethodGet, http.MethodOptions},
		HandlerFunc: SampleHandler,
	})

	// serve http
	err := router.Serve("127.0.0.1", 25565)
	if err != nil {
		logrus.Info(err)
	}
}

type Foo struct {
	Bar    string
	Sample string
}

func SampleHandler(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	model := Foo{
		Bar:    "test-bar",
		Sample: "test-sample",
	}

	// you can pass directly any data type to the handler!
	return model, nil
}
```

- Response
```json
{
    "header": {
        "is_success": true
    },
    "data": {
        "Bar": "test-bar",
        "Sample": "test-sample"
    }
}
```

# Example: Route File Server
- Code
```go
package main

import (
	"github.com/Bearaujus/brouter"
	"github.com/sirupsen/logrus"
)

func main() {
	// create new brouter instance
	router := brouter.NewBRouter()

	// route file server
	router.RouteFileServer(brouter.StructRouteFileServer{
		Pattern: "/fs",
		DirPath: "./",
	})

	// serve http
	err := router.Serve("127.0.0.1", 25565)
	if err != nil {
		logrus.Info(err)
	}
}
```

- Response

Your file server at `DirPath` can be accessed from `Pattern`.

# Example: Overriding Default Handler Success / Error Func
- Code
```go
package main

import (
	"encoding/json"
	"errors"
	"github.com/Bearaujus/brouter"
	"github.com/sirupsen/logrus"
	"net/http"
)

func main() {
	// create new brouter instance
	router := brouter.NewBRouter()

	// overriding default handler error func
	router.SetDefaultHandlerErrorFunc(HandlerErrorFunc)

	// overriding default handler success func
	router.SetDefaultHandlerSuccessFunc(HandlerSuccessFunc)

	// add routes
	router.Routes([]brouter.StructRoute{
		{
			Pattern:     "/success",
			Methods:     []string{http.MethodGet, http.MethodOptions},
			HandlerFunc: SampleHandlerSuccess,
		},
		{
			Pattern:     "/error",
			Methods:     []string{http.MethodGet, http.MethodOptions},
			HandlerFunc: SampleHandlerError,
		},
	})

	// serve http
	err := router.Serve("127.0.0.1", 25565)
	if err != nil {
		logrus.Info(err)
	}
}

type Foo struct {
	Bar    string
	Sample string
}

func SampleHandlerSuccess(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	model := Foo{
		Bar:    "test-bar",
		Sample: "test-sample",
	}

	return model, nil
}

func SampleHandlerError(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	return nil, errors.New("sample error")
}

func HandlerErrorFunc(w http.ResponseWriter, r *http.Request, err error) {
	w.Header().Set("Content-Type", "text/plain")

	w.WriteHeader(http.StatusBadRequest)
	_, _ = w.Write([]byte(err.Error()))
}

func HandlerSuccessFunc(w http.ResponseWriter, r *http.Request, data interface{}) {
	payload, _ := json.Marshal(data)

	w.Header().Set("Content-Type", "text/plain")

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(payload)
}
```

- Success Response
```text
{"Bar":"test-bar","Sample":"test-sample"}
```

- Error Response
```text
sample error
```

# Example: Overwriting Handler Success / Error Func
If you set default success / error handler. This is more prioritized to be executed.
- Code
```go
package main

import (
	"encoding/json"
	"errors"
	"github.com/Bearaujus/brouter"
	"github.com/sirupsen/logrus"
	"net/http"
)

func main() {
	// create new brouter instance
	router := brouter.NewBRouter()

	// add routes
	router.Routes([]brouter.StructRoute{
		{
			Pattern:     "/success",
			Methods:     []string{http.MethodGet, http.MethodOptions},
			HandlerFunc: SampleHandlerSuccess,
			// set handler error func
			HandlerErrorFunc: HandlerErrorFunc,
			// set handler success func
			HandlerSuccessFunc: HandlerSuccessFunc,
		},
		{
			Pattern:     "/error",
			Methods:     []string{http.MethodGet, http.MethodOptions},
			HandlerFunc: SampleHandlerError,
			// set handler error func
			HandlerErrorFunc: HandlerErrorFunc,
			// set handler success func
			HandlerSuccessFunc: HandlerSuccessFunc,
		},
	})

	// serve http
	err := router.Serve("127.0.0.1", 25565)
	if err != nil {
		logrus.Info(err)
	}
}

type Foo struct {
	Bar    string
	Sample string
}

func SampleHandlerSuccess(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	model := Foo{
		Bar:    "test-bar",
		Sample: "test-sample",
	}

	return model, nil
}

func SampleHandlerError(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	return nil, errors.New("sample error")
}

func HandlerErrorFunc(w http.ResponseWriter, r *http.Request, err error) {
	w.Header().Set("Content-Type", "text/plain")

	w.WriteHeader(http.StatusBadRequest)
	_, _ = w.Write([]byte(err.Error()))
}

func HandlerSuccessFunc(w http.ResponseWriter, r *http.Request, data interface{}) {
	payload, _ := json.Marshal(data)

	w.Header().Set("Content-Type", "text/plain")

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(payload)
}
```

- Success Response
```text
{"Bar":"test-bar","Sample":"test-sample"}
```

- Error Response
```text
sample error
```