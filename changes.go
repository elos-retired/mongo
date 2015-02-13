package mongo

import "github.com/elos/data"

func (db *MongoDB) RegisterForChanges(client data.Client) *chan *data.Change {
	db.Lock()
	defer db.Unlock()

	id := client.ID()
	c := make(chan *data.Change)

	alreadySubscribed, ok := db.subscribers[id]
	if !ok {
		alreadySubscribed = make([]*chan *data.Change, 0)
	}

	db.subscribers[id] = append(alreadySubscribed, &c)

	return &c
}

func (db *MongoDB) notify(c *data.Change) {
	db.Lock()
	defer db.Unlock()
	// FIXME need a read/write lock here

	for _, concernedId := range c.Concerned() {
		channels, ok := db.subscribers[concernedId]
		if ok {
			for _, channel := range channels {
				go send(*channel, c)
			}
		}
	}
}

func send(channel chan *data.Change, change *data.Change) {
	channel <- change
}
