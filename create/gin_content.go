package create

const ginContentTemplate = `package router

import (
	"net/http"

	"{{projectName}}/build"

	"github.com/gin-gonic/gin"
	"github.com/nichady/golte"
)

func wrapMiddleware(middleware *func(http.Handler) http.Handler, ctx *gin.Context) {
	if golte.GetRenderContext(func() *http.Request {
		(*middleware)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx.Request = r
			ctx.Next()
		})).ServeHTTP(ctx.Writer, ctx.Request)
		return ctx.Request
	}()) == nil {
		ctx.Abort()
	}
}

func GinRouter() http.Handler {
	// since gin doesm't use stdlib-compatible signatures, we have to wrap them
	// page := func(c string) gin.HandlerFunc {
	// 	return gin.WrapH(golte.Page(c))
	// }
	// layout := func(c string) gin.HandlerFunc {
	// 	return func(ctx *gin.Context) {
	// 		handler := golte.Layout(c)
	// 		wrapMiddleware(&handler, ctx)
	// 	}
	// }

	r := gin.Default()
	// register the main Golte middleware
	r.Use(func(ctx *gin.Context) {
		wrapMiddleware(&build.Golte, ctx)
	})

	defineRoutes(r)

	return r
}
`
const ginContentTemplate_sveltigo = `package router

import (
	"net/http"

	"{{projectName}}/build"

	"github.com/HazelnutParadise/sveltigo"
	"github.com/gin-gonic/gin"
)

func wrapMiddleware(middleware *func(http.Handler) http.Handler, ctx *gin.Context) {
	if sveltigo.GetRenderContext(func() *http.Request {
		(*middleware)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx.Request = r
			ctx.Next()
		})).ServeHTTP(ctx.Writer, ctx.Request)
		return ctx.Request
	}()) == nil {
		ctx.Abort()
	}
}

func GinRouter() http.Handler {
	// since gin doesm't use stdlib-compatible signatures, we have to wrap them
	// page := func(c string) gin.HandlerFunc {
	// 	return gin.WrapH(golte.Page(c))
	// }
	// layout := func(c string) gin.HandlerFunc {
	// 	return func(ctx *gin.Context) {
	// 		handler := golte.Layout(c)
	// 		wrapMiddleware(&handler, ctx)
	// 	}
	// }

	r := gin.Default()
	// register the main Golte middleware
	r.Use(func(ctx *gin.Context) {
		wrapMiddleware(&build.Sveltigo, ctx)
	})

	defineRoutes(r)

	return r
}
`
