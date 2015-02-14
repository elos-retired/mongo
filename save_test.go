package mongo

import (
	"testing"

	"github.com/elos/data"
)

func TestSave(t *testing.T) {
	Runner.ConfigFile = "./test.conf"
	Runner.Logger = NullLogger
	go Runner.Start()
	defer func() {
		Runner.Stop()
		Runner.ConfigFile = ""
		Runner.Logger = DefaultLogger
	}()

	db := NewDB()
	db.Connect("localhost")
	db.Logger = NullLogger

	testString := "asdfasdfASfaDF2309581934"
	testInt := 12341234

	model := data.NewNullModel()
	model.String = testString
	model.Int = testInt

	if err := db.Save(model); err != data.ErrInvalidDBType {
		t.Errorf("Mongo should reject NullType")
	}

	model.SetDBType(DBType)

	if err := db.Save(model); err != data.ErrUndefinedKind {
		t.Errorf("Mongo should recognize bad kind")
	}

	db.RegisterKind(data.NullKind, "nulls")

	if err := db.Save(model); err != data.ErrInvalidID {
		t.Errorf("Save errores: %s", err.Error())
	}

	if err := model.SetID(data.NewNullID("adf")); err != nil {
		t.Errorf("Had an error tryna set a null id to model, err: %s", err.Error())
	}

	if err := db.Save(model); err != data.ErrInvalidID {
		t.Errorf("Mongo should not choke on bad ids")
	}

	if err := model.SetID(db.NewID()); err != nil {
		t.Errorf("Had trouble setting model id to a bson id")
	}

	if err := db.Save(model); err != nil {
		t.Errorf("Model failed to save, err: %s", err.Error())
	}

	// retrieval
	retrievedModel := data.NewNullModel()
	retrievedModel.SetDBType(DBType)
	retrievedModel.SetID(model.ID())

	if err := db.PopulateByID(retrievedModel); err != nil {
		t.Errorf("DB failed to populate model by id, %s", err.Error())
	}

	if retrievedModel.String != testString {
		t.Errorf("Error with data, wanted %s, got %s", testString, retrievedModel.String)
	}

	if retrievedModel.Int != testInt {
		t.Errorf("Error with data, wanted %d, got %d", testInt, retrievedModel)
	}
}
