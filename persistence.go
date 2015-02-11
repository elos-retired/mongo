package mongo

import (
	"github.com/elos/data"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// Saves a model, broadcasted that save over ModelUpdates
func save(s *mgo.Session, m data.Record) error {
	collection, err := collectionFor(s, m)
	if err != nil {
		log(err)
		return err
	}

	id, ok := m.ID().(bson.ObjectId)
	if !ok {
		log("Model id was not of type bson.ObjectId")
	}

	if !id.Valid() {
		return data.ErrInvalidID
	}

	// changeInfo, err := ...
	_, err = collection.UpsertId(id, m)

	return err
}

func remove(s *mgo.Session, m data.Record) error {
	collection, err := collectionFor(s, m)
	if err != nil {
		log(err)
		return err
	}

	id, ok := m.ID().(bson.ObjectId)
	if !ok {
		log("Model id was not of the type bson.ObjectId")
	}

	if !id.Valid() {
		return data.ErrInvalidID
	}

	err = collection.RemoveId(id)
	return err
}

// Populates the model data for an empty struct with a specified id
func populateById(s *mgo.Session, m data.Record) error {
	collection, err := collectionFor(s, m)
	if err != nil {
		log(err)
		return err
	}

	id := m.ID()
	if !id.Valid() {
		return data.ErrInvalidID
	}

	return collection.FindId(m.ID()).One(m)
}

func populateByField(s *mgo.Session, m data.Record, field string, value interface{}) error {
	collection, err := collectionFor(s, m)
	if err != nil {
		log(err)
		return err
	}

	if err := collection.Find(bson.M{field: value}).One(m); err != nil {
		log(err)
		return err
	}

	return nil
}
