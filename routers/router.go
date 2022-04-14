package routers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Congregalis/gin-demo/middleware/jwt"
	"github.com/Congregalis/gin-demo/pkg/setting"
	"github.com/Congregalis/gin-demo/pkg/upload"
	v1 "github.com/Congregalis/gin-demo/routers/api/v1"
)

func InitRouter() *gin.Engine {
	r := gin.New()

	r.Use(gin.Logger(), gin.Recovery())

	gin.SetMode(setting.ServerSetting.RunMode)

	r.StaticFS(setting.AppSetting.ImageShowPath, http.Dir(upload.GetImageFullPath()))

	r.GET("/auth", v1.GetAuth)
	r.POST("/upload", v1.UploadImage)

	apiv1 := r.Group("/api/v1")
	apiv1.Use(jwt.JWT())
	{
		apiv1.GET("/tags", v1.GetTags)
		apiv1.POST("/tags", v1.AddTag)
		apiv1.PUT("/tags/:id", v1.EditTag)
		apiv1.DELETE("/tags/:id", v1.DeleteTag)

		apiv1.GET("/articles", v1.GetArticles)
		apiv1.GET("/articles/:id", v1.GetArticle)
		apiv1.POST("/articles", v1.AddArticle)
		apiv1.PUT("/articles/:id", v1.EditArticle)
		apiv1.DELETE("/articles/:id", v1.DeleteArticle)
	}

	return r
}
