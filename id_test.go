package mongo

import (
	"math/rand"
	"testing"
)

func TestNewObjectID(t *testing.T) {
	id := NewObjectID()

	if !id.Valid() {
		t.Errorf("NewObjectID should return a valid id")
	}
}

func TestIDSet(t *testing.T) {
	ids := make(IDSet, 0)

	ids = AddID(AddID(ids, NewObjectID()), NewObjectID())

	for i := 0; i < rand.Intn(10); i++ {
		ids = AddID(ids, NewObjectID())
	}

	id := NewObjectID()

	ids = AddID(ids, id)
	_, ok := ids.IndexID(id)
	if !ok {
		t.Errorf("Id be in set")
	}

	ids = AddID(ids, NewObjectID())
	ids = DropID(ids, id)

	_, ok = ids.IndexID(id)
	if ok != false {
		t.Errorf("Id should no longer be in set")
	}

	ids = make(IDSet, 0)
	i, ok := ids.IndexID(id)
	if i != -1 || ok != false {
		t.Errorf("IndexID return -1, false for a non-member id")
	}

	ids = make(IDSet, 0)
	ids = AddID(ids, id)
	ids = AddID(ids, id)
	ids = AddID(ids, id)

	if len(ids) != 1 {
		t.Errorf("IDSet should be a set, no duplicates")
	}

	ids = make(IDSet, 0)
	ids = DropID(ids, id)
	if len(ids) != 0 {
		t.Errorf("Should be able to drop a non-member id")
	}
}
