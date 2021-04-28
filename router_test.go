package acrouter

import (
	"net/http"
	"testing"

	"github.com/zfs123/go-ac-router/handle"
)

type EmptyRequest struct {
}

type EmptyResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func TestRouter(t *testing.T) {
	r, err := New()
	if err != nil {
		t.Fatal(err)
	}
	r.AddApiRoute("/hello", "GET", "hello api", EmptyRequest{}, EmptyResponse{}, hello)
	r.Run()
}

func hello(action handle.Action, response handle.Response) {
	resp := EmptyResponse{
		Code: 0,
		Msg:  "hello world",
	}
	response.Response(http.StatusOK, resp)
}
