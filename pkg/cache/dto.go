package cache

import "time"

type UserReadModel struct {
	id        uint64
	firstName string
	lastName  string
}

type objectCache struct {
    object *UserReadModel
    expireAt time.Time
}
