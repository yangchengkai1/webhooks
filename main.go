package main

import (
	"github.com/gin-gonic/gin"
	c "github.com/yangchengkai1/webhooks/controller"
)

func main() {
	router := gin.Default()

	c.RegisterRouter(router)

	router.Run(":8080")
}
