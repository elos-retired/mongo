package mongo

import (
	"github.com/elos/data"
	"gopkg.in/mgo.v2"
)

type MongoQuery struct {
	db    *MongoDB
	kind  data.Kind
	am    data.AttrMap
	limit int
	skip  int
	batch int
}

func (db *MongoDB) NewQuery(k data.Kind) data.Query {
	return &MongoQuery{
		db:   db,
		kind: k,
		am:   data.AttrMap{},
	}
}

func (q *MongoQuery) Execute() (data.RecordIterator, error) {
	s, err := newSession(q.db)
	if err != nil {
		log(err)
		return nil, err
	}
	defer s.Close()
	// fixme: ugly type assertion
	return query(s, q)
}

func (q *MongoQuery) Select(am data.AttrMap) data.Query {
	q.am = am
	return q
}

func (q *MongoQuery) Limit(i int) data.Query {
	q.limit = i
	return q
}

func (q *MongoQuery) Skip(i int) data.Query {
	q.skip = i
	return q
}

func (q *MongoQuery) Batch(i int) data.Query {
	q.batch = i
	return q
}

type Iterator struct {
	iter *mgo.Iter
}

func (i *Iterator) Next(m data.Record) bool {
	return false
}

type MongoModelIterator struct {
	iter *mgo.Iter
}

func (i *MongoModelIterator) Next(m data.Record) bool {
	return i.iter.Next(m)
}

func (i *MongoModelIterator) Close() error {
	return i.iter.Close()
}

func query(s *mgo.Session, q *MongoQuery) (data.RecordIterator, error) {
	c, err := collectionForKind(s, q.kind)
	if err != nil {
		return nil, err
	}

	mgoQuery := c.Find(q.am)

	if q.limit != 0 {
		mgoQuery.Limit(q.limit)
	}

	if q.skip != 0 {
		mgoQuery.Skip(q.skip)
	}

	if q.batch != 0 {
		mgoQuery.Batch(q.batch)
	}

	return &MongoModelIterator{iter: mgoQuery.Iter()}, nil
}
