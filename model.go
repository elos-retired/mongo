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

type Named struct {
	EName string `json:"name" bson:"name"`
}

func (n *Named) SetName(name string) {
	n.EName = name
}

func (n *Named) Name() string {
	return n.EName
}

type Timed struct {
	EStartTime time.Time `json:"start_time" bson:"start_time"`
	EEndTime   time.Time `json:"end_time" bson:"end_time"`
}

func (t *Timed) StartTime() time.Time {
	return t.EStartTime
}

func (t *Timed) SetStartTime(st time.Time) {
	t.EStartTime = st
}

func (t *Timed) EndTime() time.Time {
	return t.EEndTime
}

func (t *Timed) SetEndTime(et time.Time) {
	t.EEndTime = et
}
