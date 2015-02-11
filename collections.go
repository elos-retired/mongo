package mongo

import (
	"fmt"
	"github.com/elos/data"
	"gopkg.in/mgo.v2"
)

const DefaultDatabase string = "test"

const Database = DefaultDatabase

var Collections = map[data.Kind]string{}

func RegisterKind(k data.Kind, c string) {
	Collections[k] = c
}

func database(s *mgo.Session) *mgo.Database {
	return s.DB(DefaultDatabase)
}

func collectionForKind(s *mgo.Session, k data.Kind) (*mgo.Collection, error) {
	collectionForKind := Collections[k]

	if collectionForKind == "" {
		err := fmt.Errorf("No collection name has been specified for the kind %s", k)
		return nil, err
	}

	return database(s).C(collectionForKind), nil
}

func collectionFor(s *mgo.Session, m data.Record) (*mgo.Collection, error) {
	return collectionForKind(s, m.Kind())
}
