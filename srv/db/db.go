package dbsrv

import (
	"errors"
	"fmt"
	"github.com/reiver/logjam/cfg"
	"github.com/reiver/logjam/lib/db"
	"strings"
	"time"
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
				Type: db.TextType,
			},
			{
				Name:   "username",
				Type:   db.TextType,
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
				Type: db.TextType,
			},
		})
		if err != nil {
			return err
		}

		err = Repository.CreateTableIfNotExists("otps", []db.Field{
			{
				Name: "email",
				Type: db.EmailType,
			},
			{
				Name: "otp",
				Type: db.TextType,
			},
		})
		if err != nil {
			return err
		}

		err = Repository.CreateTableIfNotExists("blueskySessions", []db.Field{
			{
				Name:   "ownerId",
				Type:   db.TextType,
				Unique: true,
			},
			{
				Name: "accessToken",
				Type: db.TextType,
			},
			{
				Name: "refreshToken",
				Type: db.TextType,
			},
			{
				Name:   "did",
				Type:   db.TextType,
				Unique: true,
			},
			{
				Name: "handle",
				Type: db.TextType,
			},
			{
				Name: "email",
				Type: db.EmailType,
			},
			{
				Name: "serviceAddress",
				Type: db.TextType,
			},
		})
		if err != nil {
			return err
		}

		err = Repository.CreateTableIfNotExists("neynarIDs", []db.Field{
			{
				Name: "ownerId",
				Type: db.TextType,
			},
			{
				Name:   "fid",
				Type:   db.UnsignedInteger64Type,
				Unique: true,
			},
			{
				Name: "signerUUID",
				Type: db.TextType,
			},
		})
		if err != nil {
			return err
		}

		err = Repository.CreateTableIfNotExists("rooms", []db.Field{
			{
				Name:   "UID",
				Type:   db.TextType,
				Unique: true,
			},
			{
				Name: "ownerId",
				Type: db.TextType,
			},
			{
				Name: "title",
				Type: db.TextType,
			},
			{
				Name: "description",
				Type: db.TextType,
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
				Name: "roomUID",
				Type: db.TextType,
			}, {
				Name: "sent",
				Type: db.BooleanType,
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
				Type: db.TextType,
			},
			{
				Name: "ownerId",
				Type: db.TextType,
			},
			{
				Name: "data",
				Type: db.TextType,
			},
			{
				Name: "type",
				Type: db.TextType,
			},
			{
				Name: "fileHash",
				Type: db.TextType,
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
				Type: db.TextType,
			},
			{
				Name: "message",
				Type: db.TextType,
			},
			{
				Name: "signature",
				Type: db.TextType,
			},
			{
				Name: "ens",
				Type: db.TextType,
			},
			{
				Name: "ownerId",
				Type: db.TextType,
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

type PBTime time.Time

func (t PBTime) MarshalJSON() ([]byte, error) {
	// convert back to time.Time
	tt := time.Time(t)
	if tt.IsZero() {
		return []byte(`""`), nil
	}
	// format with the same layout PocketBase uses
	s := tt.Format("2006-01-02T15:04:05.000Z")
	// wrap in quotes for valid JSON string
	return []byte(`"` + s + `"`), nil
}
func (t *PBTime) UnmarshalJSON(b []byte) error {
	if len(b) == 2 {
		return nil
	}
	s := strings.Trim(string(b), `"`)
	// swap the space for a “T”
	s = strings.Replace(s, " ", "T", 1)
	// parse with high‑precision RFC3339
	tt, err := time.Parse(time.RFC3339Nano, s)
	if err != nil {
		return err
	}
	*t = PBTime(tt)
	return nil
}
