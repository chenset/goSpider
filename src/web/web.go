package web

import (
	"gopkg.in/gin-gonic/gin.v1"
	"gopkg.in/mgo.v2"
	"log"
	"gopkg.in/mgo.v2/bson"
)
type Person struct {
	Name string
	Phone string
}
func Listen() {
	router := gin.Default()
	router.LoadHTMLGlob("templates/*")

	session, err := mgo.Dial("10.0.0.2:27017")
	if err != nil {
		panic(err)
	}
	defer session.Close()

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)



	router.GET("/", func(c *gin.Context) {
		cc := session.DB("test").C("people")
		err = cc.Insert(&Person{"Ale", "+55 53 8116 9639"},
			&Person{"Cla", "+55 53 8402 8510"})
		if err != nil {
			log.Fatal(err)
		}

		result := Person{}
		err = cc.Find(bson.M{"name": "Ale"}).One(&result)
		if err != nil {
			log.Fatal(err)
		}
		c.HTML(200, "index.html", nil)
	})

	//gin.SetMode(gin.ReleaseMode)
	router.Run(":80")
}
