package mongo

import (
	"gopkg.in/mgo.v2"
)

type MongoConnection struct {
	Session *mgo.Session
}

func (mc *MongoConnection) Close() {
	if mc.Session != nil {
		mc.Session.Close()
		mc.Session = nil
	}
}

func Connect(addr string) (*MongoConnection, error) {
	session, err := mgo.Dial(addr)

	if err != nil {
		return nil, err
	}

	connection := &MongoConnection{Session: session}

	return connection, nil
}
