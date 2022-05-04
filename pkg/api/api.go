package api

import (
	"fmt"
	"net/http"

	v1 "github.com/cin-lawrence/gosandbox/pkg/api/v1"
	"github.com/cin-lawrence/gosandbox/pkg/db"
	"github.com/fvbock/endless"
	"github.com/gin-gonic/gin"
)

type APIServer struct {
	Router *gin.Engine
}

func NewAPIServer() APIServer {
	server := APIServer{
		Router: gin.Default(),
	}

	rg := server.Router.Group("/")
	rg.GET("/healthz", liveness)
	rg.GET("/healthz/readiness", readiness)
	v1.NewV1Group(rg)

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
