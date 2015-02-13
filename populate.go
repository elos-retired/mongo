package mongo

import (
	"github.com/elos/data"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func (db *MongoDB) PopulateByID(m data.Record) error {
	s, err := db.forkSession()
	if err != nil {
		return db.err(err)
	}
	defer s.Close()

	if err = db.populateById(s, m); err != nil {
		db.Printf("There was an error populating the %s model, error: %v", m.Kind(), err)
		if err == mgo.ErrNotFound {
			return data.ErrNotFound
		} else {
			return err
		}
	} else {
		return nil
	}
}

// Populates the model data for an empty struct with a specified id
func (db *MongoDB) populateById(s *mgo.Session, m data.Record) error {
	collection, err := db.collectionFor(s, m)
	if err != nil {
		return err
	}

	id, ok := m.ID().(bson.ObjectId)
	if !ok || !id.Valid() {
		return data.ErrInvalidID
	}

	return collection.FindId(id).One(m)
}

func (db *MongoDB) PopulateByField(field string, value interface{}, m data.Record) error {
	s, err := db.forkSession()
	if err != nil {
		return db.err(err)
	}
	defer s.Close()

	if err = db.populateByField(s, m, field, value); err != nil {
		db.Printf("There was an error populating the %s model, error: %v", m.Kind(), err)
		return err
	} else {
		return nil
	}
}

func (db *MongoDB) populateByField(s *mgo.Session, m data.Record, field string, value interface{}) error {
	collection, err := db.collectionFor(s, m)
	if err != nil {
		return err
	}

	return collection.Find(bson.M{field: value}).One(m)
}
