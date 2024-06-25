package mongo

import (
	"github.com/kamva/mgm/v3"
)

type Storage struct {
	db *mgm.Collection
}

func New(connectUrl string) *Storage {
	const op = "storage.mongo.New"

	db, er = mgm.NewClient()

	return &Storage{
		db: db,
	}
}
