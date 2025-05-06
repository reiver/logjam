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
		Repository = db.NewPocketBaseDBService(cfg.Config.PocketBaseURL(), cfg.Config.PocketBaseAuthToken())
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
		err := Repository.CreateTableIfNotExists("userstbl", []db.Field{
			{
				Name: "name",
				Type: db.StringType,
			},
			{
				Name:   "username",
				Type:   db.StringType,
				Unique: true,
			},
			{
				Name:   "email",
				Type:   db.EmailType,
				Unique: true,
			},
			{
				Name: "avatar",
				Type: db.FileType,
			},
			{
				Name: "bio",
				Type: db.StringType,
			},
		})
		if err != nil {
			return err
		}

		err = Repository.CreateTableIfNotExists("otps", []db.Field{
			{
				Name: "userId",
				Type: db.StringType,
			},
			{
				Name: "otp",
				Type: db.IntegerType,
			},
		})
		if err != nil {
			return err
		}

		err = Repository.CreateTableIfNotExists("blueskySessions", []db.Field{
			{
				Name:   "userId",
				Type:   db.StringType,
				Unique: true,
			},
			{
				Name: "accessToken",
				Type: db.StringType,
			},
			{
				Name: "refreshToken",
				Type: db.StringType,
			},
			{
				Name:   "did",
				Type:   db.StringType,
				Unique: true,
			},
			{
				Name: "handle",
				Type: db.StringType,
			},
			{
				Name: "email",
				Type: db.EmailType,
			},
			{
				Name: "serviceAddress",
				Type: db.StringType,
			},
		})
		if err != nil {
			return err
		}

		err = Repository.CreateTableIfNotExists("neynarIDs", []db.Field{
			{
				Name: "userId",
				Type: db.StringType,
			},
			{
				Name: "fid",
				Type: db.StringType,
			},
			{
				Name: "signerUUID",
				Type: db.StringType,
			},
		})
		if err != nil {
			return err
		}

		err = Repository.CreateTableIfNotExists("rooms", []db.Field{
			{
				Name:   "UID",
				Type:   db.StringType,
				Unique: true,
			},
			{
				Name: "ownerId",
				Type: db.StringType,
			},
			{
				Name: "title",
				Type: db.StringType,
			},
			{
				Name: "description",
				Type: db.StringType,
			},
			{
				Name: "thumbnail",
				Type: db.FileType,
			},
		})
		if err != nil {
			return err
		}

		err = Repository.CreateTableIfNotExists("scheduledMeetings", []db.Field{
			{
				Name: "scheduledTime",
				Type: db.DateType,
			},
			{
				Name: "roomId",
				Type: db.StringType,
			},
		})
		if err != nil {
			return err
		}

		err = Repository.CreateTableIfNotExists("images", []db.Field{
			{
				Name: "image",
				Type: db.FileType,
			},
		})
		if err != nil {
			return err
		}
		err = Repository.CreateTableIfNotExists("layouts", []db.Field{
			{
				Name: "name",
				Type: db.StringType,
			},
			{
				Name: "ownerId",
				Type: db.StringType,
			},
			{
				Name: "data",
				Type: db.StringType,
			},
			{
				Name: "type",
				Type: db.StringType,
			},
			{
				Name: "fileHash",
				Type: db.StringType,
			},
			{
				Name: "lastUsed",
				Type: db.BooleanType,
			},
		})
		if err != nil {
			return err
		}
		err = Repository.CreateTableIfNotExists("wallets", []db.Field{
			{
				Name: "address",
				Type: db.StringType,
			},
			{
				Name: "message",
				Type: db.StringType,
			},
			{
				Name: "signature",
				Type: db.StringType,
			},
			{
				Name: "ens",
				Type: db.StringType,
			},
			{
				Name: "ownerId",
				Type: db.StringType,
			},
			{
				Name: "isConnected",
				Type: db.BooleanType,
			},
		})
		if err != nil {
			return err
		}
	}

	return nil
}
