package routers

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "web3Tarot-backend/docs"
	"web3Tarot-backend/middleware"
	v1 "web3Tarot-backend/routers/api/v1"
	"web3Tarot-backend/setting"
)

func InitRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	gin.SetMode(setting.ServerSetting.RunMode)

	//r.StaticFS("/upload/images", http.Dir(upload.GetImageFullPath()))
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	//r.POST("/upload", api.UploadImage)

	apiv1 := r.Group("/api/v1").Use(middleware.AuthMiddleware)
	{
		//user
		apiv1.PUT("/users/:address/actions/sign-in", v1.Login)
		apiv1.GET("/users/:address", v1.GetUser)
		apiv1.PUT("/users/:address/actions/set-e2ee-key", v1.SetKeyInfo)
		apiv1.PUT("/public-users/actions/get", v1.GetUserPublicInfo)

		apiv1.POST("/nonces", v1.GetNonce)
		//apiv1.GET("/getPropositionsByPIDs", v1.GetPropostionsByPid)
		//apiv1.GET("/getPropositionByType/:type", v1.GetPropositionByType)
		//apiv1.GET("/getPropositions", v1.GetPropostions)
		//apiv1.GET("/getProposition/:pid", v1.GetProposition)
		//apiv1.POST("/addProposition", v1.AddProposition)
		//apiv1.PUT("/propositions/:id", v1.EditProposition)
		//apiv1.DELETE("/propositions/:id", v1.DeleteProposition)

	}

	return r
}
