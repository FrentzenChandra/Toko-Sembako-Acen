package main

import (
	"log"

	"github.com/gin-gonic/gin"
)

var port = "192.168.18.5:9090"

func main() {

	r := gin.Default()
	r.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "Pong",
		})
	})

	log.Println("Running in Port : " + port)
	r.Run(port)

}
