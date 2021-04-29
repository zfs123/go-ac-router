package acrouter

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli/v2"
	"github.com/zfs123/go-ac-router/handle"
)

type header struct {
	Key   string
	Value string
}

func performRequest(r *Router, method, path string, headers ...header) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, nil)
	for _, h := range headers {
		req.Header.Add(h.Key, h.Value)
	}
	w := httptest.NewRecorder()
	r.api.Engine.ServeHTTP(w, req)
	return w
}

func TestRouteOK(t *testing.T) {
	passed := false
	r, err := New()
	if err != nil {
		t.Fatal(err)
	}
	r.AddApiRoute("/test", "GET", "test api", nil, nil, func(action handle.Action, response handle.Response) {
		passed = true
	})

	w := performRequest(r, "GET", "/test")
	assert.True(t, passed)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestApiRouter(t *testing.T) {
	r, err := New()
	if err != nil {
		t.Fatal(err)
	}
	r.AddApiRoute("/hello", "GET", "hello api", nil, nil, func(action handle.Action, response handle.Response) {
		response.Response(http.StatusOK, "hello")
	})

	w := performRequest(r, "GET", "/hello")
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "\"hello\"", w.Body.String())
}

func TestNotFound(t *testing.T) {
	r, err := New()
	if err != nil {
		t.Fatal(err)
	}

	w := performRequest(r, "GET", "/xxxx")
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func ExampleRun() {
	os.Args = []string{"test"}
	r, _ := New()
	r.Run()

	//Output:
	//NAME:
	//test - A new cli application
	//
	//USAGE:
	//   test [global options] command [command options] [arguments...]
	//
	//COMMANDS:
	//   server, s        start a api server
	//   tls_server, tls  start a api tls server
	//   help, h          Shows a list of commands or help for one command
	//
	//GLOBAL OPTIONS:
	//   --help, -h  show help (default: false)
}

func ExampleRunCli() {
	os.Args = []string{"-", "hello"} // the first string is program name
	r, _ := New()
	r.AddCliCommand(&cli.Command{
		Name:  "hello",
		Usage: "hello cli command",
		Action: func(c *cli.Context) error {
			fmt.Println("hello world")
			return nil
		},
	})
	r.Run()

	//Output:
	//hello world
}
