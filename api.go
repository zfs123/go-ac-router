package acrouter

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/zfs123/go-ac-router/logger"
	"go.uber.org/zap"
)

type ApiServer struct {
	Engine   *gin.Engine
	ipAddr   string
	port     int
	certFile string
	keyFile  string
}

func NewApiServer(ipAddr string, port int) *ApiServer {
	gin.SetMode(gin.ReleaseMode)
	s := &ApiServer{
		ipAddr: ipAddr,
		port:   port,
		Engine: gin.New(),
	}
	s.setMiddleware()
	return s
}

func NewApiTlsServer(ipAddr string, port int, certFile, keyFile string) *ApiServer {
	s := NewApiServer(ipAddr, port)
	s.keyFile = keyFile
	s.certFile = certFile
	return s
}

// Initialize middleware and use zap to record log
// Recovery can record log and recover when the request crashes
func (aps *ApiServer) setMiddleware() {
	aps.Engine.Use(func(c *gin.Context) {
		start := time.Now()
		c.Next()
		end := time.Now()
		latency := end.Sub(start)
		path := c.Request.URL.Path
		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()
		logger.Info("api request",
			zap.Int("status_code", statusCode),
			zap.Duration("latency", latency),
			zap.String("client_ip", clientIP),
			zap.String("method", method),
			zap.String("path", path),
		)
	})
	aps.Engine.Use(gin.Recovery())
}

// Turn on debug mode
func (aps *ApiServer) SetDebug() {
	gin.SetMode(gin.DebugMode)
}

// Set no route handle
func (aps *ApiServer) SetNoRoute(handlerFunc gin.HandlerFunc) {
	aps.Engine.NoRoute(handlerFunc)
}

// Run api server
func (aps *ApiServer) Run() error {
	return aps.Engine.Run(aps.ipAddr + ":" + strconv.Itoa(aps.port))
}

// Run api tls server
func (aps *ApiServer) RunTLS() error {
	return aps.Engine.RunTLS(aps.ipAddr+":"+strconv.Itoa(aps.port), aps.certFile, aps.keyFile)
}
