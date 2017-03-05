package web

import (
	"gopkg.in/gin-gonic/gin.v1"
	"net/http"
)

func Listen() {
	r := gin.Default()
	r.LoadHTMLGlob("templates/*")
	r.Static("/static", "static")

	// default route
	r.NoRoute(func(c *gin.Context) {
		c.HTML(200, "index.html", nil)
	})
	r.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", nil)
	})

	http.ListenAndServe(":80", r)
	//router.Run(":80")
}
