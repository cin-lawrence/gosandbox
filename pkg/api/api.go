package api

import (
	"fmt"
	"net/http"

	v1 "github.com/cin-lawrence/gosandbox/pkg/api/v1"
	"github.com/cin-lawrence/gosandbox/pkg/db"
	"github.com/cin-lawrence/gosandbox/pkg/validator"

	"github.com/fvbock/endless"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type APIServer struct {
	Router *gin.Engine
}

func NewAPIServer() APIServer {
        router := gin.Default()
        binding.Validator = new(validator.DefaultValidator)
        router.Use(v1.CORS())
        router.Use(v1.GenerateRequestID())
        router.Use(gzip.Gzip(gzip.DefaultCompression))

	server := APIServer{
		Router: router,
	}

	rg := server.Router.Group("/")
	rg.GET("/healthz", liveness)
	rg.GET("/healthz/readiness", readiness)
	v1.NewV1Group(rg)

        router.NoRoute(func(ctx *gin.Context) {
                ctx.JSON(http.StatusNotFound, gin.H{"message": "not found"})
        })

	return server
}

func liveness(ctx *gin.Context) {
	ctx.String(http.StatusOK, "")
}

func readiness(ctx *gin.Context) {
	dbConn, err := db.DB.DB()
	if err != nil {
		ctx.String(http.StatusInternalServerError, "Can't connect to DB")
	}
	err = dbConn.Ping()
	if err != nil {
		ctx.String(http.StatusInternalServerError, "Can't ping DB")
	}
	ctx.String(http.StatusOK, "OK")
}

func (server APIServer) Run(host string, port int) error {
	addr := fmt.Sprintf("%s:%d", host, port)
	err := endless.ListenAndServe(addr, server.Router)
	return err
}
