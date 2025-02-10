package create

const defineRoutesContent = `package router

import (
	"github.com/gin-gonic/gin"
	"github.com/nichady/golte"
)

func defineRoutes(r *gin.Engine) {
	r.GET("/", func(ctx *gin.Context) {
		golte.RenderPage(ctx.Writer, ctx.Request, "pages/App", map[string]any{
			"title": "Golte",
		})
	})
}
`

const defineRoutesSveltigoContent = `package router

import (
	"github.com/HazelnutParadise/sveltigo"
	"github.com/gin-gonic/gin"
)

func defineRoutes(r *gin.Engine) {
	r.GET("/", func(ctx *gin.Context) {
		sveltigo.RenderPage(ctx.Writer, ctx.Request, "pages/App", map[string]any{
			"title": "Golte",
		})
	})
}
`
