package Mongo

import (
	"gopkg.in/mgo.v2"
	"log"
)

type M struct {
	DBName   string
	DBServer string
}

var singleton *mgo.Database

func GetDB() *mgo.Database {
	if singleton == nil {
		m := &M{DBName: "test", DBServer: "10.0.0.2:27017"}
		singleton = m.connect()
		log.Println("Connect MongoDB " + m.DBServer)
	}

	return singleton
}

func (o *M) connect() *mgo.Database {
	conn, err := mgo.Dial(o.DBServer)
	if err != nil {
		panic(err)
	}

	defer func() {
		if r := recover(); r != nil {
			conn.Clone()
		}
	}()

	// Optional. Switch the session to a monotonic behavior.
	conn.SetMode(mgo.Monotonic, true)
	return conn.DB(o.DBName)
}
