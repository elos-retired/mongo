package mongo

import (
	"log"
	"os"
	"sync"

	"github.com/elos/data"
	"gopkg.in/mgo.v2"
)

type (
	CollectionMap map[data.Kind]string
	DBName        string
)

const (
	DBType      data.DBType = "mongo"
	DefaultName DBName      = "test"
)

type MongoDB struct {
	collections CollectionMap
	subscribers map[data.ID][]*chan *data.Change
	log         *log.Logger
	connection  *MongoConnection
	Name        DBName
	*sync.Mutex
}

func NewDB() (db *MongoDB) {
	db = &MongoDB{}
	db.collections = make(CollectionMap)
	db.subscribers = make(map[data.ID][]*chan *data.Change)
	db.log = log.New(os.Stdout, "[MONGO]", log.Lshortfile)
	db.Name = DefaultName
	db.Mutex = new(sync.Mutex)
	return
}

func (db *MongoDB) SetName(n DBName) {
	db.Lock()
	defer db.Unlock()
	db.Name = n
}

func (db *MongoDB) Type() data.DBType {
	return DBType
}

func (db *MongoDB) RegisterKind(k data.Kind, collection string) {
	db.collections[k] = collection
}

func (db *MongoDB) database(s *mgo.Session) *mgo.Database {
	return s.DB(string(db.Name))
}

func (db *MongoDB) collectionForKind(s *mgo.Session, k data.Kind) (*mgo.Collection, error) {
	c, ok := db.collections[k]
	if !ok {
		return nil, data.ErrUndefinedKind
	}
	return db.database(s).C(c), nil
}

func (db MongoDB) collectionFor(s *mgo.Session, r data.Record) (*mgo.Collection, error) {
	return db.collectionForKind(s, r.Kind())
}

// Forks the session of the primary connection
//		- If the PrimaryConnection does not exist, this returns a nil session
func (db *MongoDB) forkSession() (*mgo.Session, error) {
	if db.connection != nil {
		return db.connection.Session.Copy(), nil
	} else {
		return nil, data.ErrNoConnection
	}
}

func (db *MongoDB) err(err error) error {
	db.log.Print("error: ", err.Error())
	return err
}
