package Mongo

import (
	"gopkg.in/mgo.v2"
)

type M struct {
	DBName   string
	DBServer string
}

var singleton *M

func GetDB() *mgo.Database {
	if singleton == nil {
		singleton = &M{DBName: "test", DBServer: "10.0.0.2:27017"}
	}

	return singleton.connect()
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
