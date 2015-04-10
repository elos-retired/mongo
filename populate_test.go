package mongo

import (
	"testing"

	"github.com/elos/data"
	"github.com/elos/testing/expect"
)

func TestPopulate(t *testing.T) {
	testify(Runner)
	defer detestify(Runner)
	db := NewDB()
	db.Connect("localhost")
	db.Logger = NullLogger

	testString := "aksjdf"

	model := data.NewNullModel()
	model.String = testString
	model.SetDBType(DBType)
	db.RegisterKind(data.NullKind, "nulls")
	id := db.NewID()
	model.SetID(id)

	err := db.Save(model)
	expect.NoError("saving model", err, t)

	model = data.NewNullModel()

	if err := db.PopulateByID(model); err != data.ErrInvalidDBType {
		t.Errorf("PopulateByID should reject a model with an incorrect DBType")
	}

	model.SetDBType(DBType)

	if err := db.PopulateByID(model); err != data.ErrInvalidID {
		t.Errorf("PopulateByID should reject a model with an invalid ID")
	}

	model.SetID(id)

	if err := db.PopulateByID(model); err != nil {
		t.Errorf("PopulateByID should return nil on a valid model, but got %s", err.Error())
	}

	if model.String != testString {
		t.Errorf("PopulateByID failed to correctly populate model, got %s, wanted: %s", model.String, testString)
	}

	model = data.NewNullModel()

	if err := db.PopulateByField("string", testString, model); err != data.ErrInvalidDBType {
		t.Errorf("PopulateByID should reject a model with an incorrect DBType")
	}

	model.SetDBType(DBType)

	if err := db.PopulateByField("string", testString, model); err != nil {
		t.Errorf("PopulateByField should return nil on valid model, but go %s", err.Error)
	}

	if model.String != testString {
		t.Errorf("Expected %s, got %s", testString, model.String)
	}

	// We can't check the id for correctness in this implementation because the null model
	// is mongo independent
}
