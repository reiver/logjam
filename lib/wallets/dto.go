package wallets

import "time"

type WalletDTO struct {
	ID          string     `json:"id"`
	Address     string     `json:"address"`
	Message     string     `json:"message"`
	Signature   string     `json:"signature"`
	ENS         string     `json:"ens"`
	OwnerId     string     `json:"ownerId"`
	IsConnected bool       `json:"isConnected"`
	Created     time.Time  `json:"created"`
	Updated     *time.Time `json:"updated"`
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
	ID          string  `json:"-"`
	OwnerId     string  `json:"-"`
	Address     *string `json:"address,omitempty"`
	Message     *string `json:"message,omitempty"`
	Signature   *string `json:"signature,omitempty"`
	ENS         *string `json:"ens,omitempty"`
	IsConnected *bool   `json:"isConnected,omitempty"`
}
