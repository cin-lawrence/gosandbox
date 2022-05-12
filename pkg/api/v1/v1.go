package v1

import (
	"github.com/gin-gonic/gin"
)

func NewV1Group(rg *gin.RouterGroup) {
	v1Group := rg.Group("/api/v1")
        NewV1AuthGroup(v1Group)
	NewV1UserGroup(v1Group)
	NewV1JobGroup(v1Group)
}
