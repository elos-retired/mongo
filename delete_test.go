package mongo

import (
	"testing"

	"github.com/elos/data"
)

func TestDelete(t *testing.T) {
	testify(Runner)
	defer detestify(Runner)

	db := NewDB()
	db.Connect("localhost")
	db.Logger = NullLogger

	testString := "asdfljlaksd"

	model := data.NewNullModel()
	model.String = testString

	if err := db.Delete(model); err != data.ErrInvalidDBType {
		t.Errorf("Delete should reject NullType")
	}

	model.SetDBType(DBType)

	if _, ok := db.Delete(model).(data.UndefinedKindError); !ok {
		t.Errorf("Delete should recognize bad kind")
	}

	db.RegisterKind(data.NullKind, "nulls")

	if err := db.Delete(model); err != data.ErrInvalidID {
		t.Errorf("Delete should recognize a bad id")
	}

	model.SetID(db.NewID())

	if err := db.Save(model); err != nil {
		t.Errorf("Save should go off fine")
	}

	// can assume it is there, testing in save_test

	if err := db.Delete(model); err != nil {
		t.Errorf("Delete should work, but errored: %s", err.Error())
	}

	r := data.NewNullModel()
	r.SetDBType(DBType)
	r.SetID(model.ID())

	if err := db.PopulateByID(r); err != data.ErrNotFound {
		t.Errorf("the delete should have removed the model")
	}
}
