package handle

import (
	"reflect"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/urfave/cli/v2"
	"github.com/zfs123/go-ac-router/utils"
)

const (
	ApiTypeAction uint8 = iota
	CliTypeAction
)

// Action is used to get value from command or http request
type Action interface {
	Time(name string) *time.Time
	Int(name string) int
	Int64(name string) int64
	Float64(name string) float64
	String(name string) string
	Bool(name string) bool
	Uint(name string) uint
	Uint64(name string) uint64
	GetActionType() uint8
	ShouldBind(params interface{}) error
}

// CliAction is used to get input from http request
type ApiAction struct {
	C *gin.Context
}

// Create an api action
func NewApiAction(c *gin.Context) *ApiAction {
	return &ApiAction{c}
}

// Convert timestamp in parameters to time type
func (api *ApiAction) Time(name string) *time.Time {
	val := getValueFromQueryPost(api.C, name)
	i, _ := strconv.ParseInt(val, 10, 64)
	t := time.Unix(i, 0)
	return &t
}

// Convert to int
// 0 if not found
func (api *ApiAction) Int(name string) int {
	val := getValueFromQueryPost(api.C, name)
	i, _ := strconv.Atoi(val)
	return i
}

// Convert to int64
// 0 if not found
func (api *ApiAction) Int64(name string) int64 {
	val := getValueFromQueryPost(api.C, name)
	i, _ := strconv.ParseInt(val, 10, 64)
	return i
}

// Convert to float64
// 0 if not found
func (api *ApiAction) Float64(name string) float64 {
	val := getValueFromQueryPost(api.C, name)
	f, _ := strconv.ParseFloat(val, 64)
	return f
}

// Convert to string
func (api *ApiAction) String(name string) string {
	return getValueFromQueryPost(api.C, name)
}

// Convert to bool
// false if not found
func (api *ApiAction) Bool(name string) bool {
	val := getValueFromQueryPost(api.C, name)
	b, _ := strconv.ParseBool(val)
	return b
}

// Convert to uint
// 0 if not found
func (api *ApiAction) Uint(name string) uint {
	val := getValueFromQueryPost(api.C, name)
	i, _ := strconv.Atoi(val)
	return uint(i)
}

// Convert to uint64
// 0 if not found
func (api *ApiAction) Uint64(name string) uint64 {
	val := getValueFromQueryPost(api.C, name)
	u, _ := strconv.ParseUint(val, 10, 64)
	return u
}

// Get action type
func (api *ApiAction) GetActionType() uint8 {
	return ApiTypeAction
}

// Should bind
func (api *ApiAction) ShouldBind(params interface{}) error {
	return api.C.ShouldBind(params)
}

// CliAction is used to get input from command
type CliAction struct {
	C *cli.Context
}

// Create a cli action
func NewCliAction(c *cli.Context) *CliAction {
	return &CliAction{c}
}

// Return time format
func (cli *CliAction) Time(name string) *time.Time {
	return cli.C.Timestamp(name)
}

// Return int format
func (cli *CliAction) Int(name string) int {
	return cli.C.Int(name)
}

// Return int64 format
func (cli *CliAction) Int64(name string) int64 {
	return cli.C.Int64(name)
}

// Return float64 format
func (cli *CliAction) Float64(name string) float64 {
	return cli.C.Float64(name)
}

// Return string format
func (cli *CliAction) String(name string) string {
	return cli.C.String(name)
}

// Return bool format
func (cli *CliAction) Bool(name string) bool {
	return cli.C.Bool(name)
}

// Return uint format
func (cli *CliAction) Uint(name string) uint {
	return cli.C.Uint(name)
}

// Return uint64 format
func (cli *CliAction) Uint64(name string) uint64 {
	return cli.C.Uint64(name)
}

// Get action type
func (cli *CliAction) GetActionType() uint8 {
	return CliTypeAction
}

// Currently supported field types are not perfect
func (cli *CliAction) ShouldBind(params interface{}) error {
	return utils.RangeStruct(params, func(value reflect.Value, field reflect.StructField) bool {
		alia := utils.GetForm(field)
		if alia == "" {
			return true
		}
		switch value.Interface().(type) {
		case bool:
			value.SetBool(cli.Bool(alia))
		case int:
			value.SetInt(cli.Int64(alia))
		case int64:
			value.SetInt(cli.Int64(alia))
		case float64:
			value.SetFloat(cli.Float64(alia))
		case string:
			value.SetString(cli.String(alia))
		case uint:
			value.SetUint(cli.Uint64(alia))
		case uint64:
			value.SetUint(cli.Uint64(alia))
		}
		return true
	})
}

// Preference to get parameters from post data
// get from url, if empty
func getValueFromQueryPost(c *gin.Context, key string) string {
	f := c.PostForm(key)
	if f != "" {
		return f
	}
	return c.Query(key)
}
