package rest

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/reyhanmichies/go-rest-api-boiler-plate/docs/swagger"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func (r *rest) registerSwaggerRoutes() {
	if r.ginConfig.Swagger.Enabled {
		swagger.SwaggerInfo.Title = r.ginConfig.Meta.Title
		swagger.SwaggerInfo.Description = r.ginConfig.Meta.Description
		swagger.SwaggerInfo.Version = r.ginConfig.Meta.Version
		swagger.SwaggerInfo.Host = r.ginConfig.Meta.Host
		swagger.SwaggerInfo.BasePath = r.ginConfig.Meta.BasePath

		swaggerAuth := gin.Accounts{
			r.ginConfig.Swagger.BasicAuth.Username: r.ginConfig.Swagger.BasicAuth.Password,
		}

		r.http.GET(fmt.Sprintf("%s/*any", r.ginConfig.Swagger.Path),
			gin.BasicAuthForRealm(swaggerAuth, "Restricted"),
			ginSwagger.WrapHandler(swaggerfiles.Handler))
	}
}
