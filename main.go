package main

import (
	"log"

	"github.com/dfds/manage-aadgroup-members/docs"
	"github.com/dfds/manage-aadgroup-members/pkg/aadgroup"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	err := router.SetTrustedProxies([]string{"127.0.0.1"})
	if err != nil {
		log.Fatal(err)
	}
	docs.SwaggerInfo.BasePath = "/aadgroup/api/v1"
	v1 := router.Group("/aadgroup/api/v1")
	{
		v1.GET("/user/:upn", aadgroup.GetUser)
		//v1.DELETE("/user", aadgroup.RemoveUser)
		//v1.GET("/users", aadgroup.GetUsers)
		v1.POST("/users", aadgroup.AddUsers)
	}
	router.GET("/aadgroup/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	router.GET("/aadgroup/healthz", aadgroup.GetHealth)
	router.GET("/aadgroup/", aadgroup.GetIndex)
	err = router.Run("0.0.0.0:8000")
	if err != nil {
		log.Fatal(err)
	}
}
