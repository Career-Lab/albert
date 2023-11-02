package main

import (
	"albert/controllers"
	_ "albert/docs"
	"albert/middlewares"
	"albert/models"
	"os"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func init() {
	models.ConnectDataBase()
	models.MigrateTables()
}

// @title           CareerLab - Alfred
// @version         1.0
// @description     Alfred is our rest endpoint service
// @host      		localhost:8080
// @BasePath  		/v1/
func main() {
	gin.SetMode(gin.DebugMode)

	router := gin.Default()

	router.ForwardedByClientIP = true
	router.SetTrustedProxies([]string{"127.0.0.1"})

	router.GET("/ping", controllers.Ping)
	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	v1Public := router.Group("/v1")
	authRouter := v1Public.Group("/auth")
	authRouter.POST("/signup", controllers.SignUp)
	authRouter.POST("/signin", controllers.SignIn)

	// Private routes
	v1Private := router.Group("/v1")
	v1Private.Use(middlewares.JwtAuthMiddleware())
	// Admin Routes
	adminRouter := router.Group("/v1")
	adminRouter.Use(middlewares.AdminAuthMiddleware())
	adminRouter.GET("/users", controllers.GetCountOfUsers)

	port := ":" + os.Getenv("PORT")
	router.Run(port)
}
