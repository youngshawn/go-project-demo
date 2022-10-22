package main

import (
	"log"

	_ "github.com/youngshawn/go-project-demo/course/config"
	_ "github.com/youngshawn/go-project-demo/course/models"
	"github.com/youngshawn/go-project-demo/course/routes"

	"github.com/gin-gonic/gin"
)

func main() {

	router := gin.Default()
	router.SetTrustedProxies([]string{"127.0.0.1"})

	routes.InstallRoutes(router)

	log.Fatal(router.Run(":3000"))

}
