package mongo

import (
	"github.com/elos/data"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func (db *MongoDB) Delete(m data.Record) error {
	if m.DBType() != db.Type() {
		return data.ErrInvalidDBType
	}

	s, err := db.forkSession()
	if err != nil {
		return db.err(err)
	}
	defer s.Close()

	if err = db.remove(s, m); err != nil {
		db.Printf("Error deleted record of kind %s, err: %s", m.Kind(), err)
		return err
	} else {
		db.Notify(data.NewDelete(m))
		return nil
	}
}

func (db *MongoDB) remove(s *mgo.Session, m data.Record) error {
	collection, err := db.collectionFor(s, m)
	if err != nil {
		return err
	}

	id, ok := m.ID().(bson.ObjectId)
	if !ok || !id.Valid() {
		return data.ErrInvalidID
	}

	return collection.RemoveId(id)
}
