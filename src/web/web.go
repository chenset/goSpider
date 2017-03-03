package web

import (
	"gopkg.in/gin-gonic/gin.v1"
)

func Listen() {
	//gin.SetMode(gin.ReleaseMode)
	g := gin.Default()
	g.GET("/", func(c *gin.Context) {
		c.String(200, "success")
	})
	g.Run(":80") // listen and serve on 0.0.0.0:8080
}
