package mongo

import (
	"github.com/elos/data"
	"gopkg.in/mgo.v2/bson"
)

func NewObjectID() data.ID {
	return bson.NewObjectId()
}

func NewObjectIDFromHex(hex string) data.ID {
	return bson.ObjectIdHex(hex)
}

func IsObjectIDHex(hex string) bool {
	return bson.IsObjectIdHex(hex)
}
