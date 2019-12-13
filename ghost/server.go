package ghost

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/thoas/stats"
	"golang.org/x/sync/errgroup"
)

var g errgroup.Group

//Stats is the statistics middleware
var Stats = stats.New()

//SysStatsHandler inspired by github.com/appleboy/gorush
func SysStatsHandler(c *gin.Context) {
	c.JSON(http.StatusOK, Stats.Data())
}

// StatMiddleware response time, status code count, etc.
func StatMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		beginning, recorder := Stats.Begin(c.Writer)
		c.Next()
		Stats.End(beginning, stats.WithRecorder(recorder))
	}
}

//NewDefaultHandler returns an http handler with //recovery and stats middleware
func NewDefaultHandler() http.Handler {
	handler := gin.New()
	handler.Use(gin.Recovery())
	handler.Use(StatMiddleware())
	handler.Use(gin.Logger())

	handler.GET("/stats", SysStatsHandler)
	handler.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return handler

}

//NewServer returns an http server
func NewServer(port int, rTimeout int, wTimeout int) *http.Server {
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      NewDefaultHandler(),
		ReadTimeout:  time.Duration(rTimeout) * time.Second,
		WriteTimeout: time.Duration(wTimeout) * time.Second,
	}
	return server
}

//NewDefaultServer returns a new server that listens on port 8080
func NewDefaultServer() *http.Server {
	return NewServer(8080, 5, 10)
}

//RunServers runs a given list of servers
func RunServers(admin *http.Server, server *http.Server) {
	g.Go(func() error {
		return admin.ListenAndServe()
	})

	g.Go(func() error {
		return server.ListenAndServe()
	})

	if err := g.Wait(); err != nil {
		log.Fatal(err)
	}
}
