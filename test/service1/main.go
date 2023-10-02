package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

func main() {

	r := gin.New()
	defer func() {
		r.Run(":8899")
	}()

	r.GET("/test_service1", func(c *gin.Context) {
		fmt.Println("successful!!")
		c.JSON(200, gin.H{"ok": "ok"})
	})

}
