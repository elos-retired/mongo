package mongo

import (
	"github.com/elos/data"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func (db *MongoDB) Save(m data.Record) error {
	s, err := db.forkSession()
	if err != nil {
		return db.err(err)
	}
	defer s.Close()

	if err = db.save(s, m); err != nil {
		db.log.Printf("Error saving record of kind %s, err: %s", m.Kind(), err.Error())
		return err
	} else {
		db.notify(data.NewChange(data.Update, m))
		return nil
	}
}

func (db *MongoDB) save(s *mgo.Session, r data.Record) error {
	collection, err := db.collectionFor(s, r)
	if err != nil {
		return err
	}

	id, ok := r.ID().(bson.ObjectId)
	if !ok || !id.Valid() {
		return data.ErrInvalidID
	}

	_ /*changeInfo*/, err = collection.UpsertId(id, r)

	return err
}