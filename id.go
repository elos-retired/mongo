package mongo

import (
	"sync"

	"github.com/elos/data"
	"gopkg.in/mgo.v2/bson"
)

func (db *MongoDB) NewID() data.ID {
	return NewObjectID()
}

func (db *MongoDB) CheckID(id data.ID) error {
	id, ok := id.(bson.ObjectId)
	if !ok || !id.Valid() {
		return data.ErrInvalidID
	} else {
		return nil
	}
}

func NewObjectID() data.ID {
	return bson.NewObjectId()
}

func NewObjectIDFromHex(hex string) data.ID {
	return bson.ObjectIdHex(hex)
}

func IsObjectIDHex(hex string) bool {
	return bson.IsObjectIdHex(hex)
}

type IDSet []bson.ObjectId

func AddID(s IDSet, id bson.ObjectId) IDSet {
	_, ok := s.IndexID(id)

	if !ok {
		ns := append(s, id)
		return ns
	}

	return s
}

func DropID(s IDSet, id bson.ObjectId) IDSet {
	i, ok := s.IndexID(id)

	if !ok {
		return s
	}

	s = s[:i+copy(s[i:], s[i+1:])]

	return s
}

func (s IDSet) IndexID(id bson.ObjectId) (int, bool) {
	for i, iid := range s {
		if id == iid {
			return i, true
		}
	}

	return -1, false
}

type IDIter struct {
	data.Store
	ids   IDSet
	place int
	err   error

	*sync.Mutex
}

func NewIDIter(set IDSet, s data.Store) *IDIter {
	return &IDIter{
		place: 0,
		Store: s,
		ids:   set,
		Mutex: new(sync.Mutex),
	}
}

func (i *IDIter) Next(r data.Record) bool {
	if i.place >= len(i.ids) {
		return false
	}

	r.SetID(i.ids[i.place])

	if err := i.Store.PopulateByID(r); err != nil {
		i.err = err
		return false
	} else {
		i.place += 1
		return true
	}
}

func (i *IDIter) Close() error {
	return i.err
}
