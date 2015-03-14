package mongo

import (
	"sync"

	"github.com/elos/data"
	"gopkg.in/mgo.v2/bson"
)

func EmptyID() data.ID {
	return *new(bson.ObjectId)
}

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

func NewObjectID() bson.ObjectId {
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

func (s IDSet) NotIn(i IDSet) IDSet {
	ids := make(IDSet, 0)
	for _, id := range s {
		_, in := i.IndexID(id)
		if !in {
			ids = append(ids, id)
		}
	}
	return ids
}

func (s IDSet) IDs() []data.ID {
	ids := make([]data.ID, len(s))
	for i, id := range s {
		ids[i] = id
	}
	return ids
}

type IDIter struct {
	data.Access
	ids   IDSet
	place int
	err   error

	*sync.Mutex
}

func NewIDIter(set IDSet, a data.Access) *IDIter {
	return &IDIter{
		place:  0,
		Access: a,
		ids:    set,
		Mutex:  new(sync.Mutex),
	}
}

func (i *IDIter) Next(m data.Model) bool {
	if i.place >= len(i.ids) {
		return false
	}

	m.SetID(i.ids[i.place])

	if err := i.Access.PopulateByID(m); err != nil {
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
