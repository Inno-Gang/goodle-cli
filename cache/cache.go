package cache

import (
	"github.com/Inno-Gang/goodle-cli/where"
	"github.com/charmbracelet/log"
	"github.com/philippgille/gokv"
	"github.com/philippgille/gokv/bbolt"
	"github.com/philippgille/gokv/encoding"
	"path/filepath"
)

var dbs = make([]gokv.Store, 0)

type Empty struct{}

func New(name string) gokv.Store {
	options := bbolt.Options{
		BucketName: name,
		Path:       filepath.Join(where.Cache(), name),
		Codec:      encoding.GobCodec{},
	}

	client, err := bbolt.NewStore(options)
	if err != nil {
		panic(err)
	}

	dbs = append(dbs, client)

	return client
}

func Close() {
	for _, db := range dbs {
		if err := db.Close(); err != nil {
			log.Error("closing db", "err", err.Error())
		}
	}
}
