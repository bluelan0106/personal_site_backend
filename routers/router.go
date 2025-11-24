package routers

import (
	"personal_site/config"
	"personal_site/controllers"
	"personal_site/middlewares"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Router interface {
	RegisterRoutes(r *gin.RouterGroup, db *gorm.DB)
}

func RegisterRouters(r *gin.Engine, db *gorm.DB) {
	apiPathPrefix, _ := config.GetVariableAsString("API_PATH_PREFIX")
	mainRouter := r.Group(apiPathPrefix)

	var authRouterVal Router = authRouter{}
	authRouterVal.RegisterRoutes(mainRouter.Group("/auth"), db)

	var storageRouterVal Router = storageRouter{}
	storageRouterVal.RegisterRoutes(mainRouter.Group("/storage"), db)

	var battleCatRouterVal Router = battleCatRouter{}
	battleCatRouterVal.RegisterRoutes(mainRouter.Group("/battle-cat"), db)

	var reurlRouterVal Router = reurlRouter{}
	reurlRouterVal.RegisterRoutes(mainRouter.Group("/reurl"), db)

	var postRouterVal Router = postRouter{}
	postRouterVal.RegisterRoutes(mainRouter, db)

	mainRouter.GET("/get-yt-data-api-token", middlewares.AuthOptional(), func(c *gin.Context) {
		controllers.GetYTDataAPIToken(c, db)
	})
}
