package mongo

import (
	"time"

	"github.com/elos/data"
	"gopkg.in/mgo.v2/bson"
)

type Model struct {
	id          `bson:",inline"`
	Timestamped `bson:",inline"`
}

func (m *Model) DBType() data.DBType {
	return DBType
}

func (m *Model) Valid() bool {
	return true
}

type id struct {
	EID bson.ObjectId `json:"id" bson:"_id,omitempty"`
}

func (m *id) SetID(id data.ID) error {
	vid, ok := id.(bson.ObjectId)
	if !ok {
		return data.ErrInvalidID
	}
	m.EID = vid
	return nil
}

func (m *id) ID() data.ID {
	return m.EID
}

type Timestamped struct {
	ECreatedAt time.Time `json:"created_at" bson:"created_at"`
	EUpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}

func (t *Timestamped) SetCreatedAt(ca time.Time) {
	t.ECreatedAt = ca
}

func (t *Timestamped) CreatedAt() time.Time {
	return t.ECreatedAt
}

func (t *Timestamped) SetUpdatedAt(ua time.Time) {
	t.EUpdatedAt = ua
}

func (t *Timestamped) UpdatedAt() time.Time {
	return t.EUpdatedAt
}
