package api

import (
	"Auth_service/api/handler"
	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "Auth_service/api/docs"
)

// @title Authenfication service
// @version 1.0
// @description Server for signUp, signIn, forgot password and reset password
func NewRouter(hd *handler.Handler) *gin.Engine {

	router := gin.Default()
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.POST("/user/register", hd.RegisterStudent)
	router.POST("/user/accept-code", hd.AcceptCodeToRegister)

	router.POST("/user/login", hd.Login)
	router.POST("/user/forgot-password", hd.ForgotPassword)
	router.POST("/user/reset-password", hd.ResetPassword)

	router.POST("/admin/register", hd.RegisterAdmin)

	return router
}
