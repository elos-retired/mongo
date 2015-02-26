package mongo

import (
	"github.com/elos/data"
	"gopkg.in/mgo.v2/bson"
)

type Model struct {
	ID               `bson:",inline"`
	data.Timestamped `bson:",inline"`
}

func (m *Model) DBType() data.DBType {
	return DBType
}

func (m *Model) Valid() bool {
	return true
}

type ID struct {
	EID bson.ObjectId `json:"id" bson:"_id,omitempty"`
}

func (m *ID) SetID(id data.ID) error {
	vid, ok := id.(bson.ObjectId)
	if !ok {
		return data.ErrInvalidID
	}
	m.EID = vid
	return nil
}

func (m *ID) ID() data.ID {
	return m.EID
}
