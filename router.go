package acrouter

import (
	"reflect"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/urfave/cli/v2"
	"github.com/zfs123/go-ac-router/v1/handle"
	"github.com/zfs123/go-ac-router/v1/utils"
)

type Router struct {
	api *ApiServer
	cli *CliServer
}

func NewRouter(api *ApiServer, cli *CliServer) *Router {
	return &Router{
		api: api,
		cli: cli,
	}
}

// Generate cli and api routes simultaneously
func (r *Router) AddMultiRoute(path string, method string, description string, params interface{}, response interface{}, handleFunc handle.Func) {
	r.AddApiRoute(path, method, description, params, response, handleFunc)
	r.AddCliCommandByStruct(path[1:], description, params, handleFunc)
}

// Add cli route by struct
func (r *Router) AddCliCommandByStruct(path string, description string, params interface{}, handleFunc handle.Func) {
	r.cli.AddCommand(&cli.Command{
		Name:    strings.Replace(path, "/", "_", -1),
		Aliases: []string{path},
		Usage:   description,
		Flags:   buildCliFlag(params),
		Action: func(c *cli.Context) error {
			handleFunc(handle.NewCliAction(c), handle.NewCliResponse(c))
			return nil
		},
	})
}

// Add cli route by command
func (r *Router) AddCliCommand(c *cli.Command) {
	r.cli.AddCommand(c)
}

// Add api route and add doc item for generate documents
func (r *Router) AddApiRoute(path string, method string, description string, params interface{}, response interface{}, handleFunc handle.Func) {
	autoAddApiRoute(r.api.Engine, path, method, func(context *gin.Context) {
		handleFunc(handle.NewApiAction(context), handle.NewApiResponse(context))
	})
}

func autoAddApiRoute(engine *gin.Engine, path string, method string, handleFunc gin.HandlerFunc) {
	switch method {
	case "GET":
		engine.GET(path, handleFunc)
	case "POST":
		engine.POST(path, handleFunc)
	case "DELETE":
		engine.DELETE(path, handleFunc)
	case "PATCH":
		engine.PATCH(path, handleFunc)
	case "PUT":
		engine.PUT(path, handleFunc)
	case "OPTIONS":
		engine.OPTIONS(path, handleFunc)
	case "HEAD":
		engine.HEAD(path, handleFunc)
	case "Any":
		engine.Any(path, handleFunc)
	}
}

// Generate cli command parameters by the structure
func buildCliFlag(flag interface{}) (fields []cli.Flag) {
	_ = utils.RangeStruct(flag, func(value reflect.Value, field reflect.StructField) bool {
		alia := utils.GetForm(field)
		if alia == "" {
			return true
		}
		description := utils.GetDescription(field)
		require := utils.GetRequired(field)
		var flag cli.Flag
		switch value.Interface().(type) {
		case bool:
			flag = &cli.BoolFlag{Name: alia, Usage: description, Required: require}
		case int:
			flag = &cli.IntFlag{Name: alia, Usage: description, Required: require}
		case int64:
			flag = &cli.Int64Flag{Name: alia, Usage: description, Required: require}
		case float64:
			flag = &cli.Float64Flag{Name: alia, Usage: description, Required: require}
		case string:
			flag = &cli.StringFlag{Name: alia, Usage: description, Required: require}
		case []string:
			flag = &cli.StringSliceFlag{Name: alia, Usage: description, Required: require}
		case time.Duration:
			flag = &cli.DurationFlag{Name: alia, Usage: description, Required: require}
		case uint:
			flag = &cli.UintFlag{Name: alia, Usage: description, Required: require}
		case uint64:
			flag = &cli.Uint64Flag{Name: alia, Usage: description, Required: require}
		case time.Time:
			flag = &cli.TimestampFlag{Name: alia, Usage: description, Required: require}
		}
		fields = append(fields, flag)
		return true
	})
	return
}
