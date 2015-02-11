package mongo

import (
	"github.com/elos/data"
	"gopkg.in/mgo.v2/bson"
)

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
	ids IDSet
	p   int
	err error
}

func NewIDIter(set IDSet, s data.Store) *IDIter {
	return &IDIter{
		Store: s,
		ids:   set,
	}
}

func (i *IDIter) Next(r data.Record) bool {
	if i.p >= len(i.ids) {
		return false
	}

	r.SetID(i.ids[i.p])
	err := i.Store.PopulateByID(r)

	if err != nil {
		i.err = err
		return false
	}

	return true
}

func (i *IDIter) Close() error {
	return i.err
}
