package ghost

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

//NewHeartBeatHandler returns a new heartbeat handler
func NewHeartBeatHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		now := fmt.Sprintf("%s", time.Now())
		c.JSON(200, now)
	}
}

//NewConfigHandler returns a handler that prints the configuration
func NewConfigHandler(config interface{}) gin.HandlerFunc {
	return func(c *gin.Context) { c.JSON(http.StatusOK, config) }
}

//NewDefaultAdminHandler returns an http handler with hearbeat and config
//routes set up
func NewDefaultAdminHandler(config interface{}) http.Handler {
	admin := gin.New()
	admin.Use(gin.Recovery())
	admin.Use(gin.Logger())
	admin.GET("/_health", NewHeartBeatHandler())
	admin.GET("/_config", NewConfigHandler(config))

	return admin
}

//NewAdminServer returns an http server with the default admin configuration
func NewAdminServer(config interface{}, port int, rTimeout int, wTimeout int) *http.Server {
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      NewDefaultAdminHandler(config),
		ReadTimeout:  time.Duration(rTimeout) * time.Second,
		WriteTimeout: time.Duration(wTimeout) * time.Second,
	}
	return server
}

//NewDefaultAdminServer returns a server with the given configuration that listens on port 8080
//with a read timeout of 5 and write timeout of 10
func NewDefaultAdminServer(config interface{}) *http.Server {
	return NewAdminServer(config, 8081, 5, 10)
}
