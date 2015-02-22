package mongo

import (
	"sync"

	"github.com/elos/data"
	"gopkg.in/mgo.v2"
)

func (db *MongoDB) NewQuery(k data.Kind) data.RecordQuery {
	return &MongoQuery{
		db:    db,
		kind:  k,
		match: data.AttrMap{},
		Mutex: new(sync.Mutex),
	}
}

type MongoQuery struct {
	db    *MongoDB
	kind  data.Kind
	match data.AttrMap
	limit int
	skip  int
	batch int
	*sync.Mutex
}

func (q *MongoQuery) Execute() (data.RecordIterator, error) {
	q.Lock()
	defer q.Unlock()

	s, err := q.db.forkSession()
	if err != nil {
		return nil, err
	}
	defer s.Close()

	c, err := q.db.collectionForKind(s, q.kind)
	if err != nil {
		return nil, err
	}

	mgoQuery := c.Find(q.match)

	if q.limit != 0 {
		mgoQuery.Limit(q.limit)
	}

	if q.skip != 0 {
		mgoQuery.Skip(q.skip)
	}

	if q.batch != 0 {
		mgoQuery.Batch(q.batch)
	}

	return newRecordIter(mgoQuery.Iter()), nil
}

func (q *MongoQuery) Select(am data.AttrMap) data.RecordQuery {
	q.Lock()
	defer q.Unlock()
	q.match = am
	return q
}

func (q *MongoQuery) Limit(i int) data.RecordQuery {
	q.Lock()
	defer q.Unlock()

	q.limit = i
	return q
}

func (q *MongoQuery) Skip(i int) data.RecordQuery {
	q.Lock()
	defer q.Unlock()

	q.skip = i
	return q
}

func (q *MongoQuery) Batch(i int) data.RecordQuery {
	q.Lock()
	defer q.Unlock()

	q.batch = i
	return q
}

type recordIter struct {
	iter *mgo.Iter
	*sync.Mutex
}

func newRecordIter(i *mgo.Iter) data.RecordIterator {
	return &recordIter{
		iter:  i,
		Mutex: new(sync.Mutex),
	}
}

func (i *recordIter) Next(r data.Record) bool {
	return i.iter.Next(r)
}

func (i *recordIter) Close() error {
	return i.iter.Close()
}
