package types

import "github.com/gin-gonic/gin"

// RouterGroup wraps gin.RouterGroup for DI
type RouterGroup struct {
	*gin.RouterGroup
}
