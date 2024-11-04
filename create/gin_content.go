package create

const ginContentTemplate = `package router

import (
	"net/http"

	"{{projectName}}/build"

	"github.com/gin-gonic/gin"
	"github.com/nichady/golte"
)

func GinRouter() http.Handler {
	// Gin doesn't have a function to wrap middleware, so define our own
	wrapMiddleware := func(middleware func(http.Handler) http.Handler) func(ctx *gin.Context) {
		return func(ctx *gin.Context) {
			middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				ctx.Request = r
				ctx.Next()
			})).ServeHTTP(ctx.Writer, ctx.Request)
			if golte.GetRenderContext(ctx.Request) == nil {
				ctx.Abort()
			}
		}
	}

	// since gin doesm't use stdlib-compatible signatures, we have to wrap them
	// page := func(c string) gin.HandlerFunc {
	// 	return gin.WrapH(golte.Page(c))
	// }
	// layout := func(c string) gin.HandlerFunc {
	// 	return wrapMiddleware(golte.Layout(c))
	// }

	r := gin.Default()
	// register the main Golte middleware
	r.Use(wrapMiddleware(build.Golte))

	defineRoutes(r)

	return r
}
`
