package handle

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/urfave/cli/v2"
)

// Response is used to response
type Response interface {
	Response(code int, data interface{})
	SendSimpleOk(msg string)
	SendSimpleFail(msg string)
}

// ApiResponse implemented response of http request
type ApiResponse struct {
	C *gin.Context
}

// Create an api response
func NewApiResponse(c *gin.Context) *ApiResponse {
	return &ApiResponse{c}
}

// Output json
func (resp *ApiResponse) Response(code int, data interface{}) {
	resp.C.JSON(code, data)
}

// Simple send success
func (resp *ApiResponse) SendSimpleOk(msg string) {
	resp.C.JSON(http.StatusOK, map[string]interface{}{"code": 0, "msg": msg})
}

// Simple send error
func (resp *ApiResponse) SendSimpleFail(msg string) {
	resp.C.JSON(http.StatusInternalServerError, map[string]interface{}{"code": 1, "msg": msg})
}

// ApiResponse implemented response of command
type CliResponse struct {
	C *cli.Context
}

// Create an cli response
func NewCliResponse(c *cli.Context) *CliResponse {
	return &CliResponse{c}
}

// Format cli output
func (resp *CliResponse) Response(code int, data interface{}) {
	b, err := json.Marshal(data)
	if err != nil {
		_, _ = fmt.Fprintf(resp.C.App.Writer, "code %d, msg %s, err %s\n", code, "output failed", err.Error())
	}
	_, _ = fmt.Fprintf(resp.C.App.Writer, "code %d,msg %s\n", code, string(b))
}

// Simple send success
func (resp *CliResponse) SendSimpleOk(msg string) {
	_, _ = fmt.Fprintln(resp.C.App.Writer, msg)
}

// Simple send error
func (resp *CliResponse) SendSimpleFail(msg string) {
	_, _ = fmt.Fprintln(resp.C.App.Writer, msg)
}
