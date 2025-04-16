package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/obnahsgnaw/api/service"
)

var FrontendMiddleware = map[string]func() gin.HandlerFunc{
	//
}
var FrontendMuxMiddleware = map[string]func() service.MuxRouteHandleFunc{
	//
}
var BackendMiddleware = map[string]func() gin.HandlerFunc{
	//
}
var BackendMuxMiddleware = map[string]func() service.MuxRouteHandleFunc{
	//
}
