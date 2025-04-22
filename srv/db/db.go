package dbsrv

import (
	"errors"
	"fmt"
	"github.com/reiver/logjam/cfg"
	"github.com/reiver/logjam/lib/db"
)

var Repository db.IDBService

type TDBService string

const (
	PocketBase TDBService = "pb"
	MongoDB    TDBService = "mongodb"
	Postgres   TDBService = "postgres"
)

func Initialize(dbServiceType TDBService) error {
	switch dbServiceType {
	case PocketBase:
		Repository = db.NewPocketBaseDBService(cfg.Config.PocketBaseURL(), "")
	case MongoDB:
		return errors.New("mongodb service not implemented yet")
	case Postgres:
		return errors.New("postgres service not implemented yet")
	default:
		return errors.New(fmt.Sprintf("invalid db service type: %s", dbServiceType))
	}

	err := Repository.Ping()
	if err != nil {
		return err
	}

	//table initializations here
	{
		err := Repository.CreateTableIfNotExists("neynarIDs", []db.Field{
			{
				Name: "fid",
				Type: db.UnsignedInteger64,
			},
			{
				Name: "signerUUID",
				Type: db.String,
			},
		})
		if err != nil {
			return err
		}
	}

	return nil
}
