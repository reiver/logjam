package wallets

import (
	dbsrv "github.com/reiver/logjam/srv/db"
)

type WalletDTO struct {
	ID          string        `json:"id"`
	Address     string        `json:"address"`
	Message     string        `json:"message"`
	Signature   string        `json:"signature"`
	ENS         string        `json:"ens"`
	OwnerId     string        `json:"ownerId"`
	IsConnected bool          `json:"isConnected"`
	Created     dbsrv.PBTime  `json:"created"`
	Updated     *dbsrv.PBTime `json:"updated"`
}
type AddWalletDTO struct {
	Address     string `json:"address"`
	Message     string `json:"message"`
	Signature   string `json:"signature"`
	ENS         string `json:"ens"`
	OwnerId     string `json:"ownerId"`
	IsConnected bool   `json:"isConnected"`
}

type UpdateWalletDTO struct {
	OwnerId     string  `json:"-"`
	Address     *string `json:"address,omitempty"`
	Message     *string `json:"message,omitempty"`
	Signature   *string `json:"signature,omitempty"`
	ENS         *string `json:"ens,omitempty"`
	IsConnected *bool   `json:"isConnected,omitempty"`
}
