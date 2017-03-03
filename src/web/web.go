package web

import (
	"gopkg.in/gin-gonic/gin.v1"
	"net/http"
)


func Listen() {
	router := gin.Default()
	router.LoadHTMLGlob("templates/*")
	router.Static("/static", "static")
	router.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", nil)
	})

	http.ListenAndServe(":80", router)
	//router.Run(":80")
}
