package mongo

import (
	"testing"

	"github.com/elos/data"
)

func TestNewDB(t *testing.T) {
	db := NewDB()

	if db.Name != DefaultName {
		t.Errorf("Should start with default name")
	}

	if db.Type() != DBType {
		t.Errorf("Should return DBType")
	}

	db.RegisterKind(data.NullKind, "nulls")
}

func TestDBOpsWithoutConnection(t *testing.T) {
	db := NewDB()
	db.Logger = NullLogger

	model := data.NewNullModel()
	model.SetDBType(DBType)

	if err := db.Save(model); err != data.ErrNoConnection {
		t.Errorf("Save should return an error when there is no connection")
	}

	if err := db.Delete(model); err != data.ErrNoConnection {
		t.Errorf("Delete should return an error when there is no connection")
	}

	if err := db.PopulateByID(model); err != data.ErrNoConnection {
		t.Errorf("PopulateByID should return an error when there is no connection")
	}

	if err := db.PopulateByField("hello", "world", model); err != data.ErrNoConnection {
		t.Errorf("PopulateByField should return an error when there is no conneciton")
	}
}
